package main

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity/cognitoidentityiface"
	"github.com/guregu/dynamo"

	. "github.com/myuon/portals-me/functions/authenticator/api"
	. "github.com/myuon/portals-me/functions/authenticator/lib"
	. "github.com/myuon/portals-me/functions/authenticator/verifier"
	collection "github.com/myuon/portals-me/functions/collection/lib"
)

type fakeCustomProvider struct {
	cognitoidentityiface.CognitoIdentityAPI
	customID string
}

func (provider *fakeCustomProvider) GetIdpID(logins Logins) (string, error) {
	if provider.customID != "" {
		return provider.customID, nil
	} else {
		return "fake-idp", nil
	}
}

type operation struct {
	request  string
	argument interface{}
}

type fakeSigner struct {
}

func (signer fakeSigner) Sign(payload []byte) ([]byte, error) {
	return payload, nil
}

func TestCanSignUpWithGoogle(t *testing.T) {
	ddb := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	idp := &fakeCustomProvider{}
	signer := &fakeSigner{}

	testUser := struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Picture     string `json:"picture"`
	}{
		Name:        "test_name",
		DisplayName: "test_display_name",
		Picture:     "test_picture",
	}
	input := SignUpInput{
		Form: testUser,
		Logins: Logins{
			Google: "id_token",
		},
	}
	_, identityID, err := DoSignUp(input, idp, entityTable, signer)

	if err != nil {
		t.Fatal(err)
	}
	if identityID != "fake-idp" {
		t.Fatal(err)
	}

	var userDBO UserDBO
	if err := entityTable.
		Get("id", "user##"+identityID).
		Range("sort", dynamo.Equal, "user##detail").
		One(&userDBO); err != nil {
		t.Fatal(err)
	}

	user := userDBO.FromDBO()
	if !(user.Name == testUser.Name &&
		user.DisplayName == testUser.DisplayName &&
		user.Picture == testUser.Picture) {
		t.Fatalf("Argument does not match: %+v", user)
	}

	var colDBO collection.CollectionDBO
	if err := entityTable.
		Get("id", "collection##"+user.Name).
		Range("sort", dynamo.Equal, "collection##detail").
		One(&colDBO); err != nil {
		t.Fatal(err)
	}
	col := colDBO.FromDBO()

	if !(col.ID == testUser.Name) {
		t.Errorf("Argument does not match: %+v", col)
	}
}

func TestCanSignInWithoutUserCollectionTwice(t *testing.T) {
	ddb := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	entityTable := ddb.Table(os.Getenv("EntityTable"))

	testUser := User{
		ID:   "user##user-id",
		Name: "user-name",
	}

	if err := entityTable.Put(testUser.ToDBO()).Run(); err != nil {
		t.Fatal(err)
	}

	idp := &fakeCustomProvider{
		customID: "user-id",
	}
	signer := &fakeSigner{}
	logins := Logins{
		Twitter: "id_token",
	}

	if _, err := DoSignIn(logins, idp, entityTable, signer); err != nil {
		t.Fatal(err)
	}

	var colDBO collection.CollectionDBO
	if err := entityTable.
		Get("id", "collection##"+testUser.Name).
		Range("sort", dynamo.Equal, "collection##detail").
		One(&colDBO); err != nil {
		t.Fatal(err)
	}

	col := colDBO.FromDBO()

	if !(col.ID == testUser.Name &&
		col.Title == testUser.Name) {
		t.Fatalf("Argument does not match: %+v", col)
	}

	if _, err := DoSignIn(logins, idp, entityTable, signer); err != nil {
		t.Fatal(err)
	}
}
