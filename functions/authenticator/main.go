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

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

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

type Logins struct {
	Twitter string `json:"twitter"`
	Google  string `json:"google"`
}

func (logins *Logins) ToLoginsMap() (map[string]string, error) {
	loginsMap := map[string]string{}
	if logins.Google != "" {
		verifiedID, err := GoogleVerifier{Token: logins.Google}.Verify()
		if err != nil {
			return nil, err
		}

		loginsMap["accounts.google.com"] = verifiedID
	}
	if logins.Twitter != "" {
		verifiedID, err := TwitterVerifier{Token: logins.Twitter}.Verify()
		if err != nil {
			return nil, err
		}

		loginsMap["portals.me"] = verifiedID
	}

	return loginsMap, nil
}

func GetIdpIDByLogins(idp cognitoidentityiface.CognitoIdentityAPI, logins Logins) (string, error) {
	loginsMap, err := logins.ToLoginsMap()
	if err != nil {
		return "", err
	}

	getIDReq, err := idp.GetOpenIdTokenForDeveloperIdentityRequest(&cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolId: aws.String(os.Getenv("IdentityPoolId")),
		Logins:         loginsMap,
	}).Send()
	if err != nil {
		return "", err
	}

	identityID := *getIDReq.IdentityId

	return identityID, nil
}

type IVerifier interface {
	Verify() (string, error)
}

type TwitterVerifier struct {
	Token string
}

func (str TwitterVerifier) Verify() (string, error) {
	twitterKey := strings.Split(str.Token, ".")

	var account TwitterUser
	err := GetTwitterUser(&oauth.Credentials{
		Token:  twitterKey[0],
		Secret: twitterKey[1],
	}, &account)

	if err != nil {
		return "", err
	}

	return "twitter-" + account.ID, nil
}

type GoogleVerifier struct {
	Token string
}

func (str GoogleVerifier) Verify() (string, error) {
	tokenInfo := GoogleIdTokenVerifier.Verify(str.Token, os.Getenv("GClientId"))

	if tokenInfo == nil {
		return "", errors.New("Invalid GoogleToken")
	}

	return str.Token, nil
}

type ICustomProvider interface {
	GetIdpID(Logins) (string, error)
}

type CustomProvider struct {
	IdentityPoolID          string
	CognitoIdentityInstance cognitoidentityiface.CognitoIdentityAPI
}

func (provider *CustomProvider) GetIdpID(logins Logins) (string, error) {
	loginsMap := map[string]string{}
	if logins.Google != "" {
		verified, err := GoogleVerifier{Token: logins.Google}.Verify()
		if err != nil {
			return "", err
		}

		loginsMap["accounts.google.com"] = verified
	}
	if logins.Twitter != "" {
		verified, err := TwitterVerifier{Token: logins.Twitter}.Verify()
		if err != nil {
			return "", err
		}

		loginsMap["portals.me"] = verified
	}

	getIDReq, err := provider.CognitoIdentityInstance.GetOpenIdTokenForDeveloperIdentityRequest(&cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolId: aws.String(provider.IdentityPoolID),
		Logins:         loginsMap,
	}).Send()
	if err != nil {
		return "", err
	}

	return *getIDReq.IdentityId, nil
}

func DoSignUp(
	event events.APIGatewayProxyRequest,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
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
			"Location":                    "/users/" + identityID,
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

	if event.Path == "/auth/signUp" {
		if event.HTTPMethod == "POST" {
			customProvider := &CustomProvider{
				IdentityPoolID:          os.Getenv("IdentityPoolId"),
				CognitoIdentityInstance: idp,
			}
			DoSignUp(event, customProvider, ddb)
		}
	} else if event.Path == "/auth/signIn" {
		if event.HTTPMethod == "POST" {
			var logins Logins
			err := json.Unmarshal([]byte(event.Body), &logins)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			identityID, err := GetIdpIDByLogins(idp, logins)
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
