package models

type Feed struct {
	Id  int    `json:"id" gorm:"primary_key"`
	Url string `json:"url"`
}
