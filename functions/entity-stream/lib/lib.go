package feed

type FeedEvent struct {
	UserID    string `dynamo:"user_id"`
	Timestamp int64  `dynamo:"timestamp"`
	EventName string `dynamo:"event_name"`
	ItemID    string `dynamo:"item_id"`
}
