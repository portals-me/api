package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	lambda_handler "github.com/aws/aws-lambda-go/lambda"

	"github.com/gomodule/oauth1/oauth"

	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/signer"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
	collection "github.com/myuon/portals-me/functions/collection/lib"
)

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

type CreateInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Cover       map[string]string `json:"cover"`
}

func createUserCollection(
	user User,
	ddb dynamodbiface.DynamoDBAPI,
) error {
	result, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + user.Name)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()

	if err != nil {
		return err
	}

	if result.Item != nil {
		return nil
	}

	item, err := collection.DumpCollection(collection.Collection{
		ID:          user.Name,
		Owner:       user.ID,
		Title:       user.Name,
		Description: "",
		Cover: map[string]string{
			"color": "red lighten-3",
			"sort":  "solid",
		},
		Media:          []string{},
		CommentMembers: []string{user.ID},
		CommentCount:   0,
		CreatedAt:      time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	if _, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      item,
	}).Send(); err != nil {
		return err
	}

	return err
}

func DoSignUp(
	input SignUpInput,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (string, string, error) {
	identityID, err := idp.GetIdpID(input.Logins)

	user := User{
		ID:          "user##" + identityID,
		CreatedAt:   time.Now().Unix(),
		Name:        input.Form.Name,
		DisplayName: input.Form.DisplayName,
		Picture:     input.Form.Picture,
	}

	item, err := DumpUser(user)
	if err != nil {
		return "", "", err
	}

	if _, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("EntityTable")),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}).Send(); err != nil {
		return "", "", err
	}

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return "", "", err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return "", "", err
	}

	if err = createUserCollection(user, ddb); err != nil {
		return "", "", err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	return string(body), identityID, nil
}

func DoSignIn(
	logins Logins,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (string, error) {
	identityID, err := idp.GetIdpID(logins)
	if err != nil {
		return "", err
	}

	userID := "user##" + identityID

	getItemReq, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(userID)},
			"sort": {S: aws.String("user##detail")},
		},
	}).Send()
	if err != nil {
		return "", err
	}

	if getItemReq.Item["id"].S == nil {
		return "", errors.New("UserNotExist: " + userID)
	}

	user := ParseUser(getItemReq.Item)

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return "", err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	err = createUserCollection(user, ddb)
	if err != nil {
		return "", err
	}

	return string(body), nil
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

	if event.Resource == "/auth/signUp" {
		if event.HTTPMethod == "POST" {
			var input SignUpInput
			err := json.Unmarshal([]byte(event.Body), &input)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			body, identityID, err := DoSignUp(input, customProvider, ddb, signer)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
					"Location":                    "/users/" + identityID,
				},
				Body: string(body),
			}, nil
		}
	} else if event.Resource == "/auth/signIn" {
		if event.HTTPMethod == "POST" {
			var logins Logins
			err := json.Unmarshal([]byte(event.Body), &logins)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			body, err := DoSignIn(logins, customProvider, ddb, signer)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: string(body),
			}, nil
		}
	} else if event.Resource == "/auth/twitter" {
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
	lambda_handler.Start(handler)
}
