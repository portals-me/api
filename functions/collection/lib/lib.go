package collection

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type CollectionDBO struct {
	ID             string            `json:"id" dynamo:"id"`
	CommentMembers []string          `json:"comment_members" dynamo:"comment_members"`
	CommentCount   int               `json:"comment_count" dynamo:"comment_count"`
	Media          []string          `json:"media" dynamo:"media"`
	Cover          map[string]string `json:"cover" dynamo:"cover"`
	Owner          string            `json:"sort_value" dynamo:"sort_value"`
	Title          string            `json:"title" dynamo:"title"`
	CreatedAt      int64             `json:"created_at" dynamo:"created_at"`
	Sort           string            `json:"sort" dynamo:"sort"`
	Description    string            `json:"description" dynamo:"description"`
}

type Collection struct {
	ID             string            `json:"id"`
	CommentMembers []string          `json:"comment_members"`
	CommentCount   int               `json:"comment_count"`
	Media          []string          `json:"media"`
	Cover          map[string]string `json:"cover"`
	Owner          string            `json:"owner"`
	Title          string            `json:"title"`
	CreatedAt      int64             `json:"created_at"`
	Description    string            `json:"description"`
}

func (collection Collection) ToDBO() CollectionDBO {
	return CollectionDBO{
		ID:             "collection##" + collection.ID,
		CommentMembers: collection.CommentMembers,
		CommentCount:   collection.CommentCount,
		Media:          collection.Media,
		Cover:          collection.Cover,
		Owner:          collection.Owner,
		Title:          collection.Title,
		CreatedAt:      collection.CreatedAt,
		Description:    collection.Description,
		Sort:           "collection##detail",
	}
}

func (collection CollectionDBO) FromDBO() Collection {
	return Collection{
		ID:             strings.Split(collection.ID, "collection##")[1],
		CommentMembers: collection.CommentMembers,
		CommentCount:   collection.CommentCount,
		Media:          collection.Media,
		Cover:          collection.Cover,
		Owner:          collection.Owner,
		Title:          collection.Title,
		CreatedAt:      collection.CreatedAt,
		Description:    collection.Description,
	}
}

func ParseCollections(attrs []map[string]dynamodb.AttributeValue) []Collection {
	var collectionsDBO []CollectionDBO
	dynamodbattribute.UnmarshalListOfMaps(attrs, &collectionsDBO)

	var collections []Collection
	for _, v := range collectionsDBO {
		collections = append(collections, v.FromDBO())
	}

	return collections
}

func ParseCollection(attr map[string]dynamodb.AttributeValue) Collection {
	var collectionDBO CollectionDBO
	dynamodbattribute.UnmarshalMap(attr, &collectionDBO)

	return collectionDBO.FromDBO()
}

func DumpCollection(collection Collection) (map[string]dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(collection.ToDBO())
}
