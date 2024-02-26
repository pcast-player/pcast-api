package user

// CreateRequest represents a user request
// @model CreateRequest
type CreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
