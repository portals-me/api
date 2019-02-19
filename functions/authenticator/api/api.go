package api

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"

	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/signer"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
	collection "github.com/myuon/portals-me/functions/collection/lib"
)

func createUserCollection(
	user User,
	ddb dynamodbiface.DynamoDBAPI,
) error {
	result, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String("collection##" + user.Name)},
			"sort": {S: aws.String("collection##detail")},
		},
	}).Send()

	if err != nil {
		return err
	}

	if result.Item != nil {
		return nil
	}

	item, err := collection.DumpCollection(collection.Collection{
		ID:          user.Name,
		Owner:       user.ID,
		Title:       user.Name,
		Description: "",
		Cover: map[string]string{
			"color": "red lighten-3",
			"sort":  "solid",
		},
		Media:          []string{},
		CommentMembers: []string{user.ID},
		CommentCount:   0,
		CreatedAt:      time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	if _, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Item:      item,
	}).Send(); err != nil {
		return err
	}

	return err
}

type CreateInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Cover       map[string]string `json:"cover"`
}

func DoSignUp(
	input SignUpInput,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (string, string, error) {
	identityID, err := idp.GetIdpID(input.Logins)

	user := User{
		ID:          "user##" + identityID,
		CreatedAt:   time.Now().Unix(),
		Name:        input.Form.Name,
		DisplayName: input.Form.DisplayName,
		Picture:     input.Form.Picture,
	}

	item, err := DumpUser(user)
	if err != nil {
		return "", "", err
	}

	if _, err = ddb.PutItemRequest(&dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("EntityTable")),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}).Send(); err != nil {
		return "", "", err
	}

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return "", "", err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return "", "", err
	}

	if err = createUserCollection(user, ddb); err != nil {
		return "", "", err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	return string(body), identityID, nil
}

func DoSignIn(
	logins Logins,
	idp ICustomProvider,
	ddb dynamodbiface.DynamoDBAPI,
	signer ISigner,
) (string, error) {
	identityID, err := idp.GetIdpID(logins)
	if err != nil {
		return "", err
	}

	userID := "user##" + identityID

	getItemReq, err := ddb.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("EntityTable")),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(userID)},
			"sort": {S: aws.String("user##detail")},
		},
	}).Send()
	if err != nil {
		return "", err
	}

	if getItemReq.Item["id"].S == nil {
		return "", errors.New("UserNotExist: " + userID)
	}

	user := ParseUser(getItemReq.Item)

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return "", err
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(map[string]interface{}{
		"id_token": string(token),
		"user":     string(jsn),
	})

	err = createUserCollection(user, ddb)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
