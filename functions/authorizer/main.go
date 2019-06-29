package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	. "github.com/portals-me/account/functions/authenticate/account"
	. "github.com/portals-me/account/functions/authenticate/lib"
)

var jwtPrivateKey = os.Getenv("jwtPrivate")

func handler(ctx context.Context, event map[string]interface{}) (Account, error) {
	fmt.Printf("%+v\n", event)
	// Why authorization lower-case?
	token, ok := event["request"].(map[string]interface{})["headers"].(map[string]interface{})["authorization"].(string)
	if ok != true {
		return Account{}, errors.New("Invalid Authorization Token")
	}
	raw := strings.TrimPrefix(token, "Bearer ")

	signer := ES256Signer{
		Key: jwtPrivateKey,
	}
	verified, err := signer.Verify([]byte(raw))
	if err != nil {
		return Account{}, err
	}

	var account Account
	if err := json.Unmarshal(verified, &account); err != nil {
		return Account{}, err
	}

	fmt.Printf("%+v\n", account)
	return account, err
}

func main() {
	lambda.Start(handler)
}
