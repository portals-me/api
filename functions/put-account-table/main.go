package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var svc *dynamodb.DynamoDB
var accountTable = os.Getenv("accountTableName")

func AsDynamoDBAttributeValue(old events.DynamoDBAttributeValue) *dynamodb.AttributeValue {
	if old.DataType() == events.DataTypeBoolean {
		return &dynamodb.AttributeValue{
			BOOL: aws.Bool(old.Boolean()),
		}
	}
	if old.DataType() == events.DataTypeString {
		return &dynamodb.AttributeValue{
			S: aws.String(old.String()),
		}
	}
	if old.DataType() == events.DataTypeList {
		list := old.List()
		newList := make([]*dynamodb.AttributeValue, len(list))

		for i, v := range list {
			newList[i] = AsDynamoDBAttributeValue(v)
		}

		return &dynamodb.AttributeValue{
			L: newList,
		}
	}
	if old.DataType() == events.DataTypeMap {
		kv := old.Map()
		newKv := make(map[string]*dynamodb.AttributeValue)

		for k, v := range kv {
			newKv[k] = AsDynamoDBAttributeValue(v)
		}

		return &dynamodb.AttributeValue{
			M: newKv,
		}
	}

	return nil
}

func AsDynamoDBAttributeValues(old map[string]events.DynamoDBAttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	new := map[string]*dynamodb.AttributeValue{}

	for key, value := range old {
		new[key] = AsDynamoDBAttributeValue(value)
	}

	return new, nil
}

func handler(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		message := record.SNS.Message

		var dbEvent events.DynamoDBEventRecord
		if err := json.Unmarshal([]byte(message), &dbEvent); err != nil {
			return err
		}
		item, err := AsDynamoDBAttributeValues(dbEvent.Change.NewImage)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", item)
		if _, err := svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(accountTable),
			Item:      item,
		}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	svc = dynamodb.New(session.New())

	lambda.Start(handler)
}
