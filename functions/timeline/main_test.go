package main

import (
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"
	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
	timeline "github.com/myuon/portals-me/functions/stream-timeline-feed/lib"
)

func TestCanGetTimeline(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := db.Table(os.Getenv("EntityTable"))
	timelineTable := db.Table(os.Getenv("TimelineTable"))

	testUser := authenticator.User{
		ID:   "user##" + uuid.Must(uuid.NewV4()).String(),
		Name: "name",
	}
	if err := entityTable.Put(testUser.ToDBO()).Run(); err != nil {
		t.Fatal(err)
	}

	timelineTable.Put(timeline.TimelineItem{
		ID: uuid.Must(uuid.NewV4()).String(),
		FeedEvent: feed.FeedEvent{
			UserID:    testUser.ID,
			Timestamp: time.Now().Unix(),
			EventName: "INSERT_COLLECTION",
			ItemID:    "collection##1234",
			Entity:    nil,
		},
	})

	if items, err := DoListTimeline(
		testUser.ID,
		entityTable,
		timelineTable,
	); err != nil {
		t.Fatalf("%+v", items)

		t.Fatal(err)
	}
}
