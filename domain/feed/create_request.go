package feed

// CreateRequest represents a feed request
// @model CreateRequest
type CreateRequest struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}
