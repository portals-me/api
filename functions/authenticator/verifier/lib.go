package verifier

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/GoogleIdTokenVerifier/GoogleIdTokenVerifier"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/cognitoidentityiface"
	"github.com/gomodule/oauth1/oauth"
)

type TwitterUser struct {
	ID              string `json:"id_str"`
	ScreenName      string `json:"screen_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

func GetTwitterClient() *oauth.Client {
	twitterKey := strings.Split(os.Getenv("TwitterKey"), ".")

	return &oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
		Credentials: oauth.Credentials{
			Token:  twitterKey[0],
			Secret: twitterKey[1],
		},
	}
}

func GetTwitterUser(cred *oauth.Credentials, user *TwitterUser) error {
	client := GetTwitterClient()
	resp, err := client.Get(nil, cred, "https://api.twitter.com/1.1/account/verify_credentials.json", url.Values{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return err
	}

	return nil
}

type IVerifier interface {
	Verify() (string, error)
}

type TwitterVerifier struct {
	Token string
}

func (str TwitterVerifier) Verify() (string, error) {
	twitterKey := strings.Split(str.Token, ".")

	var account TwitterUser
	err := GetTwitterUser(&oauth.Credentials{
		Token:  twitterKey[0],
		Secret: twitterKey[1],
	}, &account)

	if err != nil {
		return "", err
	}

	return "twitter-" + account.ID, nil
}

type GoogleVerifier struct {
	Token string
}

func (str GoogleVerifier) Verify() (string, error) {
	tokenInfo := GoogleIdTokenVerifier.Verify(str.Token, os.Getenv("GClientId"))

	if tokenInfo == nil {
		return "", errors.New("Invalid GoogleToken")
	}

	return str.Token, nil
}

type Logins struct {
	Twitter string `json:"twitter"`
	Google  string `json:"google"`
}

type ICustomProvider interface {
	GetIdpID(Logins) (string, error)
}

type CustomProvider struct {
	IdentityPoolID          string
	CognitoIdentityInstance cognitoidentityiface.CognitoIdentityAPI
}

func (provider *CustomProvider) GetIdpID(logins Logins) (string, error) {
	loginsMap := map[string]string{}
	if logins.Google != "" {
		verified, err := GoogleVerifier{Token: logins.Google}.Verify()
		if err != nil {
			return "", err
		}

		loginsMap["accounts.google.com"] = verified
	}
	if logins.Twitter != "" {
		verified, err := TwitterVerifier{Token: logins.Twitter}.Verify()
		if err != nil {
			return "", err
		}

		loginsMap["portals.me"] = verified
	}

	getIDReq, err := provider.CognitoIdentityInstance.GetOpenIdTokenForDeveloperIdentityRequest(&cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolId: aws.String(provider.IdentityPoolID),
		Logins:         loginsMap,
	}).Send()
	if err != nil {
		return "", err
	}

	return *getIDReq.IdentityId, nil
}
