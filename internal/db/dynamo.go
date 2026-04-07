package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	ddb       *dynamodb.Client
	tableName string
}

func NewClient(ddb *dynamodb.Client, tableName string) *Client {
	return &Client{ddb: ddb, tableName: tableName}
}

// TODO: add helpers: GetItem, PutItem, Query, etc.
func (c *Client) Ping(ctx context.Context) error {
	// TODO: maybe a lightweight call or no-op
	return nil
}
