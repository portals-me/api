package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg, _ := external.LoadDefaultAWSConfig(aws.Config{
		Region: "ap-northeast-1",
	})

	ddb := dynamodb.New(cfg)

	sourceTable := ""
	targetTable := ""

	req := ddb.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String(sourceTable),
	})
	pager := req.Paginate()

	for pager.Next() {
		page := pager.CurrentPage()

		var requests []dynamodb.WriteRequest

		for _, item := range page.Items {
			requests = append(requests, dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: item,
				},
			})
		}

		_, err := ddb.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]dynamodb.WriteRequest{
				targetTable: requests,
			},
		}).Send()

		if err != nil {
			panic(err)
		}
	}
}
