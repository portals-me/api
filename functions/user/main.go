package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/aws/aws-lambda-go/lambda"
	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"
	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
	. "github.com/myuon/portals-me/functions/user/lib"
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
		Index(os.Getenv("SortIndex")).
		Range("sort_value", dynamo.Equal, name).
		One(&user)
	if err != nil {
		return nil, errors.Wrap(err, "getUserByName failed")
	}

	var items []feed.FeedEvent
	err = feedTable.
		Get("user_id", user.ID).
		Limit(10).
		All(&items)
	if err != nil {
		return nil, errors.Wrap(err, "getFeedsByUserID failed")
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

func nameExists(
	name string,
	entityTable dynamo.Table,
) (bool, error) {
	count, err := entityTable.
		Get("sort", "user##detail").
		Index(os.Getenv("SortIndex")).
		Range("sort_value", dynamo.Equal, name).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

type UpdateInput struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func DoUpdateUser(
	userID string,
	input UpdateInput,
	entityTable dynamo.Table,
) error {
	var userDBO authenticator.UserDBO
	if err := entityTable.
		Get("id", userID).
		Range("sort", dynamo.Equal, "user##detail").
		One(&userDBO); err != nil {
		return errors.Wrap(err, "getUserByID failed")
	}
	user := userDBO.FromDBO()

	// user name existance check
	if input.Name != "" && user.Name != input.Name {
		if ex, err := nameExists(input.Name, entityTable); ex == true && err == nil {
			return errors.New("user_name " + input.Name + " already exists")
		}
	}

	// create an update query
	// NB: DynamoDB local not support updateItem onto the sortkey of index
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.DisplayName != "" {
		user.DisplayName = input.DisplayName
	}
	if input.Picture != "" {
		user.Picture = input.Picture
	}

	if err := entityTable.Put(user.ToDBO()).
		Run(); err != nil {
		return err
	}

	return nil
}

func DoFollowUser(
	source string,
	targetName string,
	entityTable dynamo.Table,
) error {
	var targetDBO authenticator.UserDBO
	if err := entityTable.
		Get("sort", "user##detail").
		Range("sort_value", dynamo.Equal, targetName).
		Index(os.Getenv("SortIndex")).
		One(&targetDBO); err != nil {
		return err
	}
	target := targetDBO.FromDBO()

	if source == target.ID {
		return errors.New("Cannot follow oneself")
	}

	if err := entityTable.Put(UserFollowRecord{
		ID:    target.ID,
		Sort:  "user##follow-" + source,
		Value: target.ID,
	}).Run(); err != nil {
		return err
	}

	return nil
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
	} else if event.HTTPMethod == "POST" {
		if event.Resource == "/users/{userId}/follow" {
			user := event.RequestContext.Authorizer

			if err := DoFollowUser(user["id"].(string), userName, entityTable); err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 204,
			}, nil
		}
	} else if event.HTTPMethod == "PUT" {
		if event.Resource == "/users/{userId}" {
			authorizedUser := event.RequestContext.Authorizer

			var input UpdateInput
			if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			if authorizedUser["id"].(string) != "user##ap-northeast-1:"+userName {
				return events.APIGatewayProxyResponse{}, errors.New("Access Denied")
			}

			if err := DoUpdateUser(authorizedUser["id"].(string), input, entityTable); err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 204,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{Body: "Invalid path: " + event.Resource, StatusCode: 400}, nil
}

func main() {
	lambda.Start(handler)
}
