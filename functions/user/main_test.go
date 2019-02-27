package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofrs/uuid"
	"github.com/guregu/dynamo"
	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"
	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
	. "github.com/myuon/portals-me/functions/user/lib"
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
		t.Fatal(err)
	}

	err = entityTable.Put(authenticator.UserDBO{
		ID:   "test-user",
		Sort: "user##detail",
		Name: "hoge",
	}).Run()
	if err != nil {
		t.Fatal(err)
	}

	items, err := DoListFeed("hoge", entityTable, feedTable)
	if err != nil {
		t.Fatal(err)
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

func TestCanFollow(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := db.Table(os.Getenv("EntityTable"))

	user1 := authenticator.User{
		ID:   "user##1",
		Name: "test-user-1",
	}
	user2 := authenticator.User{
		ID:   "user##2",
		Name: "test-user-2",
	}

	if err := entityTable.
		Put(user1.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}
	if err := entityTable.
		Put(user2.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}

	if err := DoFollowUser(
		user1.ID,
		user2.Name,
		entityTable,
	); err != nil {
		t.Fatal(err)
	}

	var record UserFollowRecord
	if err := entityTable.
		Get("id", "user##2").
		Range("sort", dynamo.Equal, "user##follow-user##1").
		One(&record); err != nil {
		t.Fatal(err)
	}
}

func TestCanUpdate(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := db.Table(os.Getenv("EntityTable"))

	testUser := authenticator.User{
		ID:          "user##" + uuid.Must(uuid.NewV4()).String(),
		Name:        "test",
		DisplayName: "test-display-name",
	}
	if err := entityTable.
		Put(testUser.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}

	newName := uuid.Must(uuid.NewV4()).String()
	if err := DoUpdateUser(
		testUser.ID,
		authenticator.User{
			Name: newName,
		},
		entityTable,
	); err != nil {
		t.Fatal(err)
	}

	updatedUser, err := DoGetUser(newName, entityTable)
	if err != nil {
		t.Fatal(err)
	}
	if !(updatedUser.ID == testUser.ID && updatedUser.DisplayName == testUser.DisplayName) {
		t.Fatalf("Unexpected Argument: %+v", updatedUser)
	}

	if err := DoUpdateUser(
		testUser.ID,
		authenticator.User{
			DisplayName: "piyo",
			Picture:     "piyo",
		},
		entityTable,
	); err != nil {
		t.Fatal(err)
	}

	updatedUser, err = DoGetUser(newName, entityTable)
	if err != nil {
		t.Fatal(err)
	}
	if !(updatedUser.ID == testUser.ID &&
		updatedUser.DisplayName == "piyo" &&
		updatedUser.Picture == "piyo") {
		t.Fatalf("Unexpected Argument: %+v", updatedUser)
	}
}

func TestCannotUpdateWithNonuniqueName(t *testing.T) {
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := db.Table(os.Getenv("EntityTable"))

	testUser1 := authenticator.User{
		ID:          "user##" + uuid.Must(uuid.NewV4()).String(),
		Name:        "test1",
		DisplayName: "test-display-name-1",
	}
	if err := entityTable.
		Put(testUser1.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}
	testUser2 := authenticator.User{
		ID:          "user##" + uuid.Must(uuid.NewV4()).String(),
		Name:        "test2",
		DisplayName: "test-display-name-2",
	}
	if err := entityTable.
		Put(testUser2.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}

	if err := DoUpdateUser(
		testUser1.ID,
		authenticator.User{
			Name: "test2",
		},
		entityTable,
	); !strings.Contains(err.Error(), "user_name already exists") {
		t.Fatal(err)
	}
}
