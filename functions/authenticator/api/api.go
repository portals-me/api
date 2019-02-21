package api

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/guregu/dynamo"

	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/signer"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
	collection_api "github.com/myuon/portals-me/functions/collection/api"
)

func createUserCollection(
	user User,
	entityTable dynamo.Table,
) (string, error) {
	return collection_api.DoCreate(
		collection_api.CreateInput{
			Title:       user.Name,
			Description: user.DisplayName + "のコレクション",
			Cover: map[string]string{
				"color": "red lighten-3",
				"sort":  "solid",
			},
		},
		user.ID,
		entityTable,
	)
}

type CreateInput struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Cover       map[string]string `json:"cover"`
}

func DoSignUp(
	input SignUpInput,
	idp ICustomProvider,
	entityTable dynamo.Table,
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

	collectionID, err := createUserCollection(user, entityTable)
	if err != nil {
		return "", "", errors.Wrap(err, "createUserCollection failed")
	}

	user.UserCollectionID = collectionID

	if err := entityTable.
		Put(user.ToDBO()).
		If("attribute_not_exists(id)").
		Run(); err != nil {
		return "", "", errors.Wrap(err, "putUser failed")
	}

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return "", "", errors.Wrap(err, "marshalUser failed")
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return "", "", errors.Wrap(err, "sign failed")
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
	entityTable dynamo.Table,
	signer ISigner,
) (string, error) {
	identityID, err := idp.GetIdpID(logins)
	if err != nil {
		return "", err
	}

	userID := "user##" + identityID

	var userDBO UserDBO
	if err := entityTable.
		Get("id", userID).
		Range("sort", dynamo.Equal, "user##detail").
		One(&userDBO); err != nil {
		return "", err
	}
	user := userDBO.FromDBO()

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

	if user.UserCollectionID == "" {
		collectionID, err := createUserCollection(user, entityTable)
		if err != nil {
			return "", err
		}

		if err := entityTable.
			Update("id", userDBO.ID).
			Range("sort", "user##detail").
			Set("user_collection_id", collectionID).
			Run(); err != nil {
			return "", err
		}
	}

	return string(body), nil
}
