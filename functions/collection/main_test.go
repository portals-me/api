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

func TestCanCreateUpdateAndDelete(t *testing.T) {
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

	collection.Title = "test-title-2"
	if err := DoUpdate(collection.Collection, entityTable); err != nil {
		t.Fatal(err)
	}

	collection, err = DoGet(
		collectionID,
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !(collection.Title == "test-title-2") {
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
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanCreateAndList(t *testing.T) {
	ddb := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	testUser := authenticator.User{
		ID:          "test-user-2",
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
	if _, err := DoCreate(
		testInput,
		testUser.ID,
		entityTable,
	); err != nil {
		t.Fatal(err)
	}

	testInput = CreateInput{
		Title:       "test-title-2",
		Description: "hoge-2",
		Cover:       map[string]string{},
	}
	if _, err := DoCreate(
		testInput,
		testUser.ID,
		entityTable,
	); err != nil {
		t.Fatal(err)
	}

	cols, err := DoList(testUser.ID, entityTable)
	if err != nil {
		t.Fatal(err)
	}

	if !(len(cols) == 2) {
		t.Fatalf("Argument does not match: %+v", cols)
	}
}
