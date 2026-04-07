package db

import (
	"context"
	"fmt"

	"github.com/ajarvis3/kickball-go/internal/domain"
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
	// Build main table keys
	game.PK = fmt.Sprintf("GAME#%s", game.GameID)
	game.SK = fmt.Sprintf("GAME#%s", game.GameID)

	// Build GSI keys for listing games by league
	game.GSILeagueGamePK = fmt.Sprintf("LEAGUE#%s", game.LeagueID)
	game.GSILeagueGameSK = fmt.Sprintf("GAME#%s", game.GameID)

	// Marshal the item
	item, err := attributevalue.MarshalMap(game)
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
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("GAME#%s", gameID)},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("GAME#%s", gameID)},
	}
	out, err := r.client.ddb.GetItem(ctx, &dynamodb.GetItemInput{TableName: aws.String(r.client.tableName), Key: key})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, apperrors.ErrNotFound
	}
	var g domain.Game
	if err := attributevalue.UnmarshalMap(out.Item, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *gameRepo) ListGamesByLeague(ctx context.Context, leagueID string) ([]domain.Game, error) {
	pk := fmt.Sprintf("LEAGUE#%s", leagueID)
	// Query for items where PK = LEAGUE#<leagueID> and SK begins_with "GAME#"
	expr := "GSILeagueGamePK = :GSILeagueGamePK"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(r.client.tableName),
		IndexName:              aws.String("GSILeagueGame"),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":GSILeagueGamePK":     &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}
	var games []domain.Game
	for _, it := range out.Items {
		var g domain.Game
		if err := attributevalue.UnmarshalMap(it, &g); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, nil
}