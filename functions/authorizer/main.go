package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID string
}

func (*User) Valid() error {
	return nil
}

func (user *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id": user.ID,
	}
}

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

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	if !strings.HasPrefix(event.AuthorizationToken, "Bearer ") {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	token := strings.Replace(event.AuthorizationToken, "Bearer ", "", 1)

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "ES256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method.Alg())
		}

		return token, nil
	})
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	if err := parsed.Method.Verify("ES256", os.Getenv("JwtPublic"), func(token *jwt.Token) (interface{}, error) {
		return token, nil
	}); err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized")
	}

	isAllowed := true

	var effect string
	if isAllowed {
		effect = "Allow"
	} else {
		effect = "Deny"
	}

	userID := parsed.Claims.(*User).ID

	return generatePolicy(userID, effect, event.MethodArn, parsed.Claims.(*User).ToMap()), nil
}

func main() {
	lambda.Start(handler)
}
