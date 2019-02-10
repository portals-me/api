package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	authenticator "../authenticator/lib"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestCanCreateAndDelete(t *testing.T) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Error(err)
	}
	cfg.Region = "ap-notheast-1"
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == "dynamodb" {
			return aws.Endpoint{
				URL:           "http://localhost:8000",
				SigningRegion: cfg.Region,
			}, nil
		}

		panic(fmt.Errorf(service, region))
	})

	ddb := dynamodb.New(cfg)
	testUser := authenticator.User{
		ID:          "test-user",
		CreatedAt:   time.Now().Unix(),
		Name:        "test-user-name",
		DisplayName: "test-user-display-name",
		Picture:     "",
	}
	testUserOut, _ := authenticator.DumpUser(testUser)
	_, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      testUserOut,
	}).Send()
	if err != nil {
		t.Error(err.Error())
	}

	testInput := CreateInput{
		Title:       "test-title",
		Description: "hoge",
		Cover:       map[string]string{},
	}
	collectionID, err := doCreate(
		testInput,
		map[string]interface{}{
			"id": testUser.ID,
		},
		ddb,
	)
	if err != nil {
		t.Error(err.Error())
	}

	collection, err := doGet(
		collectionID,
		map[string]interface{}{
			"id": testUser.ID,
		},
		ddb,
	)
	if err != nil {
		t.Error(err.Error())
	}

	if !(collectionID == collection.ID &&
		testInput.Title == collection.Title &&
		testInput.Description == collection.Description) {
		t.Errorf("Invalid collection returned: %+v", collection)
	}

	err = doDelete(
		collectionID,
		map[string]interface{}{
			"id": testUser.ID,
		},
		ddb,
	)
	if err != nil {
		t.Error(err.Error())
	}

	collection, err = doGet(
		collectionID,
		map[string]interface{}{
			"id": testUser.ID,
		},
		ddb,
	)
	if err == nil {
		t.Error(err.Error())
	}
}
