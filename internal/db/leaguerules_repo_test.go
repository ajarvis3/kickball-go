package db

import (
	"context"
	"errors"
	"testing"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func newLeagueRulesClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutLeagueRulesSuccess(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewLeagueRulesRepository(newLeagueRulesClient(ddb))
	err := repo.PutLeagueRules(context.Background(), domain.LeagueRules{LeagueID: "l1", RulesVersion: 1, MaxStrikes: 3})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutLeagueRulesError(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("put error")
		},
	}
	repo := NewLeagueRulesRepository(newLeagueRulesClient(ddb))
	err := repo.PutLeagueRules(context.Background(), domain.LeagueRules{LeagueID: "l1", RulesVersion: 1})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetLeagueRulesSuccess(t *testing.T) {
	it := storage.LeagueRulesItem{
		LeagueID:     "l1",
		RulesVersion: 1,
		MaxStrikes:   3,
		MaxInnings:   7,
	}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: item}, nil
		},
	}
	repo := NewLeagueRulesRepository(newLeagueRulesClient(ddb))
	r, err := repo.GetLeagueRules(context.Background(), "l1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected rules, got nil")
	}
	if r.MaxStrikes != 3 {
		t.Errorf("MaxStrikes = %d; want 3", r.MaxStrikes)
	}
}

func TestGetLeagueRulesNotFoundReturnsNilNil(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: nil}, nil
		},
	}
	repo := NewLeagueRulesRepository(newLeagueRulesClient(ddb))
	r, err := repo.GetLeagueRules(context.Background(), "l1", 1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if r != nil {
		t.Errorf("expected nil rules, got %+v", r)
	}
}

func TestGetLeagueRulesError(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("get error")
		},
	}
	repo := NewLeagueRulesRepository(newLeagueRulesClient(ddb))
	_, err := repo.GetLeagueRules(context.Background(), "l1", 1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
