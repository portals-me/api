package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofrs/uuid"
	"github.com/guregu/dynamo"
	dynamo_helper "github.com/portals-me/api/lib/dynamo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/portals-me/api/lib/timeline"
	"github.com/portals-me/api/lib/user"
)

var timelineTableName = os.Getenv("timelineTableName")
var userTableName = os.Getenv("userTableName")
var timelineTable dynamo.Table
var userRepository user.Repository

func createNotifiedItemID(itemID string, followerID string) string {
	return followerID + "-" + itemID
}

// item should be {id: string, owner: string, updated_at: number}
func createItemsToFollowers(item map[string]*dynamodb.AttributeValue) ([]interface{}, error) {
	ownerID := *item["owner"].S
	updatedAt, _ := strconv.ParseInt(*item["updated_at"].N, 10, 64)

	followers, err := userRepository.ListFollowers(ownerID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to listFollowers")
	}

	var items []interface{}
	for _, follower := range append(followers, ownerID) {
		items = append(items, timeline.TimelineItem{
			ID:         uuid.Must(uuid.NewV4()).String(),
			Target:     follower,
			OriginalID: *item["id"].S,
			UpdatedAt:  updatedAt,
		})
	}

	return items, nil
}

func createItemsToDelete(item map[string]*dynamodb.AttributeValue) ([]dynamo.Keyed, error) {
	itemID := item["id"].String()

	var timelineItems []timeline.TimelineItem
	if err := timelineTable.
		Get("original_id", itemID).
		Index("original_id").
		All(&timelineItems); err != nil {
		return nil, errors.Wrap(err, "Failed to query original_id items")
	}

	var items []dynamo.Keyed
	for _, timelineItem := range timelineItems {
		items = append(items, dynamo.Keys{timelineItem.ID})
	}

	return items, nil
}

func handler(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		fmt.Printf("%+v\n", record)
		message := record.SNS.Message

		var dbEvent events.DynamoDBEventRecord
		if err := json.Unmarshal([]byte(message), &dbEvent); err != nil {
			return errors.Wrap(err, "Failed to unmarshal")
		}

		if dbEvent.EventName == "MODIFY" || dbEvent.EventName == "INSERT" {
			item, err := dynamo_helper.AsDynamoDBAttributeValues(dbEvent.Change.NewImage)
			if err != nil {
				return errors.Wrap(err, "Failed to parse NewImage")
			}

			items, err := createItemsToFollowers(item)
			if err != nil {
				return errors.Wrap(err, "Failed to create itemsToFollowers")
			}

			if _, err := timelineTable.Batch().Write().Put(items...).Run(); err != nil {
				return errors.Wrap(err, "Failed to BatchWrite")
			}
		} else if dbEvent.EventName == "REMOVE" {
			item, err := dynamo_helper.AsDynamoDBAttributeValues(dbEvent.Change.Keys)
			if err != nil {
				return errors.Wrap(err, "Failed to parse Keys")
			}

			items, err := createItemsToDelete(item)
			if err != nil {
				return errors.Wrap(err, "Failed to create items to delete")
			}

			if _, err := timelineTable.Batch().Write().Delete(items...).Run(); err != nil {
				return errors.Wrap(err, "Failed to BatchWrite")
			}
		} else {
			fmt.Printf("%+v\n", dbEvent)
			panic("Not supported EventName: " + dbEvent.EventName)
		}
	}

	return nil
}

func main() {
	svc := dynamodb.New(session.New())
	userRepository = user.NewRepositoryFromAWS(svc, userTableName)
	timelineTable = dynamo.NewFromIface(svc).Table(timelineTableName)

	lambda.Start(handler)
}
