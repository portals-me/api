package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	. "github.com/myuon/portals-me/functions/authenticator/api"
	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
	collection "github.com/myuon/portals-me/functions/collection/lib"
)

type fakeCustomProvider struct {
	cognitoidentityiface.CognitoIdentityAPI
	customID string
}

func (provider *fakeCustomProvider) GetIdpID(logins Logins) (string, error) {
	if provider.customID != "" {
		return provider.customID, nil
	} else {
		return "fake-idp", nil
	}
}

type operation struct {
	request  string
	argument interface{}
}

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	callStack []operation
	payload   map[string]map[string]dynamodb.AttributeValue
	err       error
}

func (ddb *fakeDynamoDB) PutItemRequest(input *dynamodb.PutItemInput) dynamodb.PutItemRequest {
	ddb.callStack = append(ddb.callStack, operation{
		request:  "PutItemRequest",
		argument: input,
	})

	ddb.payload[*input.Item["id"].S+"-"+*input.Item["sort"].S] = input.Item

	return dynamodb.PutItemRequest{
		Request: &aws.Request{
			Data:  &dynamodb.PutItemOutput{},
			Error: ddb.err,
		},
	}
}

func (ddb *fakeDynamoDB) GetItemRequest(input *dynamodb.GetItemInput) dynamodb.GetItemRequest {
	ddb.callStack = append(ddb.callStack, operation{
		request:  "GetItemRequest",
		argument: input,
	})

	return dynamodb.GetItemRequest{
		Request: &aws.Request{
			Data: &dynamodb.GetItemOutput{
				Item: ddb.payload[*input.Key["id"].S+"-"+*input.Key["sort"].S],
			},
			Error: ddb.err,
		},
	}
}

type fakeSigner struct {
}

func (signer fakeSigner) Sign(payload []byte) ([]byte, error) {
	return payload, nil
}

func TestCanSignUpWithGoogle(t *testing.T) {
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
	idp := &fakeCustomProvider{}
	signer := &fakeSigner{}

	testUser := struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Picture     string `json:"picture"`
	}{
		Name:        "test_name",
		DisplayName: "test_display_name",
		Picture:     "test_picture",
	}
	input := SignUpInput{
		Form: testUser,
		Logins: Logins{
			Google: "id_token",
		},
	}
	_, identityID, err := DoSignUp(input, idp, ddb, signer)

	if err != nil {
		t.Error("Error", err)
	}
	if identityID != "fake-idp" {
		t.Error("Error", err)
	}

	item, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("user##" + identityID)},
			"sort": {S: aws.String("user##detail")},
		},
	}).Send()
	if err != nil {
		t.Error("error", err)
	}

	if !(*item.Item["sort_value"].S == testUser.Name &&
		*item.Item["display_name"].S == testUser.DisplayName &&
		*item.Item["picture"].S == testUser.Picture) {
		t.Errorf("Argument does not match: %+v", item.Item)
	}

	item, err = ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + testUser.Name)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()
	if err != nil {
		t.Error("error", err)
	}

	if !(*item.Item["title"].S == testUser.Name) {
		t.Errorf("Argument does not match: %+v", item.Item)
	}
}

func TestCanSignInWithoutUserCollectionTwice(t *testing.T) {
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

	testUser := User{
		ID:   "user##user-id",
		Name: "user-name",
	}
	testUserDump, _ := DumpUser(testUser)

	_, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      testUserDump,
	}).Send()
	if err != nil {
		t.Error("error", err)
	}

	idp := &fakeCustomProvider{
		customID: "user-id",
	}
	signer := &fakeSigner{}
	logins := Logins{
		Twitter: "id_token",
	}

	_, err = DoSignIn(logins, idp, ddb, signer)
	if err != nil {
		t.Error("Error", err)
	}

	item, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + testUser.Name)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()
	if err != nil {
		t.Error("error", err)
	}

	col := collection.ParseCollection(item.Item)
	if !(col.ID == testUser.Name &&
		col.Title == testUser.Name) {
		t.Errorf("Argument does not match: %+v", item.Item)
	}

	_, err = DoSignIn(logins, idp, ddb, signer)
	if err != nil {
		t.Error("Error", err)
	}
}
