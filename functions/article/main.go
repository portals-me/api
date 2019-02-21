package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/service/s3"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"

	collection "github.com/myuon/portals-me/functions/collection/lib"
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
	Owner       string            `json:"owner"`
	CreatedAt   int64             `json:"created_at"`
}

// for DynamoDB
type ArticleDBO struct {
	ID          string            `json:"id" dynamo:"id"`
	Sort        string            `json:"sort" dynamo:"sort"`
	Entity      map[string]string `json:"entity" dynamo:"entity"`
	Title       string            `json:"title" dynamo:"title"`
	Description string            `json:"description" dynamo:"description"`
	Owner       string            `json:"sort_value" dynamo:"sort_value"`
	CreatedAt   int64             `json:"created_at" dynamo:"created_at"`
}

func (article Article) ToDBO(collectionID string) ArticleDBO {
	return ArticleDBO{
		ID:          "collection##" + collectionID,
		Sort:        "article##" + article.ID,
		Entity:      article.Entity,
		Title:       article.Title,
		Description: article.Description,
		Owner:       article.Owner,
		CreatedAt:   article.CreatedAt,
	}
}

func (dto ArticleDBO) FromDBO() Article {
	return Article{
		ID:          strings.Replace(dto.Sort, "article##", "", 1),
		Entity:      dto.Entity,
		Title:       dto.Title,
		Description: dto.Description,
		Owner:       dto.Owner,
		CreatedAt:   dto.CreatedAt,
	}
}

func doList(
	collectionID string,
	entityTable dynamo.Table,
) (int, []Article, error) {
	var articleDBOs []ArticleDBO
	if err := entityTable.
		Get("id", "collection##"+collectionID).
		Range("sort", dynamo.BeginsWith, "article##").
		All(&articleDBOs); err != nil {
		return 502, nil, err
	}

	if len(articleDBOs) == 0 {
		return 200, []Article{}, nil
	}

	var articles = make([]Article, len(articleDBOs))
	for index, value := range articleDBOs {
		articles[index] = value.FromDBO()
	}

	return 200, articles, nil
}

func doCreate(
	collectionID string,
	userID string,
	createInput map[string]interface{},
	entityTable dynamo.Table,
) (int, string, error) {
	// care for Entity struct
	entityMap := map[string]string{}
	entity := createInput["entity"].(map[string]interface{})
	for key, value := range entity {
		entityMap[key] = value.(string)
	}

	article := Article{
		ID:          uuid.NewV4().String(),
		Entity:      entityMap,
		Title:       createInput["title"].(string),
		Description: createInput["description"].(string),
		CreatedAt:   time.Now().Unix(),
		Owner:       userID,
	}

	var colDBO collection.CollectionDBO
	if err := entityTable.
		Get("id", "collection##"+collectionID).
		Range("sort", dynamo.Equal, "collection##detail").
		One(&colDBO); err != nil {
		if err == dynamo.ErrNotFound {
			return 404, "", err
		}

		return 502, "", errors.Wrap(err, "getCollection failed")
	}
	col := colDBO.FromDBO()

	if userID != col.Owner {
		return 403, "", errors.New("AccessDenied")
	}

	if err := entityTable.Put(article.ToDBO(collectionID)).Run(); err != nil {
		return 502, "", errors.Wrapf(err, "putArticle failed, %+v", article.ToDBO(collectionID))
	}

	return 201, article.ID, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := event.RequestContext.Authorizer

	sess := session.Must(session.NewSession())

	ddb := dynamo.NewFromIface(dynamodb.New(sess))
	entityTable := ddb.Table(os.Getenv("EntityTable"))
	collectionID := event.PathParameters["collectionId"]

	if event.HTTPMethod == "GET" {
		if _, ok := event.PathParameters["articleId"]; ok {
		} else {
			statusCode, articles, err := doList(collectionID, entityTable)

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

			S3 := s3.New(sess)

			req, _ := S3.PutObjectRequest(&s3.PutObjectInput{
				Bucket: aws.String("portals-me-storage-users"),
				Key:    aws.String(user["id"].(string) + "/" + collectionID + "/" + event.Body),
			})
			signedURL, err := req.Presign(15 * time.Minute)

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
				user["id"].(string),
				createInput,
				entityTable,
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
