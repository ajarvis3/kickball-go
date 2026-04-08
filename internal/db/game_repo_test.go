package db

import (
	"context"
	"errors"
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/storage"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDB struct {
	putItemFn            func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	getItemFn            func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	queryFn              func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	transactWriteItemsFn func(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

func (m *mockDynamoDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFn(ctx, params, optFns...)
}

func (m *mockDynamoDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.getItemFn(ctx, params, optFns...)
}

func (m *mockDynamoDB) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return m.queryFn(ctx, params, optFns...)
}

func (m *mockDynamoDB) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return m.transactWriteItemsFn(ctx, params, optFns...)
}

func newGameClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutGame_Success(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	err := repo.PutGame(context.Background(), domain.Game{GameID: "g1", LeagueID: "l1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutGame_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("put error")
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	err := repo.PutGame(context.Background(), domain.Game{GameID: "g1"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetGame_Success(t *testing.T) {
	it := storage.GameItem{
		GameID:   "g1",
		LeagueID: "l1",
		State:    domain.GameState{Inning: 1, Half: "top"},
	}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: item}, nil
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	g, err := repo.GetGame(context.Background(), "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.GameID != "g1" {
		t.Errorf("GameID = %q; want g1", g.GameID)
	}
}

func TestGetGame_NotFound(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: nil}, nil
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	_, err := repo.GetGame(context.Background(), "g1")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetGame_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("get error")
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	_, err := repo.GetGame(context.Background(), "g1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListGamesByLeague_Success(t *testing.T) {
	it := storage.GameItem{GameID: "g1", LeagueID: "l1"}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	games, err := repo.ListGamesByLeague(context.Background(), "l1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(games) != 1 || games[0].GameID != "g1" {
		t.Errorf("unexpected games: %v", games)
	}
}

func TestListGamesByLeague_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}
	repo := NewGameRepository(newGameClient(ddb))
	_, err := repo.ListGamesByLeague(context.Background(), "l1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
