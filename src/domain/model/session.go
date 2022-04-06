package model

import (
	"github.com/uptrace/bun"
	"time"
)

// Base entity
type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:ssn"`

	SessionID string `json:"session_id" bun:"session_id,pk"`

	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
	IsLogout  bool   `json:"is_logout"`

	ExpiresAT time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,nullzero"`

	User   *User  `json:"user" bun:"rel:belongs-to"`
	UserID string `json:"user_id"`
}
