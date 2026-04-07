package db

import (
	"context"
	"fmt"

	"github.com/ajarvis3/kickball-go/internal/domain"
	apperrors "github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AtBatRepository interface {
	PutAtBat(ctx context.Context, atbat domain.AtBat) error
	ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error)
	ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error)
}

type atBatRepo struct {
	client *Client
}

func NewAtBatRepository(client *Client) AtBatRepository {
	return &atBatRepo{client: client}
}

func (r *atBatRepo) PutAtBat(ctx context.Context, atbat domain.AtBat) error {
    // Build main table keys
    atbat.PK = fmt.Sprintf("GAME#%s", atbat.GameID)
    atbat.SK = fmt.Sprintf("ATBAT#%d", atbat.Seq)

    // Build GSI keys for querying at-bats by player
    atbat.GSIPlayerAtBatPK = fmt.Sprintf("PLAYER#%s", atbat.PlayerID)
    atbat.GSIPlayerAtBatSK = fmt.Sprintf("GAME#%s#ATBAT#%d", atbat.GameID, atbat.Seq)

    // Marshal the item
    item, err := attributevalue.MarshalMap(atbat)
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

func (r *atBatRepo) ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error) {
	pk := fmt.Sprintf("GAME#%s", gameID)
	expr := "PK = :pk AND begins_with(SK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(r.client.tableName),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
			":prefix": &types.AttributeValueMemberS{Value: "ATBAT#"},
		},
	})
	if err != nil {
		return nil, err
	}
	var outAt []domain.AtBat
	for _, it := range out.Items {
		var a domain.AtBat
		if err := attributevalue.UnmarshalMap(it, &a); err != nil {
			return nil, err
		}
		outAt = append(outAt, a)
	}
	if len(outAt) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return outAt, nil
}

func (r *atBatRepo) ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error) {
	// Query the GSI that indexes at-bats by player
	indexName := "GSIPlayerAtBat"
	pk := fmt.Sprintf("PLAYER#%s", playerID)
	expr := "GSIPlayerAtBatPK = :pk"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(r.client.tableName),
		IndexName: aws.String(indexName),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}
	var outAt []domain.AtBat
	for _, it := range out.Items {
		var a domain.AtBat
		if err := attributevalue.UnmarshalMap(it, &a); err != nil {
			return nil, err
		}
		outAt = append(outAt, a)
	}
	if len(outAt) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return outAt, nil
}