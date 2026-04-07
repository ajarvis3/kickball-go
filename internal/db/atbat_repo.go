package db

import (
	"context"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/storage"
	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	// Convert domain to storage item
	it := storage.AtbatToItem(atbat)

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

func (r *atBatRepo) ListAtBatsByGame(ctx context.Context, gameID string) ([]domain.AtBat, error) {
	pk := keys.GamePK(gameID)
	expr := "PK = :pk AND begins_with(SK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
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
		var stored storage.AtbatItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		outAt = append(outAt, storage.ItemToAtbat(stored))
	}
	if len(outAt) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return outAt, nil
}

func (r *atBatRepo) ListAtBatsByPlayer(ctx context.Context, playerID string) ([]domain.AtBat, error) {
	// Query the GSI that indexes at-bats by player
	indexName := "GSIPlayerAtBat"
	pk := keys.GSI2PK(playerID)
	expr := "GSIPlayerAtBatPK = :pk"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}
	var outAt []domain.AtBat
	for _, it := range out.Items {
		var stored storage.AtbatItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		outAt = append(outAt, storage.ItemToAtbat(stored))
	}
	if len(outAt) == 0 {
		return nil, apperrors.ErrNotFound
	}
	return outAt, nil
}
