package model

// Base entity
type AuthenticationData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
