package main

import (
	"context"
	"encoding/json"
	"fmt"
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
)

func processEvent(
	records []events.DynamoDBEventRecord,
	table dynamo.Table,
) error {
	insertItems := []interface{}{}
	deleteItems := []dynamo.Keyed{}

	for _, record := range records {
		fmt.Printf("%+v\n", record)

		if record.EventName == "INSERT" &&
			strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
			record.Change.Keys["sort"].String() == "collection##detail" {

			description := ""
			if !record.Change.NewImage["description"].IsNull() {
				description = record.Change.NewImage["description"].String()
			}

			feed := feed.FeedEvent{
				UserID:    record.Change.NewImage["sort_value"].String(),
				Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
				EventName: "INSERT_COLLECTION",
				ItemID:    record.Change.Keys["id"].String(),
				Entity: map[string]interface{}{
					"title":       record.Change.NewImage["title"].String(),
					"description": description,
				},
			}

			insertItems = append(insertItems, feed)
		} else if record.EventName == "INSERT" &&
			strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") &&
			strings.HasPrefix(record.Change.Keys["sort"].String(), "article##") {

			description := ""
			if !record.Change.NewImage["description"].IsNull() {
				description = record.Change.NewImage["description"].String()
			}

			feed := feed.FeedEvent{
				UserID:    record.Change.NewImage["sort_value"].String(),
				Timestamp: record.Change.ApproximateCreationDateTime.Unix(),
				EventName: "INSERT_ARTICLE",
				ItemID:    record.Change.Keys["id"].String() + "/" + record.Change.Keys["sort"].String(),
				Entity: map[string]interface{}{
					"title":       record.Change.NewImage["title"].String(),
					"description": description,
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

var SQS sqsiface.SQSAPI
var feedTable dynamo.Table

func handler(ctx context.Context, event events.SQSEvent) error {
	var records []events.DynamoDBEventRecord

	for _, record := range event.Records {
		var eventRecord events.DynamoDBEventRecord
		err := json.Unmarshal([]byte(record.Body), eventRecord)
		if err != nil {
			return err
		}

		records = append(records, eventRecord)
	}

	processEvent(
		records,
		feedTable,
	)

	return nil
}

func main() {
	SQS = sqs.New(session.New(), &aws.Config{})
	Dynamo := dynamo.New(session.New(), &aws.Config{})
	feedTable = Dynamo.Table(os.Getenv("FeedTable"))

	lambda.Start(handler)
}
