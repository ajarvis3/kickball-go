package main

import (
    "context"
    "log"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"

    "github.com/ajarvis3/kickball-go/internal/db"
    "github.com/ajarvis3/kickball-go/internal/handlers"
)

var (
    atBatHandler   *handlers.AtBatHandlers
    gameHandler    *handlers.GameHandlers
    teamHandler    *handlers.TeamHandlers
    leagueHandler  *handlers.LeagueHandlers
    playerHandler  *handlers.PlayerHandlers
)

func init() {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        log.Fatalf("failed to load AWS config: %v", err)
    }

    ddb := dynamodb.NewFromConfig(cfg)

    tableName := os.Getenv("DYNAMODB_TABLE")
    if tableName == "" {
        log.Fatalf("DYNAMODB_TABLE environment variable is required")
    }

    // Create shared Dynamo client wrapper
    dynamoClient := db.NewClient(ddb, tableName)

    // Construct repositories
    atBatRepo := db.NewAtBatRepository(dynamoClient)
    gameRepo := db.NewGameRepository(dynamoClient)
    teamRepo := db.NewTeamRepository(dynamoClient)
    leagueRepo := db.NewLeagueRepository(dynamoClient)
    playerRepo := db.NewPlayerRepository(dynamoClient)

    // Construct handlers
    atBatHandler = handlers.NewAtBatHandlers(atBatRepo)
    gameHandler = handlers.NewGameHandlers(gameRepo)
    teamHandler = handlers.NewTeamHandlers(teamRepo)
    leagueHandler = handlers.NewLeagueHandlers(leagueRepo)
    playerHandler = handlers.NewPlayerHandlers(playerRepo)
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Received request: %s %s", req.HTTPMethod, req.Path)

    switch req.Path {
    case "/atbats":
        if req.HTTPMethod == "POST" {
            return atBatHandler.RecordAtBat(ctx, req)
        }

    case "/games":
        if req.HTTPMethod == "POST" {
            return gameHandler.CreateGame(ctx, req)
        }

    case "/teams":
        if req.HTTPMethod == "POST" {
            return teamHandler.CreateTeam(ctx, req)
        }

    case "/leagues":
        if req.HTTPMethod == "POST" {
            return leagueHandler.CreateLeague(ctx, req)
        }

    case "/players":
        if req.HTTPMethod == "POST" {
            return playerHandler.CreatePlayer(ctx, req)
        }
    }

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