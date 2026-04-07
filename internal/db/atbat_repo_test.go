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

func newAtBatClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutAtBat_Success(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	err := repo.PutAtBat(context.Background(), domain.AtBat{GameID: "g1", PlayerID: "p1", Seq: 1})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutAtBat_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("put error")
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	err := repo.PutAtBat(context.Background(), domain.AtBat{GameID: "g1"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestPutAtBatAndUpdateGame_Success(t *testing.T) {
	ddb := &mockDynamoDB{
		transactWriteItemsFn: func(_ context.Context, _ *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
			return &dynamodb.TransactWriteItemsOutput{}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	err := repo.PutAtBatAndUpdateGame(context.Background(), domain.AtBat{GameID: "g1", Seq: 1}, domain.Game{GameID: "g1"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutAtBatAndUpdateGame_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		transactWriteItemsFn: func(_ context.Context, _ *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
			return nil, errors.New("tx error")
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	err := repo.PutAtBatAndUpdateGame(context.Background(), domain.AtBat{GameID: "g1"}, domain.Game{GameID: "g1"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListAtBatsByGame_Success(t *testing.T) {
	it := storage.AtbatItem{GameID: "g1", PlayerID: "p1", Seq: 1}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	atbats, err := repo.ListAtBatsByGame(context.Background(), "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(atbats) != 1 || atbats[0].GameID != "g1" {
		t.Errorf("unexpected atbats: %v", atbats)
	}
}

func TestListAtBatsByGame_Empty(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	_, err := repo.ListAtBatsByGame(context.Background(), "g1")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListAtBatsByGame_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	_, err := repo.ListAtBatsByGame(context.Background(), "g1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListAtBatsByPlayer_Success(t *testing.T) {
	it := storage.AtbatItem{GameID: "g1", PlayerID: "p1", Seq: 1}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	atbats, err := repo.ListAtBatsByPlayer(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(atbats) != 1 {
		t.Errorf("expected 1 atbat, got %d", len(atbats))
	}
}

func TestListAtBatsByPlayer_Empty(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	_, err := repo.ListAtBatsByPlayer(context.Background(), "p1")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListAtBatsByPlayer_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}
	repo := NewAtBatRepository(newAtBatClient(ddb))
	_, err := repo.ListAtBatsByPlayer(context.Background(), "p1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
