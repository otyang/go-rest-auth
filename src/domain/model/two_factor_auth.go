package model

// TwoFactorAuthCodeVerifyReq entity for verify 2fa code request
type TwoFactorAuthCodeVerifyReq struct {
	Code2fa string `json:"code_2fa"`
}

// TwoFactorAuthDeleteReq entity for delete 2fa code request
type TwoFactorAuthDeleteReq struct {
	Password string `json:"password"`
	Type     string `json:"type"`
}

// TwoFactorAuthSetUpReq entity for set up 2fa code request
type TwoFactorAuthSetUpReq struct {
	Code2fa     string `json:"code_2fa"`
	Secret      string `json:"secret"`
	Code2faType string `json:"code_2fa_type"`
}
