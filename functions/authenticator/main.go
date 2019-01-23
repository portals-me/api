package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/GoogleIdTokenVerifier/GoogleIdTokenVerifier"
)

type SignUpInput struct {
	GoogleToken string `json:"google_token"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

type User struct {
	ID          string `json:"id"`
	Sort        string `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	ddb := dynamodb.New(cfg)
	idp := cognitoidentity.New(cfg)

	if event.Path == "/auth/signUp" {
		if event.HTTPMethod == "POST" {
			var input SignUpInput
			err := json.Unmarshal([]byte(event.Body), input)

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			result, err := idp.GetIdRequest(&cognitoidentity.GetIdInput{
				IdentityPoolId: aws.String(os.Getenv("IdentityPoolId")),
				Logins: map[string]string{
					"accounts.google.com": input.GoogleToken,
				},
			}).Send()
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			identityID := *result.IdentityId
			tokenInfo := GoogleIdTokenVerifier.Verify(input.GoogleToken, os.Getenv("GClientId"))

			if tokenInfo == nil {
				return events.APIGatewayProxyResponse{}, errors.New("Invalid GoogleToken")
			}

			item, err := dynamodbattribute.MarshalMap(User{
				ID:          "user##" + identityID,
				Sort:        "detail",
				CreatedAt:   time.Now().Unix(),
				Name:        input.Name,
				DisplayName: input.DisplayName,
				Picture:     input.Picture,
			})
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			_, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
				TableName:           aws.String(os.Getenv("EntityTable")),
				Item:                item,
				ConditionExpression: aws.String("attribute_not_exists(id)"),
			}).Send()

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 201,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
					"Location":                    "/users/" + identityID,
				},
			}, nil
		}
	}
	if event.HTTPMethod == "POST" {
		if event.PathParameters["userId"] == "me" {
			out, _ := json.Marshal(event.RequestContext.Authorizer)
			return events.APIGatewayProxyResponse{Body: string(out), StatusCode: 200}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
}

func main() {
	lambda.Start(handler)
}
