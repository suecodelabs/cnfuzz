package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suecodelabs/cnfuzz/example/controllers"
	_ "github.com/suecodelabs/cnfuzz/example/docs"
	"github.com/suecodelabs/cnfuzz/example/model"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Test api for fuzzer poc
// @version 1.2
// @description This is a test api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	seedData()

	var router = gin.Default()
	controllers.NewRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}

// Create some test data
func seedData() {
	model.CreateToken("Finish the test API")
	model.CreateToken("0d5989ed-d60c-470e-b1b5-576fcf0f5d8c")
}
