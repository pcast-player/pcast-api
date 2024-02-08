package main

import (
	"flag"
	"fmt"
	echoSwagger "github.com/swaggo/echo-swagger"
	"pcast-api/config"
	"pcast-api/db"
	_ "pcast-api/docs"
	feedController "pcast-api/domain/feed/controller"
	feedStore "pcast-api/domain/feed/store"
	"pcast-api/router"
)

const usage = `Usage:
  -c, --config Path to config file (default: config.toml)
`

// @title PCast REST-API
// @version 0.1
// @BasePath  /api
// @host localhost:8080
func main() {
	var cfgFile string

	flag.StringVar(&cfgFile, "config", "config.toml", "path to config file")
	flag.StringVar(&cfgFile, "c", "config.toml", "path to config file")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	c := config.New(cfgFile)
	r := router.New(c)
	apiV1 := r.Group("/api")
	d := db.New(c)

	db.AutoMigrate(d)

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	fs := feedStore.New(d)
	fc := feedController.New(fs)

	fc.Register(apiV1)

	r.Logger.Fatal(r.Start(c.Server.GetAddress()))
}
