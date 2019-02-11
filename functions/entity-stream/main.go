package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/guregu/dynamo"

	feed "./lib"
)

var db *dynamo.DB
var table dynamo.Table

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	insertItems := []interface{}{}
	deleteItems := []dynamo.Keyed{}

	for _, record := range event.Records {
		fmt.Printf("%+v\n", record)

		if record.EventName == "INSERT" &&
			strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
			record.Change.Keys["sort"].String() == "collection##detail" {
			feed := feed.FeedEvent{
				UserID:    record.Change.NewImage["sort_value"].String(),
				Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
				EventName: "INSERT_COLLECTION",
				ItemID:    record.Change.Keys["id"].String(),
				Entity: map[string]interface{}{
					"title":       record.Change.NewImage["title"].String(),
					"description": record.Change.NewImage["description"].String(),
				},
			}

			insertItems = append(insertItems, feed)
		} else if record.EventName == "INSERT" &&
			strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
			strings.HasPrefix(record.Change.Keys["sort"].String(), "article##") {
			feed := feed.FeedEvent{
				UserID:    record.Change.NewImage["sort_value"].String(),
				Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
				EventName: "INSERT_ARTICLE",
				ItemID:    record.Change.Keys["id"].String() + "/" + record.Change.Keys["sort"].String(),
				Entity: map[string]interface{}{
					"title":       record.Change.NewImage["title"].String(),
					"description": record.Change.NewImage["description"].String(),
				},
			}

			insertItems = append(insertItems, feed)
		} else if record.EventName == "REMOVE" &&
			strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
			record.Change.Keys["sort"].String() == "collection##detail" {
			var events []feed.FeedEvent
			err := table.
				Get("item_id", record.Change.Keys["id"].String()).
				Index("ItemID").
				All(&events)
			if err != nil {
				return err
			}

			for _, event := range events {
				deleteItems = append(deleteItems, dynamo.Keys{event.UserID, event.Timestamp})
			}
		}
	}

	if len(insertItems) > 0 {
		_, err := table.Batch().Write().Put(insertItems...).Run()
		if err != nil {
			return err
		}
	}

	if len(deleteItems) > 0 {
		fmt.Printf("%+v\n", deleteItems)
		_, err := table.Batch("user_id", "timestamp").Write().Delete(deleteItems...).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	db = dynamo.New(session.New(), &aws.Config{})
	table = db.Table(os.Getenv("FeedTable"))

	lambda.Start(handler)
}
