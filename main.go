package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"pcast-api/controllers"
	"pcast-api/models"
	"pcast-api/validation"
)

func main() {
	models.ConnectDatabase()
	e := echo.New()
	e.Validator = &validation.ApiValidator{Validator: validator.New()}

	e.GET("/feeds", controllers.GetFeeds)
	e.POST("/feeds", controllers.CreateFeed)

	e.Logger.Fatal(e.Start(":8080"))
}
