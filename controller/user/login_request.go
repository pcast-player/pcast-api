package user

// LoginRequest represents a login request
// @model LoginRequest
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
