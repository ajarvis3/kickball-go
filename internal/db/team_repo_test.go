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

func newTeamClient(ddb DynamoDBAPI) *Client {
	return NewClient(ddb, "test-table")
}

func TestPutTeamSuccess(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	repo := NewTeamRepository(newTeamClient(ddb))
	err := repo.PutTeam(context.Background(), domain.Team{TeamID: "t1", LeagueID: "l1", Name: "Tigers"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPutTeamError(t *testing.T) {
	ddb := &mockDynamoDB{
		putItemFn: func(_ context.Context, _ *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("put error")
		},
	}
	repo := NewTeamRepository(newTeamClient(ddb))
	err := repo.PutTeam(context.Background(), domain.Team{TeamID: "t1"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListTeamsByLeagueSuccess(t *testing.T) {
	it := storage.TeamItem{TeamID: "t1", LeagueID: "l1", Name: "Tigers"}
	item, _ := attributevalue.MarshalMap(it)
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil
		},
	}
	repo := NewTeamRepository(newTeamClient(ddb))
	teams, err := repo.ListTeamsByLeague(context.Background(), "l1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 1 || teams[0].TeamID != "t1" {
		t.Errorf("unexpected teams: %v", teams)
	}
}

func TestListTeamsByLeagueEmpty(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	repo := NewTeamRepository(newTeamClient(ddb))
	_, err := repo.ListTeamsByLeague(context.Background(), "l1")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListTeamsByLeagueError(t *testing.T) {
	ddb := &mockDynamoDB{
		queryFn: func(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("query error")
		},
	}
	repo := NewTeamRepository(newTeamClient(ddb))
	_, err := repo.ListTeamsByLeague(context.Background(), "l1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
