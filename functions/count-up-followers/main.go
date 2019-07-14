package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	//"github.com/portals-me/account/lib/user"
	//	userExtra "github.com/portals-me/api/functions/get-user-social"
)

var accountTable = os.Getenv("accountTableName")

func countUp(table dynamo.Table, followedUserID string, followingUserID string) error {
	if err := table.
		Update("id", followedUserID).
		Range("sort", "social").
		SetExpr("'followers' = 'followers' + ?", 1).
		Run(); err != nil {
		return err
	}

	if err := table.
		Update("id", followingUserID).
		Range("sort", "social").
		SetExpr("'followings' = 'followings' + ?", 1).
		Run(); err != nil {
		return err
	}

	return nil
}

func handler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(accountTable)

	fmt.Printf("%+v", event)

	followedUserID := event["prev"].(map[string]interface{})["result"].(map[string]interface{})["id"].(string)
	followingUserID := event["prev"].(map[string]interface{})["result"].(map[string]interface{})["follow"].(string)

	countUp(authTable, followedUserID, followingUserID)

	return map[string]interface{}{
		"id": followedUserID,
	}, nil
}

func main() {
	lambda.Start(handler)
}
