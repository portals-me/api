package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	Dynamo := dynamodb.New(cfg)

	if event.HTTPMethod == "GET" {
		if event.PathParameters == nil {
			result, err := Dynamo.QueryRequest(&dynamodb.QueryInput{
				TableName:              aws.String(os.Getenv("EntityTable")),
				IndexName:              aws.String("owner"),
				KeyConditionExpression: aws.String("owned_by = :owned_by and begins_with(id, :id)"),
				ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
					":owned_by": {
						S: aws.String(user["id"].(string)),
					},
					":id": {
						S: aws.String("collection"),
					},
				},
			}).Send()

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%+v", result.Items), StatusCode: 200}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%+v", event.PathParameters), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
