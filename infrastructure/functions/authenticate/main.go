package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/guregu/dynamo"
)

var authTableName = os.Getenv("authTable")
var jwtPrivateKey = os.Getenv("jwtPrivate")

type Record struct {
	ID        string `dynamo:"id"`
	sort      string `dynamo:"sort"`
	checkData string `dynamo:"checkData"`
}

type Input struct {
	authType string      `json:"auth_type"`
	data     interface{} `json:"data"`
}

type JwtPayload struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
}

type AuthMethod interface {
	createJwt(table dynamo.Table) (string, error)
}

type Password struct {
	userName string `json:"user_name"`
	password string `json:"password"`
}

func (password Password) createJwt(table dynamo.Table) (string, error) {
	var record Record
	if err := table.
		Get("sort", "name-pass##"+password.userName).
		Index("auth").
		One(&record); err != nil {
		return "", errors.Wrap(err, "UserName not found")
	}

	hash, _ := HashPassword(password.password)
	fmt.Println(hash)
	if err := VerifyPassword(password.password, record.checkData); err != nil {
		return "", errors.Wrap(err, "Invalid Password")
	}

	payload, err := json.Marshal(JwtPayload{
		ID:       record.ID,
		UserName: password.userName,
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

	if input.authType == "password" {
		var password Password

		data, _ := json.Marshal(input.data)
		if err := json.Unmarshal([]byte(data), &password); err != nil {
			return nil, errors.Wrap(err, "Unmarshal password failed")
		}

		return password, nil
	}

	return nil, errors.New("Unsupported auth_type: " + input.authType)
}

/*	POST /authenticate

	expects Input
	returns String (jwt)
*/
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(authTableName)

	method, err := createAuthMethod(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	jwt, err := method.createJwt(authTable)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: jwt, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
