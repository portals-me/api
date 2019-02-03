package collection

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type CollectionDBO struct {
	ID             string            `json:"id"`
	CommentMembers []string          `json:"comment_members"`
	CommentCount   int               `json:"comment_count"`
	Media          []string          `json:"media"`
	Cover          map[string]string `json:"cover"`
	Owner          string            `json:"sort_value"`
	Title          string            `json:"title"`
	CreatedAt      int64             `json:"created_at"`
	Sort           string            `json:"sort"`
	Description    string            `json:"description"`
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

func (collection Collection) toDBO() CollectionDBO {
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

func (collection CollectionDBO) fromDBO() Collection {
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
		collections = append(collections, v.fromDBO())
	}

	return collections
}

func ParseCollection(attr map[string]dynamodb.AttributeValue) Collection {
	var collectionDBO CollectionDBO
	dynamodbattribute.UnmarshalMap(attr, &collectionDBO)

	return collectionDBO.fromDBO()
}

func DumpCollection(collection Collection) (map[string]dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(collection.toDBO())
}
