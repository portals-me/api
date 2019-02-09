package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	authenticator "../authenticator/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

var db = dynamo.New(session.New(), &aws.Config{})
var sortIndexTable = db.Table(os.Getenv("SortIndex"))
var feedTable = db.Table(os.Getenv("FeedTable"))

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

func DoListFeed(name string) ([]interface{}, error) {
	var user authenticator.UserDBO
	err := sortIndexTable.Get("user##detail", name).One(&user)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	err = feedTable.Get(user.ID, nil).Limit(10).All(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName := event.PathParameters["userName"]

	if event.HTTPMethod == "GET" {
		if userName == "me" {
			out, _ := json.Marshal(event.RequestContext.Authorizer)
			return events.APIGatewayProxyResponse{Body: string(out), StatusCode: 200}, nil
		}

		if strings.HasSuffix(event.Resource, "/feed") {
			items, err := DoListFeed(userName)
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

	return events.APIGatewayProxyResponse{Body: "/user", StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
