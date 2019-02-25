package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	. "github.com/myuon/portals-me/functions/stream-timeline-feed/api"
)

var SQS sqsiface.SQSAPI
var entityTable dynamo.Table
var timelineTable dynamo.Table

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

		ProcessEvent(
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
