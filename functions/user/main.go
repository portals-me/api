package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	authenticator "../authenticator/lib"
	feed "../entity-stream/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

func getUser(s map[string]interface{}) (authenticator.User, error) {
	var user authenticator.User
	blob, err := json.Marshal(s)
	if err != nil {
		return authenticator.User{}, err
	}

	err = json.Unmarshal(blob, &user)
	if err != nil {
		return authenticator.User{}, err
	}

	return user, nil
}

func DoListFeed(
	name string,
	entityTable dynamo.Table,
	feedTable dynamo.Table,
) ([]feed.FeedEvent, error) {
	var user authenticator.UserDBO
	err := entityTable.
		Get("sort", "user##detail").
		Range("sort_value", dynamo.Equal, name).
		Index(os.Getenv("SortIndex")).
		One(&user)
	if err != nil {
		return nil, err
	}

	var items []feed.FeedEvent
	err = feedTable.
		Get("user_id", user.ID).
		Limit(10).
		All(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func DoGetUser(
	userName string,
	entityTable dynamo.Table,
) (authenticator.User, error) {
	var user authenticator.UserDBO
	err := entityTable.
		Get("sort", "user##detail").
		Range("sort_value", dynamo.Equal, userName).
		Index(os.Getenv("SortIndex")).
		One(&user)
	if err != nil {
		return authenticator.User{}, err
	}

	return user.FromDBO(), nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName := event.PathParameters["userId"]

	db := dynamo.New(session.New(), &aws.Config{})
	entityTable := db.Table(os.Getenv("EntityTable"))
	feedTable := db.Table(os.Getenv("FeedTable"))

	if event.HTTPMethod == "GET" {
		if event.Resource == "/users/{userId}" {
			if userName == "me" {
				out, _ := json.Marshal(event.RequestContext.Authorizer)
				return events.APIGatewayProxyResponse{Body: string(out), StatusCode: 200}, nil
			}

			user, err := DoGetUser(userName, entityTable)
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       err.Error(),
					StatusCode: 400,
				}, nil
			}

			out, _ := json.Marshal(user)
			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		} else if event.Resource == "/users/{userId}/feed" {
			items, err := DoListFeed(userName, entityTable, feedTable)
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       err.Error(),
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
					StatusCode: 400,
				}, nil
			}

			out, _ := json.Marshal(items)
			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: "Invalid path: " + event.Resource, StatusCode: 400}, nil
}

func main() {
	lambda.Start(handler)
}
