package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/guregu/dynamo"

	. "github.com/myuon/portals-me/functions/authenticate/lib"
)

var authTableName = os.Getenv("authTable")
var jwtPrivateKey = os.Getenv("jwtPrivate")

type Record struct {
	ID        string `dynamo:"id"`
	Sort      string `dynamo:"sort"`
	CheckData string `dynamo:"check_data"`
}

type Input struct {
	AuthType string      `json:"auth_type"`
	Data     interface{} `json:"data"`
}

type JwtPayload struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
}

type AuthMethod interface {
	createJwt(table dynamo.Table) (string, error)
}

type Password struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (password Password) createJwt(table dynamo.Table) (string, error) {
	var record Record
	if err := table.
		Get("sort", "name-pass##"+password.UserName).
		Index("auth").
		One(&record); err != nil {
		return "", errors.New("UserName not found: " + password.UserName)
	}

	if err := VerifyPassword(record.CheckData, password.Password); err != nil {
		return "", errors.Wrap(err, "Invalid Password")
	}

	payload, err := json.Marshal(JwtPayload{
		ID:       record.ID,
		UserName: password.UserName,
	})
	if err != nil {
		panic(err)
	}

	signer := ES256Signer{
		Key: jwtPrivateKey,
	}
	token, err := signer.Sign(payload)
	if err != nil {
		return "", errors.Wrap(err, "sign failed")
	}

	return string(token), nil
}

func createAuthMethod(body string) (AuthMethod, error) {
	var input Input
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		return nil, errors.Wrap(err, "Unmarshal failed")
	}

	if input.AuthType == "password" {
		var password Password

		data, _ := json.Marshal(input.Data)
		if err := json.Unmarshal([]byte(data), &password); err != nil {
			return nil, errors.Wrap(err, "Unmarshal password failed")
		}

		return password, nil
	}

	return nil, errors.New("Unsupported auth_type: " + input.AuthType)
}

func tryDecodeBase64(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}

	return string(decoded)
}

/*	POST /authenticate

	expects Input
	returns String (jwt)
*/
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// try base64 decoding
	body := tryDecodeBase64(request.Body)
	fmt.Println(body)

	method, err := createAuthMethod(body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(authTableName)

	jwt, err := method.createJwt(authTable)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: jwt, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
