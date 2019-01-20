package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

type Collection struct {
	ID             string            `json:"id"`
	CommentMembers []string          `json:"comment_members"`
	CommentCount   int               `json:"comment_count"`
	Media          []string          `json:"media"`
	Cover          map[string]string `json:"cover"`
	OwnedBy        string            `json:"owned_by"`
	Title          string            `json:"title"`
	CreatedAt      int64             `json:"created_at"`
	Sort           string            `json:"sort"`
	Description    string            `json:"description"`
}

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

			var collections []Collection
			dynamodbattribute.UnmarshalListOfMaps(result.Items, &collections)
			out, _ := json.Marshal(collections)

			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%+v", event.PathParameters),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		StatusCode: 400,
	}, nil
}

func main() {
	lambda.Start(handler)
}
