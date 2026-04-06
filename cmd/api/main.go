package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ajarvis3/kickball-go/code/internal/handlers"
)

var (
	atBatHandler   = handlers.NewAtBatHandler()
	gameHandler    = handlers.NewGameHandler()
	teamHandler    = handlers.NewTeamHandler()
	leagueHandler  = handlers.NewLeagueHandler()
	playerHandler  = handlers.NewPlayerHandler()
)

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", req.HTTPMethod, req.Path)

	switch req.Path {
	// ----- AtBats -----
	case "/atbats":
		if req.HTTPMethod == "POST" {
			return atBatHandler.CreateAtBat(ctx, req)
		}

	// ----- Games -----
	case "/games":
		if req.HTTPMethod == "POST" {
			return gameHandler.CreateGame(ctx, req)
		}

	// ----- Teams -----
	case "/teams":
		if req.HTTPMethod == "POST" {
			return teamHandler.CreateTeam(ctx, req)
		}

	// ----- Leagues -----
	case "/leagues":
		if req.HTTPMethod == "POST" {
			return leagueHandler.CreateLeague(ctx, req)
		}

	// ----- Players -----
	case "/players":
		if req.HTTPMethod == "POST" {
			return playerHandler.CreatePlayer(ctx, req)
		}
	}

	// Default 404
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       `{"error":"route not found"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(router)
}