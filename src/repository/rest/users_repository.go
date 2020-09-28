package rest

import (
	"encoding/json"
	"github.com/csrias/bookstore_oauth-api/src/domain/users"
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "localhost:8080",
		Timeout: 100 * time.Microsecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(e, p string) (*users.User, *errors.RestErr) {
	req := users.UserLoginRequest{
		Email:    e,
		Password: p,
	}
	res := usersRestClient.Post("/users/login", req)
	if res == nil || res.Response == nil {
		return nil, errors.NewInternalServerError("request timeout")
	}
	if res.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(res.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}
	var user users.User
	if err := json.Unmarshal(res.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users response")
	}
	return &user, nil
}
