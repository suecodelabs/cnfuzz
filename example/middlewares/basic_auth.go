package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suecodelabs/cnfuzz/example/model"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Search user in the slice of allowed credentials
		// user, found := pairs.searchCredential(c.requestHeader("Authorization"))
		passedToken := c.Request.Header.Get("Authorization")
		token, found := model.SearchToken(passedToken)
		if !found {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(gin.AuthUserKey, token.Token)
	}
}
