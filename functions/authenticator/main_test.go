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

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	payload map[string]string
	err     error
}

func (ddb *fakeDynamoDB) PutItemRequest(input *dynamodb.PutItemInput) dynamodb.PutItemRequest {
	return dynamodb.PutItemRequest{
		Request: &aws.Request{
			Data:  &dynamodb.PutItemOutput{},
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
	resp, err := DoSignUp(input, idp, ddb, signer)

	if err != nil {
		t.Error("Error", err)
	}
	if resp.StatusCode != 200 {
		t.Error("StatusCode != 200", resp)
	}
}
