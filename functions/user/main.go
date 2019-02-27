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
	var userDBO authenticator.UserDBO
	if err := entityTable.
		Get("sort", "user##detail").
		Index(os.Getenv("SortIndex")).
		Range("sort_value", dynamo.Equal, name).
		One(&userDBO); err != nil {
		return false, err
	}

	return true, nil
}

func DoUpdateUser(
	userID string,
	user authenticator.User,
	entityTable dynamo.Table,
) error {
	var userRecordDBO authenticator.UserDBO
	if err := entityTable.
		Get("id", userID).
		Range("sort", dynamo.Equal, "user##detail").
		One(&userRecordDBO); err != nil {
		return errors.Wrap(err, "getUserByID failed")
	}
	userRecord := userRecordDBO.FromDBO()

	if user.Name != "" {
		ex, err := nameExists(user.Name, entityTable)
		if ex == true && err == nil {
			return errors.New("Specified user_name already exists")
		}

		userRecord.Name = user.Name
	}
	if user.DisplayName != "" {
		userRecord.DisplayName = user.DisplayName
	}
	if user.Picture != "" {
		userRecord.Picture = user.Picture
	}

	if err := entityTable.Put(userRecord.ToDBO()).Run(); err != nil {
		return errors.Wrap(err, "putUser failed")
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
	} else if event.HTTPMethod == "PATCH" {
		if event.Resource == "/users/{userId}" {
			authorizedUser := event.RequestContext.Authorizer

			var user authenticator.User
			if err := json.Unmarshal([]byte(event.Body), &user); err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			if err := DoUpdateUser(authorizedUser["id"].(string), user, entityTable); err != nil {
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
