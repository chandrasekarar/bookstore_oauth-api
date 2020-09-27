package db

import (
	"github.com/csrias/bookstore_oauth-api/src/client/cassandra"
	"github.com/csrias/bookstore_oauth-api/src/domain/access_token"
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken          = "SELECT access_token, client_id, expires, user_id FROM access_tokens WHERE access_token=?;"
	queryInsertAccessToken       = "INSERT INTO access_tokens(access_token, client_id, expires, user_id) VALUES(?,?,?,?);"
	queryUpdateAccessTokenExpiry = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetByID(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetByID(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.ClientID,
		&result.Expires,
		&result.UserID,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFound("no access token found with given id")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryInsertAccessToken,
		at.AccessToken,
		at.ClientID,
		at.Expires,
		at.UserID,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateAccessTokenExpiry,
		at.AccessToken,
		at.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
