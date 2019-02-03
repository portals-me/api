package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	collection "../collection/lib"
	. "./lib"
	. "./verifier"
)

type fakeCustomProvider struct {
	cognitoidentityiface.CognitoIdentityAPI
}

func (provider *fakeCustomProvider) GetIdpID(logins Logins) (string, error) {
	return "fake-idp", nil
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
	idp := &fakeCustomProvider{}
	ddb := &fakeDynamoDB{}
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

	callStackIndex := 0

	// user existence check
	if ddb.callStack[callStackIndex].request != "PutItemRequest" {
		t.Error("Invalid callStack order: ", ddb.callStack[callStackIndex])
	}

	user := ParseUser(ddb.callStack[callStackIndex].argument.(*dynamodb.PutItemInput).Item)
	if !(user.ID == "user##fake-idp" &&
		user.Name == testUser.Name &&
		user.DisplayName == testUser.DisplayName &&
		user.Picture == testUser.Picture) {
		t.Errorf("Argument does not match: %+v", ddb.callStack[callStackIndex].argument)
	}

	callStackIndex = 1

	// user collection check
	if ddb.callStack[callStackIndex].request != "GetItemRequest" {
		t.Error("Invalid callStack order: ", ddb.callStack[callStackIndex])
	}

	getItemInput := ddb.callStack[callStackIndex].argument.(*dynamodb.GetItemInput)
	if !(*getItemInput.Key["id"].S == "collection##"+testUser.Name) {
		t.Errorf("Argument does not match: %+v", ddb.callStack[callStackIndex].argument)
	}

	callStackIndex = 2

	// user collection put check
	if ddb.callStack[callStackIndex].request != "PutItemRequest" {
		t.Error("Invalid callStack order: ", ddb.callStack[callStackIndex])
	}

	col := collection.ParseCollection(ddb.callStack[callStackIndex].argument.(*dynamodb.PutItemInput).Item)
	if !(col.ID == "collection#"+testUser.Name) {
		t.Errorf("Argument does not match: %+v", ddb.callStack[callStackIndex].argument)
	}
}

func TestCanSignUpWithTwitter(t *testing.T) {
	idp := &fakeCustomProvider{}
	ddb := &fakeDynamoDB{}
	signer := &fakeSigner{}

	input := SignUpInput{
		Form: struct {
			Name        string `json:"name"`
			DisplayName string `json:"display_name"`
			Picture     string `json:"picture"`
		}{
			Name:        "test",
			DisplayName: "test",
			Picture:     "test",
		},
		Logins: Logins{
			Twitter: "id_token",
		},
	}
	_, identityID, err := DoSignUp(input, idp, ddb, signer)

	if err != nil {
		t.Error("Error", err)
	}
	if identityID != "fake-idp" {
		t.Error("Error", err)
	}
}
