package db

import (
    "context"
    "errors"
    "fmt"
    "strconv"

    errors2 "github.com/ajarvis3/kickball-go/pkg/errors"

    "github.com/ajarvis3/kickball-go/internal/domain"

    "github.com/aws/aws-sdk-go-v2/aws"
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
	// TODO: marshal and PutItem
    item := map[string]types.AttributeValue{
        "PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LEAGUE#%s", league.LeagueID)},
        "SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LEAGUE#%s", league.LeagueID)},
        "leagueId": &types.AttributeValueMemberS{Value: league.LeagueID},
        "name": &types.AttributeValueMemberS{Value: league.Name},
        "currentRulesVersion": &types.AttributeValueMemberN{Value: strconv.Itoa(league.CurrentRulesVersion)},
    }

    _, err := r.client.ddb.PutItem(ctx, &dynamodb.PutItemInput{
        TableName:           aws.String(r.client.tableName),
        Item:                item,
        ConditionExpression: aws.String("attribute_not_exists(PK)"),
    })

    if err != nil {
        var cce *types.ConditionalCheckFailedException
        if errors.As(err, &cce) {
            return errors2.ErrConflict
        }
        return err
    }

    return nil
}

func (r *leagueRepo) GetLeague(ctx context.Context, leagueID string) (*domain.League, error) {
	// TODO: GetItem and unmarshal
	return nil, nil
}