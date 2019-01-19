package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func handler(ctx context.Context, event Event) (Response, error) {
	fmt.Println("hey!")

	return Response{StatusCode: 200, Body: "body!"}, nil
}

func main() {
	lambda.Start(handler)
}
