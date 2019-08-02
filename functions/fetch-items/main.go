package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	"github.com/portals-me/api/lib/timeline"
)

var accountTableName = os.Getenv("accountTableName")
var postTableName = os.Getenv("postTableName")
var timelineTableName = os.Getenv("timelineTableName")
var ownerCache = make(map[string]map[string]interface{})
var postCache = make(map[string]map[string]interface{})

var accountTable dynamo.Table
var postTable dynamo.Table
var timelineTable dynamo.Table

func replaceOwner(cache map[string]map[string]interface{}, items []map[string]interface{}) ([]map[string]interface{}, error) {
	for _, item := range items {
		var record map[string]interface{}

		if owner, ok := ownerCache[item["owner"].(string)]; ok {
			record = owner
		} else {
			if err := accountTable.
				Get("id", item["owner"]).
				Range("sort", dynamo.Equal, "detail").
				One(&record); err != nil {
				return nil, err
			}

			ownerCache[item["owner"].(string)] = record
		}

		item["owner_user"] = record
	}

	return items, nil
}

func replacePost(cache map[string]map[string]interface{}, itemIDs []string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}

	for _, itemID := range itemIDs {
		var record map[string]interface{}

		if item, ok := cache[itemID]; ok {
			record = item
		} else {
			if err := postTable.
				Get("id", itemID).
				Range("sort", dynamo.Equal, "summary").
				One(&record); err != nil {
				return nil, err
			}

			cache[itemID] = record
		}

		items = append(items, record)
	}

	return items, nil
}

func fetchTimeline(userID string) ([]map[string]interface{}, error) {
	var items []timeline.TimelineItem
	if err := timelineTable.
		Get("target", userID).
		Index("target").
		All(&items); err != nil {
		return nil, err
	}

	itemIDs := make([]string, len(items))
	for index, item := range items {
		itemIDs[index] = item.OriginalID
	}

	posts, err := replacePost(postCache, itemIDs)
	if err != nil {
		return nil, err
	}

	postsWithOwner, err := replaceOwner(ownerCache, posts)
	if err != nil {
		return nil, err
	}

	return postsWithOwner, nil
}

func handler(ctx context.Context, event map[string]interface{}) ([]map[string]interface{}, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	accountTable = db.Table(accountTableName)
	postTable = db.Table(postTableName)
	timelineTable = db.Table(timelineTableName)

	return fetchTimeline(event["id"].(string))
}

func main() {
	lambda.Start(handler)
}
