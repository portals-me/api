package feed

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type FeedEvent struct {
	UserID    string                 `json:"user_id" dynamo:"user_id"`
	Timestamp int64                  `json:"timestamp" dynamo:"timestamp"`
	EventName string                 `json:"event_name" dynamo:"event_name"`
	ItemID    string                 `json:"item_id" dynamo:"item_id"`
	Entity    map[string]interface{} `json:"entity" dynamo:"entity"`
}

func FeedEventFromDynamoEvent(record events.DynamoDBEventRecord) (FeedEvent, error) {
	if record.EventName == "INSERT" &&
		strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
		record.Change.Keys["sort"].String() == "collection##detail" {

		description := ""
		if !record.Change.NewImage["description"].IsNull() {
			description = record.Change.NewImage["description"].String()
		}

		return FeedEvent{
			UserID:    record.Change.NewImage["sort_value"].String(),
			Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
			EventName: "INSERT_COLLECTION",
			ItemID:    record.Change.Keys["id"].String(),
			Entity: map[string]interface{}{
				"title":       record.Change.NewImage["title"].String(),
				"description": description,
			},
		}, nil
	} else if record.EventName == "INSERT" &&
		strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
		strings.HasPrefix(record.Change.Keys["sort"].String(), "article##") {

		description := ""
		if !record.Change.NewImage["description"].IsNull() {
			description = record.Change.NewImage["description"].String()
		}

		return FeedEvent{
			UserID:    record.Change.NewImage["sort_value"].String(),
			Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
			EventName: "INSERT_ARTICLE",
			ItemID:    record.Change.Keys["id"].String() + "/" + record.Change.Keys["sort"].String(),
			Entity: map[string]interface{}{
				"title":       record.Change.NewImage["title"].String(),
				"description": description,
			},
		}, nil
	} else if record.EventName == "REMOVE" &&
		strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
		record.Change.Keys["sort"].String() == "collection##detail" {

		return FeedEvent{
			UserID:    "",
			Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
			EventName: "DELETE_COLLECTION",
			ItemID:    record.Change.Keys["id"].String(),
			Entity:    nil,
		}, nil
	}

	return FeedEvent{}, fmt.Errorf("Invalid event: %+v", record)
}
