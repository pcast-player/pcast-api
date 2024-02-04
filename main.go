package main

import (
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/router"
)

func main() {
	r := router.New()

	d := db.New()
	db.AutoMigrate(d)

	feedsController := controller.NewFeedController(d)

	r.GET("/feeds", feedsController.GetFeeds)
	r.POST("/feeds", feedsController.CreateFeed)

	r.Logger.Fatal(r.Start(":8080"))
}
