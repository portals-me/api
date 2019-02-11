package api

import (
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/gofrs/uuid"

	. "../lib"
)

func DoList(
	user map[string]interface{},
	ddb dynamodbiface.DynamoDBAPI,
) ([]Collection, error) {
	result, err := ddb.QueryRequest(&dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("EntityTable")),
		IndexName:              aws.String(os.Getenv("SortIndex")),
		KeyConditionExpression: aws.String("sort = :sort and sort_value = :sort_value"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":sort":       {S: aws.String("collection##detail")},
			":sort_value": {S: aws.String(user["id"].(string))},
		},
	}).Send()

	if err != nil {
		return nil, err
	}

	return ParseCollections(result.Items), nil
}

func DoGet(
	collectionID string,
	user map[string]interface{},
	ddb dynamodbiface.DynamoDBAPI,
) (Collection, error) {
	result, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + collectionID)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()
	if err != nil {
		return Collection{}, err
	}
	if len(result.Item) == 0 {
		return Collection{}, errors.New("Not exist")
	}

	collection := ParseCollection(result.Item)

	result, err = ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(collection.Owner)},
			"sort": {S: aws.String("user##detail")},
		},
	}).Send()

	// First-aid
	collection.Owner = *result.Item["sort_value"].S

	return collection, nil
}

type CreateInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Cover       map[string]string `json:"cover"`
}

func DoCreate(
	createInput CreateInput,
	user map[string]interface{},
	ddb dynamodbiface.DynamoDBAPI,
) (string, error) {
	collectionID := uuid.Must(uuid.NewV4()).String()

	item, err := DumpCollection(Collection{
		ID:             collectionID,
		Owner:          user["id"].(string),
		Title:          createInput.Title,
		Description:    createInput.Description,
		Cover:          createInput.Cover,
		Media:          []string{},
		CommentMembers: []string{user["id"].(string)},
		CommentCount:   0,
		CreatedAt:      time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}

	if _, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      item,
	}).Send(); err != nil {
		return "", err
	}

	return collectionID, nil
}

func DoDelete(
	collectionID string,
	user map[string]interface{},
	ddb dynamodbiface.DynamoDBAPI,
) error {
	result, err := ddb.QueryRequest(&dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("EntityTable")),
		KeyConditionExpression: aws.String("id = :id and sort = :user_id"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":id":      {S: aws.String("collection##" + collectionID)},
			":user_id": {S: aws.String(user["id"].(string))},
		},
		ProjectionExpression: aws.String("sort"),
	}).Send()

	if err != nil {
		return err
	}

	writeRequest := []dynamodb.WriteRequest{
		dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]dynamodb.AttributeValue{
					"id":   {S: aws.String("collection##" + collectionID)},
					"sort": {S: aws.String("collection##detail")},
				},
			},
		},
	}

	for _, item := range result.Items {
		writeRequest = append(writeRequest, dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]dynamodb.AttributeValue{
					"id":   {S: aws.String("collection##" + collectionID)},
					"sort": item["sort"],
				},
			},
		})
	}

	requestItems := map[string][]dynamodb.WriteRequest{}
	requestItems[os.Getenv("EntityTable")] = writeRequest
	_, err = ddb.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
		RequestItems: requestItems,
	}).Send()

	if err != nil {
		return err
	}

	return nil
}
