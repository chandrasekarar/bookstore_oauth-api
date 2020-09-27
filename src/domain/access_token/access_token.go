package access_token

import (
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"strings"
	"time"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) validate() *errors.RestErr {
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

func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
