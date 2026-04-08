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

func newLeagueClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutLeague_Success(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	err := repo.PutLeague(context.Background(), domain.League{LeagueID: "l1", Name: "Test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutLeague_Conflict(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, &types.ConditionalCheckFailedException{}
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	err := repo.PutLeague(context.Background(), domain.League{LeagueID: "l1", Name: "Test"})
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestPutLeague_OtherError(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("other error")
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	err := repo.PutLeague(context.Background(), domain.League{LeagueID: "l1", Name: "Test"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrConflict) {
		t.Errorf("should not be ErrConflict for non-conditional error")
	}
}

func TestGetLeague_Success(t *testing.T) {
	it := storage.LeagueItem{
		LeagueID:            "l1",
		Name:                "Test League",
		CurrentRulesVersion: 1,
	}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: item}, nil
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	l, err := repo.GetLeague(context.Background(), "l1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.LeagueID != "l1" {
		t.Errorf("LeagueID = %q; want l1", l.LeagueID)
	}
	if l.Name != "Test League" {
		t.Errorf("Name = %q; want Test League", l.Name)
	}
}

func TestGetLeague_NotFound(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: nil}, nil
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	_, err := repo.GetLeague(context.Background(), "l1")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetLeague_Error(t *testing.T) {
	ddb := &mockDynamoDB{
		getItemFn: func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("get error")
		},
	}
	repo := NewLeagueRepository(newLeagueClient(ddb))
	_, err := repo.GetLeague(context.Background(), "l1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
