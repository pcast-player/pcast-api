package models

type Feed struct {
	Id  int    `json:"id" gorm:"primary_key"`
	Url string `json:"url"`
}

type CreateFeedInput struct {
	Url string `json:"url" binding:"required"`
}
