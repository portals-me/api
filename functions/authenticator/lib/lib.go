package authenticator

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"

	. "../verifier"
)

type SignUpInput struct {
	Form struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Picture     string `json:"picture"`
	} `json:"form"`
	Logins Logins `json:"logins"`
}

type User struct {
	ID          string `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

type UserDBO struct {
	ID          string `json:"id"`
	Sort        string `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"sort_value"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func (user User) ToDBO() UserDBO {
	return UserDBO{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
		Sort:        "user##detail",
	}
}

func (user UserDBO) FromDBO() User {
	return User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
	}
}

func ParseUser(attr map[string]dynamodb.AttributeValue) User {
	var userDBO UserDBO
	dynamodbattribute.UnmarshalMap(attr, &userDBO)

	return userDBO.FromDBO()
}

func DumpUser(user User) (map[string]dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(user.ToDBO())
}

type JwtPayload struct {
	ID          string `json:"id"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func (user User) ToJwtPayload() JwtPayload {
	return JwtPayload{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
	}
}
