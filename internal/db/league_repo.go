package db

import (
	"context"
	"errors"
	"strings"

	"github.com/ajarvis3/kickball-go/pkg/apperrors"

	"github.com/ajarvis3/kickball-go/internal/data/domain"
	"github.com/ajarvis3/kickball-go/internal/data/storage"
	"github.com/ajarvis3/kickball-go/internal/keys"
	"github.com/ajarvis3/kickball-go/internal/mappers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LeagueRepository interface {
	PutLeague(ctx context.Context, league domain.League) error
	GetLeague(ctx context.Context, leagueID string) (*domain.League, error)
	ListLeagues(ctx context.Context) ([]domain.League, error)
	ListLeaguesByName(ctx context.Context, namePrefix string) ([]domain.League, error)
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

func (r *leagueRepo) ListLeagues(ctx context.Context) ([]domain.League, error) {
	// Scan for items where PK begins with "LEAGUE#" and SK begins with "LEAGUE#"
	expr := "begins_with(PK, :pkprefix) AND begins_with(SK, :skprefix)"
	out, err := r.client.ddb.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.client.tableName),
		FilterExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkprefix": &types.AttributeValueMemberS{Value: "LEAGUE#"},
			":skprefix": &types.AttributeValueMemberS{Value: "LEAGUE#"},
		},
	})
	if err != nil {
		return nil, err
	}

	var leagues []domain.League
	for _, it := range out.Items {
		var stored storage.LeagueItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		leagues = append(leagues, mappers.ItemToLeague(stored))
	}
	return leagues, nil
}

func (r *leagueRepo) ListLeaguesByName(ctx context.Context, namePrefix string) ([]domain.League, error) {
	// Query the GSI that stores lowercase league names to support prefix search
	expr := "GSILeagueNamePK = :pk AND begins_with(GSILeagueNameSK, :prefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
		IndexName:              aws.String("GSILeagueByName"),
		KeyConditionExpression: aws.String(expr),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: "LEAGUE_NAME"},
			":prefix": &types.AttributeValueMemberS{Value: strings.ToLower(namePrefix)},
		},
	})
	if err != nil {
		return nil, err
	}

	var leagues []domain.League
	for _, it := range out.Items {
		var stored storage.LeagueItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		leagues = append(leagues, mappers.ItemToLeague(stored))
	}
	return leagues, nil
}
