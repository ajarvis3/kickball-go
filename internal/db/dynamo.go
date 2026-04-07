package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBAPI is the subset of dynamodb.Client methods used by the repositories.
// Using an interface here allows tests to inject mock implementations without
// creating real AWS resources.
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

type Client struct {
	ddb       DynamoDBAPI
	tableName string
}

func NewClient(ddb DynamoDBAPI, tableName string) *Client {
	return &Client{ddb: ddb, tableName: tableName}
}

// TODO: add helpers: GetItem, PutItem, Query, etc.
func (c *Client) Ping(ctx context.Context) error {
	// TODO: maybe a lightweight call or no-op
	return nil
}
