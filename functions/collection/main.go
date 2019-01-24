package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"

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
				FilterExpression:       aws.String("sort = :sort"),
				ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
					":owned_by": {S: aws.String(user["id"].(string))},
					":id":       {S: aws.String("collection")},
					":sort":     {S: aws.String("detail")},
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
		} else {
			result, err := Dynamo.GetItemRequest(&dynamodb.GetItemInput{
				TableName: aws.String(os.Getenv("EntityTable")),
				Key: map[string]dynamodb.AttributeValue{
					"id":   {S: aws.String("collection##" + event.PathParameters["collectionId"])},
					"sort": {S: aws.String("detail")},
				},
			}).Send()

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			var collection Collection
			dynamodbattribute.UnmarshalMap(result.Item, &collection)
			out, _ := json.Marshal(collection)

			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		}
	} else if event.HTTPMethod == "POST" {
		var createInput map[string]interface{}
		json.Unmarshal([]byte(event.Body), &createInput)

		collectionID := uuid.Must(uuid.NewV4()).String()

		// care for Cover struct
		// You cannot cast map[string]interface{} as map[string]string...
		coverMap := map[string]string{}
		cover := createInput["cover"].(map[string]interface{})
		for key, value := range cover {
			coverMap[key] = value.(string)
		}

		collection, err := dynamodbattribute.MarshalMap(Collection{
			ID:             "collection##" + collectionID,
			Sort:           "detail",
			OwnedBy:        user["id"].(string),
			Title:          createInput["title"].(string),
			Description:    createInput["description"].(string),
			Cover:          coverMap,
			Media:          []string{},
			CommentMembers: []string{user["id"].(string)},
			CommentCount:   0,
			CreatedAt:      time.Now().Unix(),
		})

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		_, err = Dynamo.PutItemRequest(&dynamodb.PutItemInput{
			TableName:           aws.String(os.Getenv("EntityTable")),
			Item:                collection,
			ConditionExpression: aws.String("owned_by = :user_id"),
			ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
				":user_id": {S: aws.String(user["id"].(string))},
			},
		}).Send()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Location":                    "/collections/" + collectionID,
			},
			StatusCode: 201,
		}, nil
	} else if event.HTTPMethod == "DELETE" {
		result, err := Dynamo.QueryRequest(&dynamodb.QueryInput{
			TableName:              aws.String(os.Getenv("EntityTable")),
			KeyConditionExpression: aws.String("id = :id and owned_by = :user_id"),
			ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
				":id":      {S: aws.String("collection##" + event.PathParameters["collectionId"])},
				":user_id": {S: aws.String(user["id"].(string))},
			},
			ProjectionExpression: aws.String("sort"),
		}).Send()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		writeRequest := []dynamodb.WriteRequest{}
		for _, item := range result.Items {
			writeRequest = append(writeRequest, dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]dynamodb.AttributeValue{
						"id":   {S: aws.String("collection##" + event.PathParameters["collectionId"])},
						"sort": item["sort"],
					},
				},
			})
		}

		requestItems := map[string][]dynamodb.WriteRequest{}
		requestItems[os.Getenv("EntityTable")] = writeRequest
		_, err = Dynamo.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		}).Send()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
			StatusCode: 204,
		}, nil
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
