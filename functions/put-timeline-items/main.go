package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamo_helper "github.com/portals-me/api/lib/dynamo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var svc *dynamodb.DynamoDB
var tableName = os.Getenv("tableName")

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

			if _, err := svc.PutItem(&dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      item,
			}); err != nil {
				return err
			}
		} else if dbEvent.EventName == "REMOVE" {
			item, err := dynamo_helper.AsDynamoDBAttributeValues(dbEvent.Change.Keys)
			if err != nil {
				return err
			}

			if _, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key:       item,
			}); err != nil {
				return err
			}

			// Skip error check
			svc.DeleteItem(&dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
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

	lambda.Start(handler)
}
