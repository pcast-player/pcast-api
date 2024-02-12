package request

// CreateFeedRequest represents a feed request
// @model CreateFeedRequest
type CreateFeedRequest struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}
