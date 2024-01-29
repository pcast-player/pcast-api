package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcast-api/models"
)

func GetFeeds(c *gin.Context) {
	var feeds []models.Feed
	models.DB.Find(&feeds)

	c.JSON(http.StatusOK, feeds)
}
