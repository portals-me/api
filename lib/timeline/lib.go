package timeline

type TimelineItem struct {
	ID         string `dynamo:"id"`
	Target     string `dynamo:"target"`
	OriginalID string `dynamo:"original_id"`
	UpdatedAt  int64  `dynamo:"updated_at"`
}
