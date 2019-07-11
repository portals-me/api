package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"

	"github.com/portals-me/account/lib/user"
)

var accountTable = os.Getenv("accountTableName")

type UserSocialRecord struct {
	Followers   int  `dynamo:"followers" json:"followers"`
	Followings  int  `dynamo:"followings" json:"followings"`
	IsFollowing bool `json:"is_following"`
}

type UserSocialRecordDDB struct {
	UserSocialRecord
	ID   string `dynamo:"id"`
	Sort string `dynamo:"sort"`
}

func (record UserSocialRecord) toDDB(id string) UserSocialRecordDDB {
	return UserSocialRecordDDB{
		UserSocialRecord: record,
		Sort:             "social",
		ID:               id,
	}
}

type UserMore struct {
	user.UserInfo
	UserSocialRecord
}

func getUserMore(table dynamo.Table, targetUserName string, requestUserID string) (UserMore, error) {
	var targetUser user.UserInfoDDB
	if err := table.
		Get("name", targetUserName).
		Index("name").
		One(&targetUser); err != nil {
		fmt.Printf("getTargetUser: %+v\n", err.Error())
		return UserMore{}, err
	}
	fmt.Printf("target: %+v\n", targetUser)

	var socialRecord UserSocialRecord
	if err := table.
		Get("id", targetUser.ID).
		Range("sort", dynamo.Equal, "social").
		One(&socialRecord); err != nil {
		if err == dynamo.ErrNotFound {
			socialRecord = UserSocialRecord{
				Followers:  0,
				Followings: 0,
			}

			if err := table.Put(socialRecord.toDDB(targetUser.ID)).Run(); err != nil {
				fmt.Printf("PutSocialRecord: %+v\n", err.Error())
				fmt.Printf("%+v\n", socialRecord)
				return UserMore{}, err
			}
		} else {
			fmt.Printf("GetTagetSocialRecord error: %+v\n", err.Error())
			return UserMore{}, err
		}
	}

	var followRecord interface{}
	if err := table.
		Get("id", targetUser.ID).
		Range("sort", dynamo.Equal, "follow@@"+requestUserID).
		One(&followRecord); err != nil {
		if err == dynamo.ErrNotFound {
			socialRecord.IsFollowing = false
		} else {
			fmt.Printf("Error in GetFollowRecord: %+v\n", err.Error())
			return UserMore{}, err
		}
	} else {
		socialRecord.IsFollowing = true
	}

	return UserMore{
		UserInfo:         targetUser.UserInfo,
		UserSocialRecord: socialRecord,
	}, nil
}

func handler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	sess := session.Must(session.NewSession())

	db := dynamo.NewFromIface(dynamodb.New(sess))
	authTable := db.Table(accountTable)

	targetUserName := event["arguments"].(map[string]interface{})["name"].(string)
	requestUserID := event["prev"].(map[string]interface{})["result"].(map[string]interface{})["id"].(string)

	userMore, err := getUserMore(authTable, targetUserName, requestUserID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v", userMore)

	return userMore, nil
}

func main() {
	lambda.Start(handler)
}
