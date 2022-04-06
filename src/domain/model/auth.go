package model

// AuthenticationReq entity for auth request
type AuthenticationReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthQrCodeReq entity for qr code auth request
type AuthQrCodeReq struct {
	Token string `json:"token"`
}
