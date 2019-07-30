package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var SNS snsiface.SNSAPI
var topicArn = os.Getenv("topicArn")

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		jsn, _ := json.Marshal(record)

		_, err := SNS.Publish(&sns.PublishInput{
			Message:  aws.String(string(jsn)),
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			return errors.Wrapf(err, "SNS publich failed: %+v", record)
		}
	}

	return nil
}

func main() {
	SNS = sns.New(session.New(), &aws.Config{})

	lambda.Start(handler)
}
