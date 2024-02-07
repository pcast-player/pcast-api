package main

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"pcast-api/config"
	"pcast-api/controller"
	"pcast-api/db"
	_ "pcast-api/docs"
	"pcast-api/router"
	"pcast-api/store"
)

// @title PCast REST-API
// @version 0.1
// @BasePath  /api
// @host localhost:8080
func main() {
	c := config.New("config.toml")
	r := router.New(c)
	apiV1 := r.Group("/api")
	d := db.New(c)

	db.AutoMigrate(d)

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	feedStore := store.New(d)
	feedController := controller.New(feedStore)

	feedController.Register(apiV1)

	r.Logger.Fatal(r.Start(c.Server.GetAddress()))
}
