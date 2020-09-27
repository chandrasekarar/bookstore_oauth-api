package http

import (
	"github.com/csrias/bookstore_oauth-api/src/domain/access_token"
	"github.com/csrias/bookstore_oauth-api/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccessTokenHandler interface {
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	//UpdateExpiryTime(ctx *gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetByID(c *gin.Context) {
	accessToken, err := handler.service.GetByID(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at access_token.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequest("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	if err := handler.service.Create(at); err != nil {
		restErr := errors.NewInternalServerError("error while creating access token")
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusCreated, at)
}

func (handler *accessTokenHandler) UpdateExpiryTime(c *gin.Context) {
	accessToken, err := handler.service.GetByID(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}
