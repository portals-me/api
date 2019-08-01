package user

import (
	aws_dynamo "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/portals-me/account/lib/user"
)

type FollowRelation struct {
	ID     string `dynamo:"id"`
	Follow string `dynamo:"follow"`
}

type Repository user.Repository

// NewRepository creates a new instance
func NewRepository(table dynamo.Table) Repository {
	return Repository(user.NewRepository(table))
}

// NewRepositoryFromAWS creates a new instance from aws-sdk DynamoDB instance
func NewRepositoryFromAWS(client *aws_dynamo.DynamoDB, tableName string) Repository {
	return NewRepository(dynamo.NewFromIface(client).Table(tableName))
}

// ListFollowers lists the id numbers of followers
func (repo Repository) ListFollowers(userID string) ([]string, error) {
	table := user.Repository(repo).AsDynamoTable()

	var records []FollowRelation
	if err := table.
		Get("id", userID).
		Range("sort", dynamo.BeginsWith, "follow@@").
		All(&records); err != nil {
		return nil, err
	}

	userIDs := make([]string, len(records))
	for index, record := range records {
		userIDs[index] = record.Follow
	}

	return userIDs, nil
}
