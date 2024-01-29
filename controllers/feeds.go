package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pcast-api/models"
)

func GetFeeds(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetFeeds())
}

func CreateFeed(c *gin.Context) {
	var input models.CreateFeedInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed := models.Feed{Url: input.Url}

	models.CreateFeed(&feed)

	c.JSON(http.StatusOK, feed)
}
