package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	. "github.com/portals-me/account/functions/authenticate/lib"
)

var jwtPrivateKey = os.Getenv("jwtPrivate")

func handler(ctx context.Context, event map[string]interface{}) ([]byte, error) {
	fmt.Printf("%+v\n", event)
	// Why authorization lower-case?
	token, ok := event["request"].(map[string]interface{})["headers"].(map[string]interface{})["authorization"].(string)
	if ok != true {
		return nil, errors.New("Invalid Authorization Token")
	}
	raw := strings.TrimPrefix(token, "Bearer ")

	signer := ES256Signer{
		Key: jwtPrivateKey,
	}
	verified, err := signer.Verify([]byte(raw))
	if err != nil {
		return nil, err
	}

	var payload JwtPayload
	if err := json.Unmarshal(verified, &payload); err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", payload)
	return verified, err
}

func main() {
	lambda.Start(handler)
}
