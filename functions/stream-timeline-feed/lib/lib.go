package timeline

import (
	feed "github.com/myuon/portals-me/functions/stream-activity-feed/lib"
)

type TimelineItem struct {
	ID string `json:"id" dynamo:"id"`
	feed.FeedEvent
}

func BuildTimelineItem(ownerID string, item feed.FeedEvent) TimelineItem {
	return TimelineItem{
		ID:        ownerID,
		FeedEvent: item,
	}
}
