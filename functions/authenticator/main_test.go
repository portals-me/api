package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

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

func TestSignUpWithGoogle(t *testing.T) {
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
}
