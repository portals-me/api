package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	lambda_handler "github.com/aws/aws-lambda-go/lambda"

	"github.com/gomodule/oauth1/oauth"

	. "github.com/myuon/portals-me/functions/authenticator/api"
	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/signer"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
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

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())

	ddb := dynamo.NewFromIface(dynamodb.New(sess))
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	idp := cognitoidentity.New(sess)

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

			body, identityID, err := DoSignUp(input, customProvider, entityTable, signer)
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

			body, err := DoSignIn(logins, customProvider, entityTable, signer)
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
