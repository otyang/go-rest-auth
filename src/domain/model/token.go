package model

import (
	"github.com/uptrace/bun"
	"time"
)

const (
	TokenReasonVerification  = "verification"
	TokenReasonTwoFactorAuth = "two_factor_auth"
	TokenReasonResetPassword = "reset_password"
	TokenReasonSignUp        = "sign_up"
	TokenReasonAuthByQrCode  = "auth_qr_code"

	TokenTypeGoogle = "google"
	TokenTypeEmail  = "email"
	TokenTypePhone  = "phone"
)

var (
	TokenTimeSendErr = ""
)

// Base entity
type Token struct {
	bun.BaseModel `bun:"table:tokens,alias:tkn"`

	ID     string `json:"id" bun:"id,pk"`
	UserID string `json:"user_id" bun:",nullzero"`
	Target string `json:"target" bun:",nullzero"`
	Value  string `json:"value" bun:",nullzero"`
	Reason string `json:"reason" bun:",nullzero"`
	Type   string `json:"type" bun:",nullzero"`
	IsUsed bool   `json:"is_used"`

	ExpiresAT time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:now()"`
}

// SendTarget2faCodeReq entity for send target 2fa code request
type SendTarget2faCodeReq struct {
	Target      string `json:"target"`
	Code2faType string `json:"code_2fa_type"`
}

// Send2faCodeData entity for send target 2fa code data for function
type Send2faCodeData struct {
	UserID      string `json:"user_id"`
	Target      string `json:"target"`
	Code2faType string `json:"code_2fa_type"`
	Reason      string `json:"reason"`
}

// Verify2faCodeReq entity for verify 2fa code request
type Verify2faCodeReq struct {
	Code2fa     string `json:"code_2fa"`
	Code2faType string `json:"code_2fa_type"`
}

// VerifyCodeData entity for verify 2fa code data for function
type VerifyCodeData struct {
	UserID   string `json:"user_id"`
	Target   string `json:"target"`
	Code     string `json:"code"`
	CodeType string `json:"code_type"`
	Reason   string `json:"reason"`
}
