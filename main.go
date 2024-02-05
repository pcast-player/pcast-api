package main

import (
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/router"
	"pcast-api/store"
)

func main() {
	r := router.New()
	apiV1 := r.Group("/api")
	d := db.New()

	db.AutoMigrate(d)

	feedStore := store.NewFeedStore(d)
	feedController := controller.NewFeedController(feedStore)

	feedController.Register(apiV1)

	r.Logger.Fatal(r.Start(":8080"))
}
