package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
	user "github.com/myuon/portals-me/functions/user/lib"
)

var SQS sqsiface.SQSAPI
var entityTable dynamo.Table
var timelineTable dynamo.Table

type TimelineItem struct {
	ID string `json:"id" dynamo:"id"`
	feed.FeedEvent
}

func BuildTimelineItem(ownerID string, item feed.FeedEvent) TimelineItem {
	return TimelineItem{
		ID:        ownerID,
		FeedEvent: item,
	}
}

func processEvent(
	event events.DynamoDBEventRecord,
	entityTable dynamo.Table,
	timelineTable dynamo.Table,
) error {
	feedEvent, err := feed.FeedEventFromDynamoEvent(event)
	if err != nil {
		return err
	}

	if strings.HasPrefix(feedEvent.EventName, "INSERT_") {
		var followers []user.UserFollowRecord
		if err := entityTable.
			Get("id", feedEvent.UserID).
			Range("sort", dynamo.BeginsWith, "user##follow-").
			All(&followers); err != nil {
			return err
		}

		items := make([]interface{}, len(followers))
		for index, follower := range followers {
			target := strings.Split(follower.Sort, "user##follow-")[1]
			items[index] = BuildTimelineItem(target, feedEvent)
		}

		if _, err = timelineTable.
			Batch().
			Write().
			Put(items...).
			Run(); err != nil {
			return err
		}
	} else if strings.HasPrefix(feedEvent.EventName, "DELETE_") {
		var items []TimelineItem
		err := timelineTable.
			Get("item_id", feedEvent.ItemID).
			Index("ItemID").
			Project("id", "timestamp").
			All(&items)
		if err != nil {
			return err
		}

		deleteItems := make([]dynamo.Keyed, len(items))
		for index, item := range items {
			deleteItems[index] = dynamo.Keys{item.ID, item.Timestamp}
		}

		if _, err := timelineTable.
			Batch("id", "timestamp").
			Write().
			Delete(deleteItems...).
			Run(); err != nil {
			return err
		}
	}

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var snsEntity events.SNSEntity
		err := json.Unmarshal([]byte(record.Body), &snsEntity)
		if err != nil {
			return err
		}

		var dbEvent events.DynamoDBEventRecord
		err = json.Unmarshal([]byte(snsEntity.Message), &dbEvent)
		if err != nil {
			return err
		}

		processEvent(
			dbEvent,
			entityTable,
			timelineTable,
		)
	}

	return nil
}

func main() {
	SQS = sqs.New(session.New(), &aws.Config{})
	Dynamo := dynamo.New(session.New(), &aws.Config{})
	entityTable = Dynamo.Table(os.Getenv("EntityTable"))
	timelineTable = Dynamo.Table(os.Getenv("TimelineTable"))

	lambda.Start(handler)
}
