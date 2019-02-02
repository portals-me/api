package main

import (
	"errors"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type OldUser struct {
	CreatedAt   int64  `json:"created_at"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Sort        string `json:"sort"`
}

type NewUser struct {
	CreatedAt   int64  `json:"created_at"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Picture     string `json:"picture"`
	Sort        string `json:"sort"`
	Name        string `json:"sort_value"`
}

func (data OldUser) transform() NewUser {
	return NewUser{
		CreatedAt:   data.CreatedAt,
		DisplayName: data.DisplayName,
		ID:          data.ID,
		Name:        data.Name,
		Picture:     data.Picture,
		Sort:        "user##" + data.Sort,
	}
}

func transformUser(attr map[string]dynamodb.AttributeValue) (map[string]dynamodb.AttributeValue, error) {
	var oldUser OldUser
	err := dynamodbattribute.UnmarshalMap(attr, &oldUser)
	if err != nil {
		return nil, err
	}

	transformed, err := dynamodbattribute.MarshalMap(oldUser.transform())
	if err != nil {
		panic(err)
	}
	return transformed, nil
}

type OldCollection struct {
	CommentCount   int               `json:"comment_count"`
	CommentMembers []string          `json:"comment_members"`
	Cover          map[string]string `json:"cover"`
	CreatedAt      int64             `json:"created_at"`
	Description    string            `json:"description"`
	Entity         map[string]string `json:"entity"`
	ID             string            `json:"id"`
	Media          []string          `json:"media"`
	OwnedBy        string            `json:"owned_by"`
	Sort           string            `json:"sort"`
	Title          string            `json:"title"`
}

type NewCollection struct {
	CommentCount   int               `json:"comment_count"`
	CommentMembers []string          `json:"comment_members"`
	Cover          map[string]string `json:"cover"`
	CreatedAt      int64             `json:"created_at"`
	Description    string            `json:"description"`
	ID             string            `json:"id"`
	Media          []string          `json:"media"`
	Sort           string            `json:"sort"`
	Owner          string            `json:"sort_value"`
	Title          string            `json:"title"`
}

func (data OldCollection) transform() NewCollection {
	return NewCollection{
		CommentCount:   data.CommentCount,
		CommentMembers: data.CommentMembers,
		Cover:          data.Cover,
		CreatedAt:      data.CreatedAt,
		Description:    data.Description,
		ID:             data.ID,
		Media:          data.Media,
		Sort:           "collection##" + data.Sort,
		Owner:          data.OwnedBy,
		Title:          data.Title,
	}
}

func transformCollection(attr map[string]dynamodb.AttributeValue) (map[string]dynamodb.AttributeValue, error) {
	var oldCollection OldCollection
	err := dynamodbattribute.UnmarshalMap(attr, &oldCollection)
	if err != nil {
		return nil, err
	}

	transformed, err := dynamodbattribute.MarshalMap(oldCollection.transform())
	if err != nil {
		panic(err)
	}
	return transformed, nil
}

type OldArticle struct {
	CreatedAt   int64             `json:"created_at"`
	Description string            `json:"description"`
	Entity      map[string]string `json:"entity"`
	ID          string            `json:"id"`
	OwnedBy     string            `json:"owned_by"`
	Sort        string            `json:"sort"`
	Title       string            `json:"title"`
}

type NewArticle struct {
	CreatedAt   int64             `json:"created_at"`
	Description string            `json:"description"`
	Entity      map[string]string `json:"entity"`
	ID          string            `json:"id"`
	OwnedBy     string            `json:"sort_value"`
	Sort        string            `json:"sort"`
	Title       string            `json:"title"`
}

func (data OldArticle) transform() NewArticle {
	return NewArticle{
		CreatedAt:   data.CreatedAt,
		Description: data.Description,
		Entity:      data.Entity,
		ID:          data.ID,
		OwnedBy:     data.OwnedBy,
		Sort:        data.Sort,
		Title:       data.Title,
	}
}

func transformArticle(attr map[string]dynamodb.AttributeValue) (map[string]dynamodb.AttributeValue, error) {
	var oldCollection OldArticle
	err := dynamodbattribute.UnmarshalMap(attr, &oldCollection)
	if err != nil {
		return nil, err
	}

	transformed, err := dynamodbattribute.MarshalMap(oldCollection.transform())
	if err != nil {
		panic(err)
	}
	return transformed, nil
}

func transform(attr map[string]dynamodb.AttributeValue) (map[string]dynamodb.AttributeValue, error, map[string]dynamodb.AttributeValue) {
	ID := *attr["id"].S
	Sort := *attr["sort"].S
	if strings.HasPrefix(ID, "user##") {
		x, y := transformUser(attr)

		transformed, err := dynamodbattribute.MarshalMap(NewCollection{
			CommentCount: 0,
			CommentMembers: []string{
				ID,
			},
			Cover: map[string]string{
				"color": "teal darken-2",
				"sort":  "solid",
			},
			CreatedAt:   time.Now().Unix(),
			Description: "",
			ID:          "collection##" + uuid.Must(uuid.NewV4()).String(),
			Media:       []string{},
			Sort:        "collection##detail",
			Owner:       *attr["id"].S,
			Title:       *attr["display_name"].S,
		})
		if err != nil {
			panic(err)
		}

		return x, y, transformed
	} else if strings.HasPrefix(ID, "collection##") && strings.HasPrefix(Sort, "article##") {
		x, y := transformArticle(attr)
		return x, y, nil
	} else if strings.HasPrefix(ID, "collection##") {
		x, y := transformCollection(attr)
		return x, y, nil
	}

	return nil, errors.New("Unsupported record: " + ID), nil
}

func main() {
	cfg, _ := external.LoadDefaultAWSConfig(aws.Config{
		Region: "ap-northeast-1",
	})

	ddb := dynamodb.New(cfg)

	sourceTable := "portals-me-prod-entities"
	targetTable := "portals-me-dev-entities"

	req := ddb.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String(sourceTable),
	})
	pager := req.Paginate()

	for pager.Next() {
		page := pager.CurrentPage()

		var requests []dynamodb.WriteRequest

		for _, item := range page.Items {
			newData, err, optional := transform(item)
			if err != nil {
				panic(err)
			}

			requests = append(requests, dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: newData,
				},
			})

			if optional != nil {
				requests = append(requests, dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: optional,
					},
				})
			}
		}

		_, err := ddb.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]dynamodb.WriteRequest{
				targetTable: requests,
			},
		}).Send()

		if err != nil {
			panic(err)
		}
	}
}
