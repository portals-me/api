package main

import (
	"os"
	"testing"
	"time"

	authenticator "../authenticator/lib"
	feed "../entity-stream/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func TestListFeed(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := db.Table(os.Getenv("EntityTable"))
	feedTable := db.Table(os.Getenv("FeedTable"))

	testEvent := feed.FeedEvent{
		UserID:    "test-user",
		Timestamp: time.Now().Unix(),
		EventName: "INSERT_COLLECTION",
		ItemID:    "test-item-id",
	}

	err := feedTable.Put(testEvent).Run()
	if err != nil {
		t.Error(err)
	}

	err = entityTable.Put(authenticator.UserDBO{
		ID:   "test-user",
		Sort: "user##detail",
		Name: "hoge",
	}).Run()
	if err != nil {
		t.Error(err)
	}

	items, err := DoListFeed("hoge", entityTable, feedTable)
	if err != nil {
		t.Error(err)
	}

	if len(items) < 1 {
		t.Fatalf("Items has a wrong size: %+v", items)
	}

	event := items[0]
	if !(event.UserID == testEvent.UserID &&
		event.Timestamp == testEvent.Timestamp &&
		event.EventName == testEvent.EventName &&
		event.ItemID == testEvent.ItemID) {
		t.Fatalf("Argument does not match: %+v vs %+v", event, testEvent)
	}
}
