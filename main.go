package main

import (
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/router"
	"pcast-api/store"
)

func main() {
	r := router.New()

	d := db.New()
	db.AutoMigrate(d)

	feedStore := store.NewFeedStore(d)
	feedController := controller.NewFeedController(feedStore)

	r.GET("/feeds", feedController.GetFeeds)
	r.POST("/feeds", feedController.CreateFeed)

	r.Logger.Fatal(r.Start(":8080"))
}
