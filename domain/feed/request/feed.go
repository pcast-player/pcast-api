package request

// Feed represents a feed request
// @model Feed
type Feed struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}
