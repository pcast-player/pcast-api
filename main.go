package main

import (
	"flag"
	"fmt"
	"github.com/swaggo/echo-swagger"
	"pcast-api/config"
	"pcast-api/db"
	_ "pcast-api/docs"
	"pcast-api/domain/feed"
	"pcast-api/router"
	feedStore "pcast-api/store/feed"
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
	apiGroup := r.Group("/api")
	d := db.New(c)
	fs := feedStore.New(d)

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	feed.New(apiGroup, fs)

	r.Logger.Fatal(r.Start(c.Server.GetAddress()))
}
