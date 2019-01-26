package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/GoogleIdTokenVerifier/GoogleIdTokenVerifier"

	"github.com/gomodule/oauth1/oauth"
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

func sign(payload []byte, keyEncoded string) ([]byte, error) {
	block, _ := pem.Decode([]byte(keyEncoded))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	es256 := jwt.NewES256(privateKey, &privateKey.PublicKey)

	if err != nil {
		return []byte{}, err
	}

	header, _ := json.Marshal(map[string]string{
		"alg": "ES256",
		"typ": "JWT",
	})

	headerEnc := base64.StdEncoding.EncodeToString(header)
	payloadEnc := base64.StdEncoding.EncodeToString(payload)
	signed, err := es256.Sign([]byte(headerEnc + "." + payloadEnc))

	if err != nil {
		return []byte{}, err
	}

	return signed, err
}

func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func GetTwitterClient() *oauth.Client {
	twitterKey := strings.Split(os.Getenv("TwitterKey"), ".")

	return &oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		Credentials: oauth.Credentials{
			Token:  twitterKey[0],
			Secret: twitterKey[1],
		},
	}
}

func GetAccessToken(cred *oauth.Credentials, oauthVerifier string) (*oauth.Credentials, error) {
	client := GetTwitterClient()
	at, _, err := client.RequestToken(nil, cred, oauthVerifier)

	return at, err
}

type TwitterUser struct {
	ID              string `json:"id_str"`
	ScreenName      string `json:"screen_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

func GetTwitterUser(cred *oauth.Credentials, user *TwitterUser) error {
	client := GetTwitterClient()
	resp, err := client.Get(nil, cred, "https://api.twitter.com/1.1/account/verify_credentials.json", url.Values{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return err
	}

	return nil
}

type SignInInput struct {
	Twitter string `json:"twitter"`
	Google  string `json:"google"`
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
	} else if event.Path == "/auth/signIn" {
		if event.HTTPMethod == "POST" {
			var input SignInInput
			err := json.Unmarshal([]byte(event.Body), &input)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			logins := map[string]string{}
			if input.Google != "" {
				logins["accounts.google.com"] = input.Google
			}
			if input.Twitter != "" {
				twitterKey := strings.Split(input.Twitter, ".")
				twitter := GetTwitterClient()

				fmt.Println(twitterKey)

				tokenCred, _, err := twitter.RequestToken(nil, &oauth.Credentials{
					Token:  twitterKey[0],
					Secret: "",
				}, twitterKey[1])
				if err != nil {
					return events.APIGatewayProxyResponse{}, err
				}

				var account TwitterUser
				GetTwitterUser(tokenCred, &account)

				logins["portals.me"] = "twitter-" + account.ID
			}

			getIDReq, err := idp.GetOpenIdTokenForDeveloperIdentityRequest(&cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
				IdentityPoolId: aws.String(os.Getenv("IdentityPoolId")),
				Logins:         logins,
			}).Send()

			fmt.Printf("%+v", getIDReq)

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			identityID := *getIDReq.IdentityId
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

			token, err := sign(jsn, os.Getenv("JwtPrivate"))
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

			tokenCred, _, err := twitter.RequestToken(nil, &oauth.Credentials{
				Token:  event.QueryStringParameters["oauth_token"],
				Secret: "",
			}, event.QueryStringParameters["oauth_verifier"])
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			var account TwitterUser
			GetTwitterUser(tokenCred, &account)

			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
				},
				Body: "Please wait...",
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
}

func main() {
	lambda.Start(handler)
}
