package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event interface{}) error {
	fmt.Printf("%+v\n", event)
	return nil
}

func main() {
	lambda.Start(handler)
}
