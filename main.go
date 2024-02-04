package main

import (
	"pcast-api/controllers"
	"pcast-api/db"
	"pcast-api/router"
)

func main() {
	r := router.New()

	d := db.New()
	db.AutoMigrate(d)

	feedsController := controllers.NewFeedsController(d)

	r.GET("/feeds", feedsController.GetFeeds)
	r.POST("/feeds", feedsController.CreateFeed)

	r.Logger.Fatal(r.Start(":8080"))
}
