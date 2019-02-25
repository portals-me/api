package api

import (
	"strings"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-lambda-go/events"

	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
	. "github.com/myuon/portals-me/functions/stream-timeline-feed/lib"
	user "github.com/myuon/portals-me/functions/user/lib"
)

func ProcessEvent(
	event events.DynamoDBEventRecord,
	entityTable dynamo.Table,
	timelineTable dynamo.Table,
) error {
	feedEvent, err := feed.FeedEventFromDynamoEvent(event)
	if err != nil {
		return err
	}

	if strings.HasPrefix(feedEvent.EventName, "INSERT_") {
		var followers []user.UserFollowRecord
		if err := entityTable.
			Get("id", feedEvent.UserID).
			Range("sort", dynamo.BeginsWith, "user##follow-").
			All(&followers); err != nil {
			return err
		}

		items := make([]interface{}, len(followers))
		for index, follower := range followers {
			target := strings.Split(follower.Sort, "user##follow-")[1]
			items[index] = BuildTimelineItem(target, feedEvent)
		}

		if _, err = timelineTable.
			Batch().
			Write().
			Put(items...).
			Run(); err != nil {
			return err
		}
	} else if strings.HasPrefix(feedEvent.EventName, "DELETE_") {
		var items []TimelineItem
		err := timelineTable.
			Get("item_id", feedEvent.ItemID).
			Index("ItemID").
			Project("id", "timestamp").
			All(&items)
		if err != nil {
			return err
		}

		deleteItems := make([]dynamo.Keyed, len(items))
		for index, item := range items {
			deleteItems[index] = dynamo.Keys{item.ID, item.Timestamp}
		}

		if _, err := timelineTable.
			Batch("id", "timestamp").
			Write().
			Delete(deleteItems...).
			Run(); err != nil {
			return err
		}
	}

	return nil
}
