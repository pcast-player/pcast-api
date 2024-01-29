package main

import (
	"github.com/gin-gonic/gin"
	"pcast-api/controllers"
	"pcast-api/models"
)

func main() {
	models.ConnectDatabase()
	router := gin.Default()

	router.GET("/feeds", controllers.GetFeeds)

	router.Run()
}
