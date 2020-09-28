package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"github.com/federicoleon/bookstore_users-api/utils/crypto_utils"
)

const (
	expirationTime            = 24
	grantTypePassword         = "password"
	grantTypeClientCredential = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// used for client credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredential:
		break
	default:
		return errors.NewBadRequest("invalid grant_type parameter")
	}

	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequest("invalid access token id")
	}
	if at.UserID <= 0 {
		return errors.NewBadRequest("invalid user id")
	}
	if at.ClientID <= 0 {
		return errors.NewBadRequest("invalid client id")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequest("invalid expire time")
	}
	return nil
}

func GetNewAccessToken(userID int64) AccessToken {
	return AccessToken{
		UserID:  userID,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
