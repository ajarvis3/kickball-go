package db

import (
	"context"
	"errors"
	"testing"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/storage"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func newPlayerClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutPlayerSuccess(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewPlayerRepository(newPlayerClient(ddb))
	err := repo.PutPlayer(context.Background(), domain.Player{PlayerID: "p1", TeamID: "t1", Name: "Alice"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutPlayerError(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("put error")
		},
	}
	repo := NewPlayerRepository(newPlayerClient(ddb))
	err := repo.PutPlayer(context.Background(), domain.Player{PlayerID: "p1"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListPlayersByTeamSuccess(t *testing.T) {
	it := storage.PlayerItem{PlayerID: "p1", TeamID: "t1", Name: "Alice"}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
		},
	}
	repo := NewPlayerRepository(newPlayerClient(ddb))
	players, err := repo.ListPlayersByTeam(context.Background(), "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(players) != 1 || players[0].PlayerID != "p1" {
		t.Errorf("unexpected players: %v", players)
	}
}

func TestListPlayersByTeamEmpty(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	repo := NewPlayerRepository(newPlayerClient(ddb))
	players, err := repo.ListPlayersByTeam(context.Background(), "t1")
	if err != nil {
		t.Errorf("expected nil error for empty result, got %v", err)
	}
	if len(players) != 0 {
		t.Errorf("expected empty slice, got %v", players)
	}
}

func TestListPlayersByTeamError(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}
	repo := NewPlayerRepository(newPlayerClient(ddb))
	_, err := repo.ListPlayersByTeam(context.Background(), "t1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
