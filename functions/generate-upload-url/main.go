package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-lambda-go/lambda"
)

var storageBucket = os.Getenv("storageBucket")
var ownerCache = make(map[string]map[string]interface{})

func handler(ctx context.Context, event map[string]interface{}) ([]string, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	keys := event["arguments"].(map[string]interface{})["keys"].([]interface{})
	urls := make([]string, len(keys))

	for index, key := range keys {
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(storageBucket),
			Key:    aws.String(key.(string)),
		})
		url, err := req.Presign(time.Hour * 1)

		if err != nil {
			fmt.Println(err.Error())

			return nil, errors.New("Cannot generate URL")
		}

		urls[index] = url
	}

	return urls, nil
}

func main() {
	lambda.Start(handler)
}
