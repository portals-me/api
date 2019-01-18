package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
}

func handler(ctx context.Context, event Event) (string, error) {
	fmt.Println("hey!")

	return "testy", nil
}

func main() {
	lambda.Start(handler)
}
