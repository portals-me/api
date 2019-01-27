package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gomodule/oauth1/oauth"

	. "./signer"
	. "./verifier"
)

type SignUpInput struct {
	Form struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Picture     string `json:"picture"`
	} `json:"form"`
	Logins Logins `json:"logins"`
}

type User struct {
	ID          string `json:"id"`
	Sort        string `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

type JwtPayload struct {
	ID          string `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func (user User) ToJwtPayload() JwtPayload {
	return JwtPayload{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
	}
}

func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func GetAccessToken(cred *oauth.Credentials, oauthVerifier string) (*oauth.Credentials, error) {
	client := GetTwitterClient()
	at, _, err := client.RequestToken(nil, cred, oauthVerifier)

	return at, err
}
func DoSignUp(
	event events.APIGatewayProxyRequest,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (events.APIGatewayProxyResponse, error) {
	var input SignUpInput
	err := json.Unmarshal([]byte(event.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	identityID, err := idp.GetIdpID(input.Logins)

	item, err := dynamodbattribute.MarshalMap(User{
		ID:          "user##" + identityID,
		Sort:        "detail",
		CreatedAt:   time.Now().Unix(),
		Name:        input.Form.Name,
		DisplayName: input.Form.DisplayName,
		Picture:     input.Form.Picture,
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

	var user User
	err = dynamodbattribute.UnmarshalMap(item, &user)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Location":                    "/users/" + identityID,
		},
		Body: string(body),
	}, nil
}

func DoSignIn(
	event events.APIGatewayProxyRequest,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (events.APIGatewayProxyResponse, error) {
	var logins Logins
	err := json.Unmarshal([]byte(event.Body), &logins)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	identityID, err := idp.GetIdpID(logins)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	userID := "user##" + identityID

	getItemReq, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(userID)},
			"sort": {S: aws.String("detail")},
		},
	}).Send()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if getItemReq.Item["id"].S == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			Body: "UserNotExist: " + userID,
		}, nil
	}

	var user User
	err = dynamodbattribute.UnmarshalMap(getItemReq.Item, &user)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	ddb := dynamodb.New(cfg)
	idp := cognitoidentity.New(cfg)

	customProvider := &CustomProvider{
		IdentityPoolID:          os.Getenv("IdentityPoolId"),
		CognitoIdentityInstance: idp,
	}
	signer := &ES256Signer{
		Key: os.Getenv("JwtPrivate"),
	}

	if event.Path == "/auth/signUp" {
		if event.HTTPMethod == "POST" {
			return DoSignUp(event, customProvider, ddb, signer)
		}
	} else if event.Path == "/auth/signIn" {
		if event.HTTPMethod == "POST" {
			return DoSignIn(event, customProvider, ddb, signer)
		}
	} else if event.Path == "/auth/twitter" {
		if event.HTTPMethod == "POST" {
			fmt.Println(event.Headers)
			twitter := GetTwitterClient()
			result, err := twitter.RequestTemporaryCredentials(
				nil,
				event.Headers["Referer"]+"/twitter-callback",
				nil,
			)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			url := twitter.AuthorizationURL(result, nil)

			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: url,
			}, nil
		} else if event.HTTPMethod == "GET" {
			twitter := GetTwitterClient()
			fmt.Println(event.QueryStringParameters)

			tokenCred, _, err := twitter.RequestToken(nil, &oauth.Credentials{
				Token:  event.QueryStringParameters["oauth_token"],
				Secret: "",
			}, event.QueryStringParameters["oauth_verifier"])
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			var account TwitterUser
			GetTwitterUser(tokenCred, &account)
			jsn, _ := json.Marshal(TwitterCallbackOutput{
				Credential: tokenCred.Token + "." + tokenCred.Secret,
				Account:    account,
			})

			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: string(jsn),
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
}

type TwitterCallbackOutput struct {
	Credential string      `json:"credential"`
	Account    TwitterUser `json:"account"`
}

func main() {
	lambda.Start(handler)
}
