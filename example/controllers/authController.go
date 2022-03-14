package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suecodelabs/cnfuzz/example/model"
)

func AddAuthController(routes *gin.RouterGroup) {
	authGroup := routes.Group("/auth")
	authGroup.GET("/check", checkToken)
	authGroup.GET("/ping", ping)
}

// checkToken godoc
// @Summary Check if a token is valid
// @Description Endpoint that can be used to check if a token is valid. Returns 200 if valid and 401 if not valid
// @Tags tokens
// @Param authorization header string true "Authorization token"
// @Success 200 {string} string
// @Failure 401 {string} string
// @Router /api/auth/check [GET]
func checkToken(ctx *gin.Context) {
	authToken := ctx.GetHeader("Authorization")
	if authToken == "" {
		ctx.String(http.StatusBadRequest, "Authorization header is empty")
	}

	_, err := model.GetToken(authToken)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Token is invalid")
	} else {
		ctx.String(http.StatusOK, "Token is valid")
	}
}

// ping godoc
// @Summary Ping endpoint, logged-in users only
// @Description Ping endpoint for logged-in users. Returns pong.
// @Security ApiKeyAuth
// @Tags tokens
// @Success 200 {string} string
// @Failure 401 {string} string
// @Router /api/auth/ping [GET]
func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}
