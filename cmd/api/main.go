package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type Request struct {
	// TODO: define request shape
}

type Response struct {
	// TODO: define response shape
}

func handler(ctx context.Context, req Request) (Response, error) {
	log.Println("kickball api invoked")
	// TODO: route to handlers
	return Response{}, nil
}