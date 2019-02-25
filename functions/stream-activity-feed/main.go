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

		event, err := feed.FeedEventFromDynamoEvent(record)
		if err != nil {
			return err
		}

		if strings.HasPrefix(event.EventName, "INSERT") {
			insertItems = append(insertItems, event)
		} else if strings.HasPrefix(event.EventName, "DELETE") {
			var events []feed.FeedEvent
			err := table.
				Get("item_id", event.ItemID).
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

		records = append(records, dbEvent)
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
