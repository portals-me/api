package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go-v2/service/s3"

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

func doList(
	collectionID string,
	ddb dynamodbiface.DynamoDBAPI,
) (int, []Article, error) {
	result, err := ddb.QueryRequest(&dynamodb.QueryInput{
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
		return 502, nil, err
	}

	if len(result.Items) == 0 {
		return 200, []Article{}, nil
	}

	var articles []Article
	var articleDTOs []ArticleDTO
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &articleDTOs)
	for _, value := range articleDTOs {
		articles = append(articles, value.FromDTO())
	}

	return 200, articles, nil
}

func doCreate(
	collectionID string,
	user map[string]interface{},
	createInput map[string]interface{},
	ddb dynamodbiface.DynamoDBAPI,
) (int, string, error) {
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
		return 400, "", err
	}

	result, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + collectionID)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()
	if len(result.Item) == 0 {
		return 404, "", errors.New("Not Found")
	}
	if err != nil {
		return 502, "", err
	}

	if user["id"].(string) != *result.Item["sort_value"].S {
		return 403, "", errors.New("AccessDenied")
	}

	_, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      collection,
	}).Send()

	if err != nil {
		return 502, "", err
	}

	return 201, articleID, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	ddb := dynamodb.New(cfg)
	collectionID := event.PathParameters["collectionId"]

	if event.HTTPMethod == "GET" {
		if _, ok := event.PathParameters["articleId"]; ok {
		} else {
			statusCode, articles, err := doList(collectionID, ddb)

			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       err.Error(),
					Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
					StatusCode: statusCode,
				}, nil
			}

			out, _ := json.Marshal(articles)
			return events.APIGatewayProxyResponse{
				Body:       string(out),
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: statusCode,
			}, nil
		}
	} else if event.HTTPMethod == "POST" {
		// OMG
		if strings.HasSuffix(event.Resource, "/articles-presigned") {
			if event.Body == "" {
				return events.APIGatewayProxyResponse{}, errors.New("Empty filepath")
			}

			S3 := s3.New(cfg)

			signedURL, _ := S3.PutObjectRequest(&s3.PutObjectInput{
				Bucket: aws.String("portals-me-storage-users"),
				Key:    aws.String(user["id"].(string) + "/" + collectionID + "/" + event.Body),
			}).Presign(15 * time.Minute)

			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			return events.APIGatewayProxyResponse{
				Body:       signedURL,
				Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
				StatusCode: 200,
			}, nil
		} else {
			var createInput map[string]interface{}
			json.Unmarshal([]byte(event.Body), &createInput)

			statusCode, articleID, err := doCreate(
				collectionID,
				user,
				createInput,
				ddb,
			)

			if err != nil {
				return events.APIGatewayProxyResponse{
					Body: err.Error(),
					Headers: map[string]string{
						"Access-Control-Allow-Origin": "*",
					},
					StatusCode: statusCode,
				}, nil
			}

			return events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Access-Control-Allow-Origin": "*",
					"Location":                    "/collections/" + collectionID + "/articles/" + articleID,
				},
				StatusCode: statusCode,
			}, nil
		}
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
