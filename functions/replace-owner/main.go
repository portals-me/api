package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	"github.com/portals-me/account/lib/user"
)

var accountTable = os.Getenv("accountTableName")

func handler(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(accountTable)

	var record user.UserInfo
	if err := authTable.
		Get("id", event["owner"]).
		Range("sort", dynamo.Equal, "detail").
		One(&record); err != nil {
		return nil, err
	}

	event["owner"] = record.Name
	return event, nil
}

func main() {
	lambda.Start(handler)
}
