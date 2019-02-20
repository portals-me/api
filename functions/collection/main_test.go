package main

import (
	"os"
	"testing"
	"time"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"

	. "github.com/myuon/portals-me/functions/collection/api"
)

func TestCanCreateAndDelete(t *testing.T) {
	ddb := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	testUser := authenticator.User{
		ID:          "test-user",
		CreatedAt:   time.Now().Unix(),
		Name:        "test-user-name",
		DisplayName: "test-user-display-name",
		Picture:     "",
	}
	if err := entityTable.Put(testUser.ToDBO()).Run(); err != nil {
		t.Fatal(err)
	}

	testInput := CreateInput{
		Title:       "test-title",
		Description: "hoge",
		Cover:       map[string]string{},
	}
	collectionID, err := DoCreate(
		testInput,
		testUser.ID,
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}

	collection, err := DoGet(
		collectionID,
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !(collectionID == collection.ID &&
		testInput.Title == collection.Title &&
		testInput.Description == collection.Description) {
		t.Fatalf("Invalid collection returned: %+v", collection)
	}

	err = DoDelete(
		collectionID,
		testUser.ID,
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}

	collection, err = DoGet(
		collectionID,
		entityTable,
	)
	if err == nil {
		t.Fatal(err)
	}
}
