package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/suecodelabs/cnfuzz/example/middlewares"
)

func NewRoutes(router *gin.Engine) *gin.Engine {
	mainGroup := router.Group("/api")
	mainGroup.Use(middlewares.BasicAuth())

	AddTodoController(mainGroup)
	AddAuthController(mainGroup)

	return router
}
