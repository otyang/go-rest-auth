package model

import "github.com/golang-jwt/jwt/v4"

const (
	PostfixRefreshToken          = "_refresh"
	AccessTokenTypeAuth          = "auth"
	AccessTokenTypeTwoFactorAuth = "two_factor_auth"
)

// AccessClaims a custom access token claims structure.
type AccessClaims struct {
	jwt.RegisteredClaims
	Authorized bool   `json:"authorized" `
	AtID       string `json:"at_id" validate:"required"`
	UserID     string `json:"usr_id" validate:"required"`
	SessionID  string `json:"session_id" validate:"required"`
	Exp        int64  `json:"exp" validate:"required"`
	Type       string `json:"type" validate:"required"`
}

// RefreshClaims a custom refresh token claims structure.
type RefreshClaims struct {
	jwt.RegisteredClaims
	RtID      string `json:"rt_id" validate:"required"`
	SessionID string `json:"session_id" validate:"required"`
	UserID    string `json:"usr_id" validate:"required"`
	Exp       int64  `json:"exp" validate:"required"`
}

type TokenDetails struct {
	SessionID    string
	AccessToken  string
	RefreshToken string
	AtID         string
	RtID         string
	AtExpires    int64
	RtExpires    int64
}

type AccessTokenDetails struct {
	AccessToken string
	AtID        string
	AtExpires   int64
}
