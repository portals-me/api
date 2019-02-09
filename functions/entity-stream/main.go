package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		fmt.Printf("%+v\n", record)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
