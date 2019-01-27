package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"

	"github.com/aws/aws-lambda-go/events"
)

type fakeCognitoIdentity struct {
	cognitoidentityiface.CognitoIdentityAPI
	payload map[string]string
	err     error
}

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	payload map[string]string
	err     error
}

func TestSignUpWithGoogle(t *testing.T) {
	idp := &fakeCognitoIdentity{}
	ddb := &fakeDynamoDB{}

	jsn, _ := json.Marshal(SignUpInput{
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
	})
	req := events.APIGatewayProxyRequest{
		Body: string(jsn),
	}

	resp, err := DoSignUp(req, idp, ddb)
	if err != nil {
		t.Error("Error", err)
	}
	if resp.StatusCode != 200 {
		t.Error("StatusCode != 200", resp)
	}
}
