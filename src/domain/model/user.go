package model

import (
	"github.com/uptrace/bun"
	"time"
)

const (
	UserRole = "user"
)

// Base entity
type User struct {
	bun.BaseModel `bun:"table:users,alias:usr"`

	ID               string    `json:"id" bun:"id,pk"`
	FullName         string    `json:"full_name" bun:",nullzero"`
	UserName         string    `json:"user_name" bun:",nullzero"`
	Email            string    `json:"email" bun:",nullzero"`
	Password         string    `json:"password" bun:",nullzero"`
	Phone            string    `json:"phone" bun:",nullzero"`
	Hash             string    `json:"hash" bun:",nullzero"`
	ReferralLink     string    `json:"referral_link" bun:",nullzero"`
	Role             string    `json:"role" bun:",nullzero"`
	IsActive         bool      `json:"is_active"`
	IsEmailVerified  bool      `json:"is_email_verified"`
	IsPhoneVerified  bool      `json:"is_phone_verified"`
	IsGoogleVerified bool      `json:"is_google_verified"`
	GoogleSecret     string    `json:"google_secret" bun:",nullzero"`
	CreatedAt        time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:now()"`
	UpdatedAt        time.Time `json:"updated_at" bun:"updated_at,nullzero"`

	ReferralUser *User  `json:"referral_user" bun:"rel:belongs-to,join:referral=referral_link"`
	Referral     string `json:"referral" bun:",nullzero"`
}

// UserGetMyProfileResp entity for get my profile resp
type UserGetMyProfileResp struct {
	ID               string `json:"id"`
	FullName         string `json:"full_name"`
	UserName         string `json:"user_name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	ReferralLink     string `json:"referral_link"`
	Role             string `json:"role"`
	IsEmailVerified  bool   `json:"is_email_verified"`
	IsPhoneVerified  bool   `json:"is_phone_verified"`
	IsGoogleVerified bool   `json:"is_google_verified"`

	CreatedAt time.Time `json:"created_at"`

	ReferralUser *UserReferralGet `json:"referral_user" bun:"rel:belongs-to,join:referral=referral_link"`
	Referral     string           `json:"referral" bun:",nullzero"`
}

// UserReferralGet entity for get my profile referral resp
type UserReferralGet struct {
	ID string `json:"id"`
}

// UserCreateReq entity of the create request
type UserCreateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserUpdResp entity of the update resp
type UserUpdResp struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Referral string `json:"referral"`
}

// UserUpdateInfoData entity of the update info data for function
type UserUpdateInfoData struct {
	FullName string `json:"full_name"`
}

// UserUpdateInfoReq entity of the update info request
type UserUpdateInfoReq struct {
	FullName string `json:"full_name"`
}

// UserPhoneUpdateReq entity of the update phone request
type UserPhoneUpdateReq struct {
	Code2fa string `json:"code_2fa"`
	Phone   string `json:"phone"`
}

// UserEmailUpdateReq entity of the update email request
type UserEmailUpdateReq struct {
	Code2fa string `json:"code_2fa"`
	Email   string `json:"email"`
}

// UserChangePasswordReq entity of the change password request
type UserChangePasswordReq struct {
	Code2fa     string `json:"code_2fa"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserChangePasswordData entity of the change password data for function
type UserChangePasswordData struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserSessionData entity of the user  session data
type UserSessionData struct {
	UserID    string `redis:"user_id"`
	UserAgent string `json:"user_agent"`
	ClientIp  string `json:"client_ip"`
}

// UserRedisSessionData entity of the user redis session data
type UserRedisSessionData struct {
	AtID string `redis:"at_id"`
	RtID string `redis:"rt_id"`
}

// SignUpReq entity of the sign-up request
type SignUpReq struct {
	Login     string `json:"login"`
	LoginType string `json:"login_type"`
	Code2fa   string `json:"code_2fa"`
	Password  string `json:"password"`
	Referral  string `json:"referral"`
}

// SignUpActivateUserData entity of to activate sign up user data for function
type SignUpActivateUserData struct {
	Login     string `json:"login"`
	LoginType string `json:"login_type"`
	Password  string `json:"password"`
	Referral  string `json:"referral"`
}

// SignUpSend2faCodeReq entity of to send 2fa sign-up code request
type SignUpSend2faCodeReq struct {
	Login     string `json:"login"`
	LoginType string `json:"login_type"`
}

// Send2faCodeForResetUserPasswordReq entity of to send 2fa reset password code request
type Send2faCodeForResetUserPasswordReq struct {
	Target     string `json:"target"`
	TargetType string `json:"target_type"`
}

// VerifyResetUserPassword2faСodeReq entity of to verify 2fa reset password code request
type VerifyResetUserPassword2faСodeReq struct {
	Target      string `json:"target"`
	Code2fa     string `json:"code_2fa"`
	Code2faType string `json:"code_2fa_type"`
	NewPassword string `json:"new_password"`
}
