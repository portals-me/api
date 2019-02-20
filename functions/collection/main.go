package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"

	. "github.com/myuon/portals-me/functions/collection/api"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer
	sess := session.Must(session.NewSession())

	ddb := dynamo.NewFromIface(dynamodb.New(sess))
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	if event.HTTPMethod == "GET" {
		if event.PathParameters == nil {
			collections, err := DoList(user["id"].(string), entityTable)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			out, _ := json.Marshal(collections)
			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		} else {
			collection, err := DoGet(event.PathParameters["collectionId"], entityTable)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			out, _ := json.Marshal(collection)
			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		}
	} else if event.HTTPMethod == "POST" {
		var createInput CreateInput
		json.Unmarshal([]byte(event.Body), &createInput)

		collectionID, err := DoCreate(createInput, user["id"].(string), entityTable)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Location":                    "/collections/" + collectionID,
			},
			StatusCode: 201,
		}, nil
	} else if event.HTTPMethod == "DELETE" {
		err := DoDelete(event.PathParameters["collectionId"], user["id"].(string), entityTable)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
			StatusCode: 204,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%+v", event.PathParameters),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		StatusCode: 400,
	}, nil
}

func main() {
	lambda.Start(handler)
}
