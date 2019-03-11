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

type AccountForm struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Picture     string `json:"picture"`
}

type SignUpInput struct {
	Form   AccountForm `json:"form"`
	Logins Logins      `json:"logins"`
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

type SignInOutput struct {
	IDToken string `json:"id_token"`
	User    User   `json:"user"`
}

func DoSignIn(
	logins Logins,
	idp ICustomProvider,
	entityTable dynamo.Table,
	signer ISigner,
) (SignInOutput, error) {
	identityID, err := idp.GetIdpID(logins)
	if err != nil {
		return SignInOutput{}, errors.Wrap(err, "getIdpID failed")
	}

	var userDBO UserDBO
	if err := entityTable.
		Get("id", "user##"+identityID).
		Range("sort", dynamo.Equal, "user##detail").
		One(&userDBO); err != nil {
		return SignInOutput{}, errors.Wrap(err, "getUserByIdpID failed")
	}
	user := userDBO.FromDBO()

	jsn, err := json.Marshal(user.ToJwtPayload())
	if err != nil {
		return SignInOutput{}, errors.Wrap(err, "json.Marshal failed")
	}

	token, err := signer.Sign(jsn)
	if err != nil {
		return SignInOutput{}, errors.Wrap(err, "signer.Sign failed")
	}

	if user.UserCollectionID == "" {
		collectionID, err := createUserCollection(user, entityTable)
		if err != nil {
			return SignInOutput{}, errors.Wrap(err, "createUserCollection failed")
		}

		if err := entityTable.
			Update("id", userDBO.ID).
			Range("sort", "user##detail").
			Set("user_collection_id", collectionID).
			Run(); err != nil {
			return SignInOutput{}, errors.Wrap(err, "updateUserCollectionID failed")
		}

		user.UserCollectionID = collectionID
	}

	return SignInOutput{
		IDToken: string(token),
		User:    user,
	}, nil
}
