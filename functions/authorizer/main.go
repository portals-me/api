package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt"
)

func generatePolicy(principalID, effect, resource string, context map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authReponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authReponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authReponse.Context = context

	return authReponse
}

func verify(token string, keyEncoded string) (string, error) {
	block, _ := pem.Decode([]byte(keyEncoded))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	es256 := jwt.NewES256(privateKey, &privateKey.PublicKey)

	decoded, sig, err := jwt.Parse(token)

	if err != nil {
		return "", err
	}
	if err = es256.Verify(decoded, sig); err != nil {
		return "", err
	}

	payloadEncoded := strings.Split(string(decoded), ".")[1]
	decoded64, err := base64.StdEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return "", err
	}

	return string(decoded64), nil
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	if !strings.HasPrefix(event.AuthorizationToken, "Bearer ") {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	token := strings.Replace(event.AuthorizationToken, "Bearer ", "", 1)
	payload, err := verify(token, os.Getenv("JwtPrivate"))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	var data interface{}
	json.Unmarshal([]byte(payload), &data)

	isAllowed := true

	var effect string
	if isAllowed {
		effect = "Allow"
	} else {
		effect = "Deny"
	}
	userID := data.(map[string]interface{})["id"].(string)

	return generatePolicy(userID, effect, event.MethodArn, data.(map[string]interface{})), nil
}

func main() {
	lambda.Start(handler)
}