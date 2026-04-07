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

type PlayerRepository interface {
	PutPlayer(ctx context.Context, player domain.Player) error
	ListPlayersByTeam(ctx context.Context, teamID string) ([]domain.Player, error)
}

type playerRepo struct {
	client *Client
}

func NewPlayerRepository(client *Client) PlayerRepository {
	return &playerRepo{client: client}
}

func (r *playerRepo) PutPlayer(ctx context.Context, player domain.Player) error {
	// Build main table keys
	player.PK = fmt.Sprintf("TEAM#%s", player.TeamID)
	player.SK = fmt.Sprintf("PLAYER#%s", player.PlayerID)

	// Marshal the item
	item, err := attributevalue.MarshalMap(player)
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

func (r *playerRepo) ListPlayersByTeam(ctx context.Context, teamID string) ([]domain.Player, error) {
	pk := fmt.Sprintf("TEAM#%s", teamID)
	expr := "PK = :pk AND begins_with(SK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(r.client.tableName),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "PLAYER#"},
		},
	})
	if err != nil {
		return nil, err
	}
	var players []domain.Player
	for _, it := range out.Items {
		var p domain.Player
		if err := attributevalue.UnmarshalMap(it, &p); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}