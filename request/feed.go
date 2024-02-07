package request

// Feed represents a feed request
// @model Feed
type Feed struct {
	URL string `json:"url" validate:"required,url"`
}
