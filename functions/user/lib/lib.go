package user

type UserFollowRecord struct {
	ID    string `dynamo:"id"`
	Sort  string `dynamo:"sort"`
	Value string `dynamo:"sort_value"`
}
