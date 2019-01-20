package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

type Collection struct {
	ID             string            `json:"id"`
	CommentMembers []string          `json:"comment_members"`
	CommentCount   int               `json:"comment_count"`
	Media          []string          `json:"media"`
	Cover          map[string]string `json:"cover"`
	OwnedBy        string            `json:"owned_by"`
	Title          string            `json:"title"`
	CreatedAt      int64             `json:"created_at"`
	Sort           string            `json:"sort"`
	Description    string            `json:"description"`
}

/* Entity Examples

{
	type: share
	format: oembed
	url: https://~~~
}

{
	type: document
	format: markdown
	url: https://~~~
}
*/
type Article struct {
	ID          string            `json:"id"`
	Entity      map[string]string `json:"entity"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	OwnedBy     string            `json:"owned_by"`
	CreatedAt   int64             `json:"created_at"`
}

// for DynamoDB
type ArticleDTO struct {
	ID          string            `json:"id"`
	Sort        string            `json:"sort"`
	Entity      map[string]string `json:"entity"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	OwnedBy     string            `json:"owned_by"`
	CreatedAt   int64             `json:"created_at"`
}

func (article Article) ToDTO(collectionID string) ArticleDTO {
	return ArticleDTO{
		ID:          "collection##" + collectionID,
		Sort:        "article##" + article.ID,
		Entity:      article.Entity,
		Title:       article.Title,
		Description: article.Description,
		OwnedBy:     article.OwnedBy,
		CreatedAt:   article.CreatedAt,
	}
}

func (dto ArticleDTO) FromDTO() Article {
	return Article{
		ID:          strings.Replace(dto.Sort, "article##", "", 1),
		Entity:      dto.Entity,
		Title:       dto.Title,
		Description: dto.Description,
		OwnedBy:     dto.OwnedBy,
		CreatedAt:   dto.CreatedAt,
	}
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	Dynamo := dynamodb.New(cfg)
	collectionID := event.PathParameters["collectionId"]

	if event.HTTPMethod == "GET" {
		if _, ok := event.PathParameters["articleId"]; ok {
		} else {
			result, err := Dynamo.QueryRequest(&dynamodb.QueryInput{
				TableName:              aws.String(os.Getenv("EntityTable")),
				KeyConditionExpression: aws.String("id = :id and begins_with(sort, :sort)"),
				ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
					":id": {
						S: aws.String("collection##" + collectionID),
					},
					":sort": {
						S: aws.String("article##"),
					},
				},
			}).Send()

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			var articles []Article
			var articleDTOs []ArticleDTO
			dynamodbattribute.UnmarshalListOfMaps(result.Items, &articleDTOs)
			for _, value := range articleDTOs {
				articles = append(articles, value.FromDTO())
			}

			out, _ := json.Marshal(articles)

			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		}
	} else if event.HTTPMethod == "POST" {
		var createInput map[string]interface{}
		json.Unmarshal([]byte(event.Body), &createInput)

		articleID := uuid.Must(uuid.NewV4()).String()

		// care for Entity struct
		entityMap := map[string]string{}
		entity := createInput["entity"].(map[string]interface{})
		for key, value := range entity {
			entityMap[key] = value.(string)
		}

		collection, err := dynamodbattribute.MarshalMap(Article{
			ID:          articleID,
			Entity:      entityMap,
			Title:       createInput["title"].(string),
			Description: createInput["description"].(string),
			CreatedAt:   time.Now().Unix(),
			OwnedBy:     user["id"].(string),
		}.ToDTO(collectionID))

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		_, err = Dynamo.PutItemRequest(&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("EntityTable")),
			Item:      collection,
		}).Send()

		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Location":                    "/collections/" + collectionID + "/articles/" + articleID,
			},
			StatusCode: 201,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%+v", event.PathParameters),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		StatusCode: 400,
	}, nil
}

func main() {
	lambda.Start(handler)
}
