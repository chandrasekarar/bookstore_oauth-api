package access_token

import (
	"github.com/csrias/bookstore_oauth-api/src/domain/access_token"
	"github.com/csrias/bookstore_oauth-api/src/repository/db"
	"github.com/csrias/bookstore_oauth-api/src/repository/rest"
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"strings"
)

type Service interface {
	GetByID(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if accessTokenID == "" {
		return nil, errors.NewBadRequest("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, err
}

func (s *service) Create(atr access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {

	if err := atr.Validate(); err != nil {
		return nil, err
	}

	// Authenticate the user against the User's API
	// Authenticate based on password
	user, err := s.restUsersRepo.LoginUser(atr.Username, atr.Password)
	if err != nil {
		return nil, err
	}

	// Generate the new access token in cassandra
	at := access_token.GetNewAccessToken(user.ID)
	at.Generate()
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
