package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer
	out, _ := json.Marshal(user)
	return events.APIGatewayProxyResponse{Body: string(out), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
