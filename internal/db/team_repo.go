package db

import (
	"context"
	"fmt"

	"github.com/ajarvis3/kickball-go/internal/domain"

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
	// Build main table keys
	team.PK = fmt.Sprintf("LEAGUE#%s", team.LeagueID)
	team.SK = fmt.Sprintf("TEAM#%s", team.TeamID)

	// Marshal the item
	item, err := attributevalue.MarshalMap(team)
	if err != nil {
		return err
	}

	// Write to DynamoDB
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
	pk := fmt.Sprintf("LEAGUE#%s", leagueID)
	expr := "PK = :pk AND begins_with(SK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(r.client.tableName),
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
		var t domain.Team
		if err := attributevalue.UnmarshalMap(it, &t); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}