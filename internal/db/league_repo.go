package db

import (
	"context"
	"errors"

	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/ajarvis3/kickball-go/internal/domain"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/mappers"
	"github.com/ajarvis3/kickball-go/internal/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LeagueRepository interface {
	PutLeague(ctx context.Context, league domain.League) error
	GetLeague(ctx context.Context, leagueID string) (*domain.League, error)
}

type leagueRepo struct {
	client *Client
}

func NewLeagueRepository(client *Client) LeagueRepository {
	return &leagueRepo{client: client}
}

func (r *leagueRepo) PutLeague(ctx context.Context, league domain.League) error {
	it := mappers.LeagueToItem(league)
	item, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.client.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})

	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if errors.As(err, &cce) {
			return apperrors.ErrConflict
		}
		return err
	}

	return nil
}

func (r *leagueRepo) GetLeague(ctx context.Context, leagueID string) (*domain.League, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: keys.LeaguePK(leagueID)},
		"SK": &types.AttributeValueMemberS{Value: keys.LeagueSK(leagueID)},
	}
	out, err := r.client.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.client.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, apperrors.ErrNotFound
	}

	var stored storage.LeagueItem
	if err := attributevalue.UnmarshalMap(out.Item, &stored); err != nil {
		return nil, err
	}
	l := mappers.ItemToLeague(stored)
	return &l, nil
}
