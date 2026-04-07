package responses

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func JsonResponse(status int, body interface{}) events.APIGatewayProxyResponse {
	b, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(b),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
