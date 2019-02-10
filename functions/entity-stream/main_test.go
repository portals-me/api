package main

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	feed "./lib"
)

func TestSendInsertAndRemove(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	table = db.Table(os.Getenv("FeedTable"))

	err := handler(nil, events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			events.DynamoDBEventRecord{
				EventID:   "1",
				EventName: "INSERT",
				Change: events.DynamoDBStreamRecord{
					ApproximateCreationDateTime: events.SecondsEpochTime{Time: time.Now()},
					Keys: map[string]events.DynamoDBAttributeValue{
						"id": events.NewStringAttribute("collection##aaaa"),
					},
					NewImage: map[string]events.DynamoDBAttributeValue{
						"sort_value":  events.NewStringAttribute("user##u"),
						"title":       events.NewStringAttribute("title"),
						"description": events.NewStringAttribute("description"),
					},
				},
			},
			events.DynamoDBEventRecord{
				EventID:   "5",
				EventName: "INSERT",
				Change: events.DynamoDBStreamRecord{
					ApproximateCreationDateTime: events.SecondsEpochTime{Time: time.Now()},
					Keys: map[string]events.DynamoDBAttributeValue{
						"id": events.NewStringAttribute("collection##aaaa"),
					},
					NewImage: map[string]events.DynamoDBAttributeValue{
						"sort_value":  events.NewStringAttribute("user##k"),
						"title":       events.NewStringAttribute("title"),
						"description": events.NewStringAttribute("description"),
					},
				},
			},
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	var items []feed.FeedEvent
	table.Get("item_id", "collection##aaaa").Index("ItemID").All(&items)
	if len(items) != 2 {
		t.Errorf("Invalid items: %+v", items)
	}

	err = handler(nil, events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			events.DynamoDBEventRecord{
				EventID:   "10",
				EventName: "REMOVE",
				Change: events.DynamoDBStreamRecord{
					ApproximateCreationDateTime: events.SecondsEpochTime{Time: time.Now()},
					Keys: map[string]events.DynamoDBAttributeValue{
						"id": events.NewStringAttribute("collection##aaaa"),
					},
				},
			},
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	items = []feed.FeedEvent{}
	table.Get("item_id", "collection##aaaa").Index("ItemID").All(&items)
	if len(items) != 0 {
		t.Errorf("Invalid items: %+v", items)
	}
}
