package app

import (
	"github.com/csrias/bookstore_oauth-api/src/domain/access_token"
	"github.com/csrias/bookstore_oauth-api/src/http"
	"github.com/csrias/bookstore_oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.Create)
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
