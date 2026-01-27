package oauth

// LoginResponse represents a successful OAuth login response
// @model LoginResponse
type LoginResponse struct {
	Token string `json:"token"`
}
