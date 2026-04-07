package handlers

import (
	"net/http"

	"github.com/ajarvis3/kickball-go/pkg/responses"
	"github.com/aws/aws-lambda-go/events"
)

// fetchResource runs loader and converts common outcomes into an API response.
// If the loader returns an error, fetchResource returns (nil, response) where
// response is a 500 JSON error. If the loader returns nil (not found), fetchResource
// returns (nil, response) where response is a 404 JSON error. On success it returns
// (value, nil).
func fetchResource[T any](loader func() (*T, error), notFoundMsg string) (*T, *events.APIGatewayProxyResponse) {
	v, err := loader()
	if err != nil {
		r := responses.JsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return nil, &r
	}
	if v == nil {
		r := responses.JsonResponse(http.StatusNotFound, map[string]string{"error": notFoundMsg})
		return nil, &r
	}
	return v, nil
}
