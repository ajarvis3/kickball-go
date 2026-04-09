package db

import (
	"context"
	"fmt"

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

type LeagueRulesRepository interface {
	PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error
	GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error)
	GetLatestLeagueRules(ctx context.Context, leagueID string) (*domain.LeagueRules, error)
}

type leagueRulesRepo struct {
	client *Client
}

func NewLeagueRulesRepository(client *Client) LeagueRulesRepository {
	return &leagueRulesRepo{client: client}
}

func (r *leagueRulesRepo) PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error {
	// Convert to storage item
	it := mappers.LeagueRulesToItem(rules)
	item, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.client.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
	})

	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if fmt.Errorf("%w", err) != nil && (apperrors.ErrConflict != nil) {
			_ = cce // keep lint happy; falling through to default
		}
		// Use simple conflict detection similar to others
		var typedCce *types.ConditionalCheckFailedException
		if err != nil && (fmt.Sprintf("%T", err) == fmt.Sprintf("%T", typedCce)) {
			return apperrors.ErrConflict
		}
		return err
	}

	return nil
}

func (r *leagueRulesRepo) GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error) {
	sk := keys.RulesSK(rulesVersion)
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: keys.LeaguePK(leagueID)},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	out, err := r.client.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.client.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var stored storage.LeagueRulesItem
	if err := attributevalue.UnmarshalMap(out.Item, &stored); err != nil {
		return nil, err
	}
	lr := mappers.ItemToLeagueRules(stored)
	return &lr, nil
}

// GetLatestLeagueRules finds the highest rules version for a league and returns it.
func (r *leagueRulesRepo) GetLatestLeagueRules(ctx context.Context, leagueID string) (*domain.LeagueRules, error) {
	// Query for items where PK = LEAGUE#<leagueID> and SK begins_with("RULES#")
	keyCond := "PK = :pk AND begins_with(SK, :skprefix)"
	out, err := r.client.ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.client.tableName),
		KeyConditionExpression: aws.String(keyCond),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: keys.LeaguePK(leagueID)},
			":skprefix": &types.AttributeValueMemberS{Value: "RULES#"},
		},
	})
	if err != nil {
		return nil, err
	}
	var latest *domain.LeagueRules
	maxVersion := -1
	for _, it := range out.Items {
		var stored storage.LeagueRulesItem
		if err := attributevalue.UnmarshalMap(it, &stored); err != nil {
			return nil, err
		}
		// stored.RulesVersion is numeric in the item mapping
		if stored.RulesVersion > maxVersion {
			v := mappers.ItemToLeagueRules(stored)
			latest = &v
			maxVersion = stored.RulesVersion
		}
	}
	return latest, nil
}
