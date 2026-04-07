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
    leagueRulesHandler *handlers.LeagueRulesHandlers
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
    leagueRulesRepo := db.NewLeagueRulesRepository(dynamoClient)

    // Construct handlers
    atBatHandler = handlers.NewAtBatHandlers(atBatRepo)
    gameHandler = handlers.NewGameHandlers(gameRepo)
    teamHandler = handlers.NewTeamHandlers(teamRepo)
    leagueHandler = handlers.NewLeagueHandlers(leagueRepo)
    playerHandler = handlers.NewPlayerHandlers(playerRepo)
    leagueRulesHandler = handlers.NewLeagueRulesHandlers(leagueRulesRepo)
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Received request: %s %s", req.HTTPMethod, req.Path)

    switch req.Path {
    case "/atbats":
        if req.HTTPMethod == "POST" {
            return atBatHandler.RecordAtBat(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return atBatHandler.GetAtBats(ctx, req)
        }

    case "/games":
        if req.HTTPMethod == "POST" {
            return gameHandler.CreateGame(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return gameHandler.GetGame(ctx, req)
        }

    case "/teams":
        if req.HTTPMethod == "POST" {
            return teamHandler.CreateTeam(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return teamHandler.GetTeams(ctx, req)
        }

    case "/leagues":
        if req.HTTPMethod == "POST" {
            return leagueHandler.CreateLeague(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return leagueHandler.GetLeague(ctx, req)
        }

    case "/leaguerules":
        if req.HTTPMethod == "POST" {
            return leagueRulesHandler.CreateLeagueRules(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return leagueRulesHandler.GetLeagueRules(ctx, req)
        }

    case "/players":
        if req.HTTPMethod == "POST" {
            return playerHandler.CreatePlayer(ctx, req)
        }
        if req.HTTPMethod == "GET" {
            return playerHandler.GetPlayers(ctx, req)
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