package db

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"

    "github.com/ajarvis3/kickball-go/pkg/apperrors"
    "github.com/ajarvis3/kickball-go/internal/domain"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LeagueRulesRepository interface {
    PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error
    GetLeagueRules(ctx context.Context, leagueID string, rulesVersion int) (*domain.LeagueRules, error)
}

type leagueRulesRepo struct {
    client *Client
}

func NewLeagueRulesRepository(client *Client) LeagueRulesRepository {
    return &leagueRulesRepo{client: client}
}

func (r *leagueRulesRepo) PutLeagueRules(ctx context.Context, rules domain.LeagueRules) error {
    // marshal the full rules into a JSON blob to keep storage simple
    blob, err := json.Marshal(rules)
    if err != nil {
        return err
    }

    item := map[string]types.AttributeValue{
        "PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LEAGUE#%s", rules.LeagueID)},
        "SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("RULES#%d", rules.RulesVersion)},
        "leagueId": &types.AttributeValueMemberS{Value: rules.LeagueID},
        "rulesVersion": &types.AttributeValueMemberN{Value: strconv.Itoa(rules.RulesVersion)},
        "rules": &types.AttributeValueMemberS{Value: string(blob)},
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
    sk := fmt.Sprintf("RULES#%d", rulesVersion)
    key := map[string]types.AttributeValue{
        "PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LEAGUE#%s", leagueID)},
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

    av, ok := out.Item["rules"].(*types.AttributeValueMemberS)
    if !ok {
        return nil, nil
    }

    var rules domain.LeagueRules
    if err := json.Unmarshal([]byte(av.Value), &rules); err != nil {
        return nil, err
    }

    return &rules, nil
}
