package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/lambdaiface"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	lambda_handler "github.com/aws/aws-lambda-go/lambda"

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
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

type UserDBO struct {
	ID          string `json:"id"`
	Sort        string `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"sort_value"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func (user User) toDBO() UserDBO {
	return UserDBO{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
		Sort:        "user##detail",
	}
}

func (user UserDBO) fromDBO() User {
	return User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
	}
}

func parseUser(attr map[string]dynamodb.AttributeValue) User {
	var userDBO UserDBO
	dynamodbattribute.UnmarshalMap(attr, &userDBO)

	return userDBO.fromDBO()
}

func dumpUser(user User) (map[string]dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(user.toDBO())
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

type CreateInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Cover       map[string]string `json:"cover"`
}

func createUserCollection(
	user User,
	ddb dynamodbiface.DynamoDBAPI,
	lam lambdaiface.LambdaAPI,
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

	// Call collection handler
	// Isn't there a better way?
	input, _ := json.Marshal(CreateInput{
		Title:       user.Name,
		Description: "",
		Cover: map[string]string{
			"color": "red lighten-3",
			"sort":  "solid",
		},
	})

	userData, _ := json.Marshal(user)
	var authorizer map[string]interface{}
	json.Unmarshal(userData, &authorizer)

	payload, _ := json.Marshal(events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       string(input),
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: authorizer,
		},
	})

	funcResult, err := lam.InvokeRequest(&lambda.InvokeInput{
		FunctionName: aws.String(strings.Replace(os.Getenv("LAMBDA_FUNCTION_NAME"), "authenticator", "collection", 1)),
		Payload:      payload,
	}).Send()
	fmt.Println(funcResult)

	return err
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

	item, err := dumpUser(User{
		ID:          "user##" + identityID,
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
	lam lambdaiface.LambdaAPI,
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
			"sort": {S: aws.String("user##detail")},
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

	user := parseUser(getItemReq.Item)

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

	err = createUserCollection(user, ddb, lam)
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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	ddb := dynamodb.New(cfg)
	idp := cognitoidentity.New(cfg)
	lam := lambda.New(cfg)

	customProvider := &CustomProvider{
		IdentityPoolID:          os.Getenv("IdentityPoolId"),
		CognitoIdentityInstance: idp,
	}
	signer := &ES256Signer{
		Key: os.Getenv("JwtPrivate"),
	}

	if event.Resource == "/auth/signUp" {
		if event.HTTPMethod == "POST" {
			return DoSignUp(event, customProvider, ddb, signer)
		}
	} else if event.Resource == "/auth/signIn" {
		if event.HTTPMethod == "POST" {
			return DoSignIn(event, customProvider, ddb, lam, signer)
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
