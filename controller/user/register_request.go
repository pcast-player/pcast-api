package user

// RegisterRequest represents a user request
// @model RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
