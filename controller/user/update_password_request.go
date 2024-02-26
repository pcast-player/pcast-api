package user

// UpdatePasswordRequest represents a user request
// @model UpdatePasswordRequest
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}
