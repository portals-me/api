package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/portals-me/account/lib/jwt"
	"github.com/portals-me/account/lib/user"
)

var jwtPrivateKey = os.Getenv("jwtPrivate")

func handler(ctx context.Context, event map[string]interface{}) (user.UserInfo, error) {
	fmt.Printf("%+v\n", event)
	// Why authorization lower-case?
	token, ok := event["request"].(map[string]interface{})["headers"].(map[string]interface{})["authorization"].(string)
	if ok != true {
		return user.UserInfo{}, errors.New("Invalid Authorization Token")
	}
	raw := strings.TrimPrefix(token, "Bearer ")

	signer := jwt.ES256Signer{
		Key: jwtPrivateKey,
	}
	verified, err := signer.Verify([]byte(raw))
	if err != nil {
		return user.UserInfo{}, err
	}

	var account user.UserInfo
	if err := json.Unmarshal(verified, &account); err != nil {
		return user.UserInfo{}, err
	}

	fmt.Printf("%+v\n", account)
	return account, err
}

func main() {
	lambda.Start(handler)
}
