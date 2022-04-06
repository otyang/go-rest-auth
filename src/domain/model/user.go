package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

// Base entity
type User struct {
	bun.BaseModel `bun:"table:users,alias:usr"`

	ID           uuid.UUID    `json:"id" bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	FullName     string       `json:"full_name"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Phone        string       `json:"phone"`
	Pin          string       `json:"pin"`
	ReferralLink string       `json:"referral_link"`
	CreatedAt    time.Time    `json:"created_at" bun:",notnull,default:now()"`
	UpdatedAt    bun.NullTime `json:"updated_at"`
}

// The entity of the create request
type UserCreateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// The entity of the updat resp
type UserUpdResp struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	ReferralLink string    `json:"referral_link"`
}

// The entity of the update user info request
type UserUpdateReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// The entity of the updqte pin request
type UserPinUpdateReq struct {
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

// The entity of the change password request
type UserPinDeleteReq struct {
	Pin string `json:"pin"`
}

// The entity of the change password request
type UserChangePasswordReq struct {
	ID          uuid.UUID `json:"id"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
}

// The entity of the user redis session data
type UserRedisSessionData struct {
	UserID string `redis:"user_id"`
	AtID   string `redis:"at_id"`
	RtID   string `redis:"rt_id"`
}
