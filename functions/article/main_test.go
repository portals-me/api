package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"
	collection_api "github.com/myuon/portals-me/functions/collection/api"
)

func TestCreate(t *testing.T) {
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

	collectionID, err := collection_api.DoCreate(
		collection_api.CreateInput{
			Title:       "test-title",
			Description: "test-description",
			Cover:       nil,
		},
		map[string]interface{}{
			"id": "test-user",
		},
		ddb,
	)
	if err != nil {
		t.Error(err)
	}

	statusCode, _, err := doCreate(
		collectionID,
		map[string]interface{}{
			"id": "test-user",
		},
		map[string]interface{}{
			"title":       "hoge",
			"description": "description",
			"entity": map[string]interface{}{
				"hoge": "piyo",
			},
		},
		ddb,
	)
	if err != nil {
		t.Error(err)
	}
	if statusCode != 201 {
		t.Errorf("Invalid StatusCode: %v", statusCode)
	}
}
