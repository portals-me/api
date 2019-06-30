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
var ownerCache = make(map[string]user.UserInfo)

func handler(ctx context.Context, event []map[string]interface{}) ([]map[string]interface{}, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(accountTable)

	for _, item := range event {
		var record user.UserInfo

		if owner, ok := ownerCache[item["owner"].(string)]; ok {
			record = owner
		} else {
			if err := authTable.
				Get("id", item["owner"]).
				Range("sort", dynamo.Equal, "detail").
				One(&record); err != nil {
				return nil, err
			}

			ownerCache[item["owner"].(string)] = record
		}

		item["owner_name"] = record.Name
		item["owner_display_name"] = record.DisplayName
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
