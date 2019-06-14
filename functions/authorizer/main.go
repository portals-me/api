package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event map[string]interface{}) (string, error) {
	fmt.Printf("event: %+v\n", event)
	fmt.Printf("arguments: %+v\n", event["arguments"])
	fmt.Printf("header: %+v\n", event["request"].(map[string]interface{})["headers"])
	return "", errors.New("Constant error message")
}

func main() {
	lambda.Start(handler)
}
