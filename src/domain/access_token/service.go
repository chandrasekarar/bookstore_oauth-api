package access_token

import (
	"fmt"
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"strings"
)

type Repository interface {
	GetByID(string) (*AccessToken, *errors.RestErr)
	Create(AccessToken) *errors.RestErr
	UpdateExpirationTime(AccessToken) *errors.RestErr
}

type Service interface {
	GetByID(string) (*AccessToken, *errors.RestErr)
	Create(AccessToken) *errors.RestErr
	UpdateExpirationTime(AccessToken) *errors.RestErr
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetByID(accessTokenID string) (*AccessToken, *errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if accessTokenID == "" {
		return nil, errors.NewBadRequest("invalid access token id")
	}
	accessToken, err := s.repository.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, err
}

func (s *service) Create(at AccessToken) *errors.RestErr {
	fmt.Println("Create call")
	if err := at.validate(); err != nil {
		return err
	}
	return s.repository.Create(at)
}

func (s *service) UpdateExpirationTime(at AccessToken) *errors.RestErr {
	if err := at.validate(); err != nil {
		return err
	}
	return s.repository.UpdateExpirationTime(at)
}
