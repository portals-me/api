package main

import (
	"os"
	"testing"
	"time"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	authenticator "github.com/myuon/portals-me/functions/authenticator/lib"
	collection_api "github.com/myuon/portals-me/functions/collection/api"
)

func TestCreate(t *testing.T) {
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
	if err := entityTable.
		Put(testUser.ToDBO()).
		Run(); err != nil {
		t.Fatal(err)
	}

	collectionID, err := collection_api.DoCreate(
		collection_api.CreateInput{
			Title:       "test-title",
			Description: "test-description",
			Cover:       nil,
		},
		"test-user",
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}

	statusCode, _, err := doCreate(
		collectionID,
		"test-user",
		map[string]interface{}{
			"title":       "hoge",
			"description": "description",
			"entity": map[string]interface{}{
				"hoge": "piyo",
			},
		},
		entityTable,
	)
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 201 {
		t.Fatalf("Invalid StatusCode: %v", statusCode)
	}
}
