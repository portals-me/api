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
	timeline "github.com/myuon/portals-me/functions/stream-timeline-feed/lib"
)

type ExtendedTimelineItem struct {
	timeline.TimelineItem
	UserName        string `json:"user_name"`
	UserDisplayName string `json:"user_display_name"`
}

func DoListTimeline(
	userID string,
	entityTable dynamo.Table,
	timelineTable dynamo.Table,
) ([]ExtendedTimelineItem, error) {
	var items []timeline.TimelineItem
	err := timelineTable.
		Get("id", userID).
		Limit(15).
		All(&items)
	if err != nil {
		return nil, errors.Wrap(err, "getFeedsByUserID failed")
	}

	userCache := map[string]authenticator.User{}
	for _, item := range items {
		if _, ok := userCache[item.UserID]; ok {
			continue
		}

		var userDBO authenticator.UserDBO
		if err := entityTable.
			Get("id", item.UserID).
			Range("sort", dynamo.Equal, "user##detail").
			One(&userDBO); err != nil {
			return nil, errors.Wrap(err, "getUserForTimelineItem failed")
		}
		user := userDBO.FromDBO()

		userCache[user.ID] = user
	}

	exItems := make([]ExtendedTimelineItem, len(items))
	for index, item := range items {
		exItems[index] = ExtendedTimelineItem{
			TimelineItem:    item,
			UserName:        userCache[item.UserID].Name,
			UserDisplayName: userCache[item.UserID].DisplayName,
		}
	}

	return exItems, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer

	db := dynamo.New(session.New(), &aws.Config{})
	timelineTable := db.Table(os.Getenv("TimelineTable"))
	entityTable := db.Table(os.Getenv("EntityTable"))

	if event.HTTPMethod == "GET" {
		if event.Resource == "/timeline" {
			items, err := DoListTimeline(user["id"].(string), entityTable, timelineTable)
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
