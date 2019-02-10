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

var db = dynamo.New(session.New(), &aws.Config{})
var table = db.Table(os.Getenv("FeedTable"))

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	items := []interface{}{}

	for _, record := range event.Records {
		fmt.Printf("%+v\n", record)

		if record.EventName == "INSERT" && strings.HasPrefix(record.Change.Keys["id"].String(), "collection##") {
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
			items = append(items, feed)
		}
	}

	_, err := table.Batch().Write().Put(items...).Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
