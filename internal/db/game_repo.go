package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/mappers"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GameRepository interface {
	PutGame(ctx context.Context, game domain.Game) error
	GetGame(ctx context.Context, gameID string) (*domain.Game, error)
	ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error)
}

type gameRepo struct {
	client *Client
}

func NewGameRepository(client *Client) GameRepository {
	return &gameRepo{client: client}
}

func (r *gameRepo) PutGame(ctx context.Context, game domain.Game) error {
	// Convert domain -> storage item
	it := mappers.GameToItem(game)

	// Marshal the item
	item, err := attributevalue.MarshalMap(it)
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

func (r *gameRepo) GetGame(ctx context.Context, gameID string) (*domain.Game, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: keys.GamePK(gameID)},
		"SK": &types.AttributeValueMemberS{Value: keys.GameSK(gameID)},
	}
	out, err := r.client.ddb.GetItem(ctx, &dynamodb.GetItemInput{TableName: aws.String(r.client.tableName), Key: key})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, apperrors.ErrNotFound
	}
	var stored storage.GameItem
	if err := attributevalue.UnmarshalMap(out.Item, &stored); err != nil {
		return nil, err
	}
	g := mappers.ItemToGame(stored)
	return &g, nil
}

func (r *gameRepo) ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error) {
	pk := keys.LeaguePK(leagueID)
	// Query for items where PK = LEAGUE#<leagueID> and SK begins_with "GAME#"
	expr := "GSILeagueGamePK = :GSILeagueGamePK"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
		IndexName:              aws.String("GSILeagueGame"),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":GSILeagueGamePK": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}
	var games []domain.Game
	for _, it := range out.Items {
		var stored storage.GameItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		games = append(games, mappers.ItemToGame(stored))
	}
	return games, nil
}
