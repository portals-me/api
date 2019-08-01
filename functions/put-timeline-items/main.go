package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofrs/uuid"
	dynamo_helper "github.com/portals-me/api/lib/dynamo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/portals-me/api/lib/user"
)

var timelineTableName = os.Getenv("timelineTableName")
var userTableName = os.Getenv("userTableName")
var svc *dynamodb.DynamoDB
var userRepository user.Repository

func createNotifiedItemID(itemID string, followerID string) string {
	return followerID + "-" + itemID
}

func createBatchRequestsToFollowers(item map[string]*dynamodb.AttributeValue) (*dynamodb.BatchWriteItemInput, error) {
	followers, err := userRepository.ListFollowers(item["owner"].String())
	if err != nil {
		return nil, err
	}

	// Move id to original_id and fill new uuid in id
	item["original_id"] = item["id"]
	item["id"] = &dynamodb.AttributeValue{S: aws.String(uuid.Must(uuid.NewV4()).String())}

	var items map[string][]*dynamodb.WriteRequest
	var requests []*dynamodb.WriteRequest = make([]*dynamodb.WriteRequest, len(followers))

	for index, followerID := range followers {
		// target modification
		item["target"] = &dynamodb.AttributeValue{S: aws.String(followerID)}

		requests[index] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		}
	}

	// post to owner themselves
	item["target"] = &dynamodb.AttributeValue{S: aws.String(item["owner"].String())}
	requests = append(requests, &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{
			Item: item,
		},
	})

	items[timelineTableName] = requests

	return &dynamodb.BatchWriteItemInput{
		RequestItems: items,
	}, nil
}

func handler(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		message := record.SNS.Message

		var dbEvent events.DynamoDBEventRecord
		if err := json.Unmarshal([]byte(message), &dbEvent); err != nil {
			return err
		}

		if dbEvent.EventName == "MODIFY" || dbEvent.EventName == "INSERT" {
			item, err := dynamo_helper.AsDynamoDBAttributeValues(dbEvent.Change.NewImage)
			if err != nil {
				return err
			}

			req, err := createBatchRequestsToFollowers(item)
			if err != nil {
				return err
			}

			if _, err := svc.BatchWriteItem(req); err != nil {
				return err
			}
		} else if dbEvent.EventName == "REMOVE" {
			item, err := dynamo_helper.AsDynamoDBAttributeValues(dbEvent.Change.Keys)
			if err != nil {
				return err
			}

			if _, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
				TableName: aws.String(timelineTableName),
				Key:       item,
			}); err != nil {
				return err
			}

			// Skip error check
			svc.DeleteItem(&dynamodb.DeleteItemInput{
				TableName: aws.String(timelineTableName),
				Key: map[string]*dynamodb.AttributeValue{
					"id":   item["id"],
					"sort": &dynamodb.AttributeValue{S: aws.String("social")},
				},
				ConditionExpression: aws.String("attribute_exists(id)"),
			})
		} else {
			fmt.Printf("%+v\n", dbEvent)
			panic("Not supported EventName: " + dbEvent.EventName)
		}
	}

	return nil
}

func main() {
	svc = dynamodb.New(session.New())
	userRepository = user.NewRepositoryFromAWS(svc, userTableName)

	lambda.Start(handler)
}
