package feed

type FeedEvent struct {
	UserID    string                 `json:"user_id" dynamo:"user_id"`
	Timestamp int64                  `json:"timestamp" dynamo:"timestamp"`
	EventName string                 `json:"event_name" dynamo:"event_name"`
	ItemID    string                 `json:"item_id" dynamo:"item_id"`
	Entity    map[string]interface{} `json:"entity" dynamo:"entity"`
}
