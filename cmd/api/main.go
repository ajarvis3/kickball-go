package main

import (
	"context"
	"log"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"net/http"

	"github.com/ajarvis3/kickball-go/internal/db"
	"github.com/ajarvis3/kickball-go/internal/handlers"
	"github.com/ajarvis3/kickball-go/internal/services"
	"github.com/ajarvis3/kickball-go/pkg/responses"
)

var (
	atBatHandler       *handlers.AtBatHandlers
	gameHandler        *handlers.GameHandlers
	teamHandler        *handlers.TeamHandlers
	leagueHandler      *handlers.LeagueHandlers
	playerHandler      *handlers.PlayerHandlers
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

	// Construct services
	rulesEngine := services.NewRulesEngine()
	gameEngine := services.NewGameEngine(rulesEngine)

	// Construct handlers
	atBatHandler = handlers.NewAtBatHandlers(atBatRepo, gameRepo, leagueRulesRepo, gameEngine)
	gameHandler = handlers.NewGameHandlers(gameRepo, leagueRulesRepo)
	teamHandler = handlers.NewTeamHandlers(teamRepo)
	leagueHandler = handlers.NewLeagueHandlers(leagueRepo)
	playerHandler = handlers.NewPlayerHandlers(playerRepo)
	leagueRulesHandler = handlers.NewLeagueRulesHandlers(leagueRulesRepo)
}


func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", req.HTTPMethod, req.Path)

	switch req.Path {
	case "/atbats":
		return handleMethods(ctx, req, atBatHandler.RecordAtBat, atBatHandler.GetAtBats)

	case "/games":
		return handleMethods(ctx, req, gameHandler.CreateGame, gameHandler.GetGame)

	case "/teams":
		return handleMethods(ctx, req, teamHandler.CreateTeam, teamHandler.GetTeams)

	case "/leagues":
		return handleMethods(ctx, req, leagueHandler.CreateLeague, leagueHandler.GetLeague)

	case "/leaguerules":
		return handleMethods(ctx, req, leagueRulesHandler.CreateLeagueRules, leagueRulesHandler.GetLeagueRules)

	case "/players":
		return handleMethods(ctx, req, playerHandler.CreatePlayer, playerHandler.GetPlayers)
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

// methodHandler is the signature used by our route handlers.
type methodHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// handleMethods dispatches to the provided POST or GET handler functions
// based on req.HTTPMethod. If the method is not supported, returns 405.
func handleMethods(ctx context.Context, req events.APIGatewayProxyRequest, post methodHandler, get methodHandler) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return post(ctx, req)
	case "GET":
		return get(ctx, req)
	case "OPTIONS":
		// Respond to CORS preflight requests with 200 and appropriate headers
		return responses.JsonResponse(http.StatusOK, map[string]string{"status": "ok"}), nil
	default:
		return responses.JsonResponse(http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}), nil
	}
}
