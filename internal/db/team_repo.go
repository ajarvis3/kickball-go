package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TeamRepository interface {
	PutTeam(ctx context.Context, team domain.Team) error
	ListTeamsByLeague(ctx context.Context, leagueID string) ([]domain.Team, error)
}

type teamRepo struct {
	client *Client
}

func NewTeamRepository(client *Client) TeamRepository {
	return &teamRepo{client: client}
}

func (r *teamRepo) PutTeam(ctx context.Context, team domain.Team) error {
	it := storage.TeamToItem(team)

	item, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}

	_, err = r.client.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.client.tableName),
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepo) ListTeamsByLeague(ctx context.Context, leagueID string) ([]domain.Team, error) {
	pk := keys.LeaguePK(leagueID)
	expr := "PK = :pk AND begins_with(SK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "TEAM#"},
		},
	})
	if err != nil {
		return nil, err
	}
	var teams []domain.Team
	for _, it := range out.Items {
		var stored storage.TeamItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		teams = append(teams, storage.ItemToTeam(stored))
	}
	if len(teams) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return teams, nil
}
