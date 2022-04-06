package model

// AccessClaims a custom access token claims structure.
type AccessClaims struct {
	Authorized bool   `json:"authorized,required"`
	AtID       string `json:"at_id,required"`
	UserID     string `json:"usr_id,required"`
	Exp        int64  `json:"exp,required"`
}

// RefreshClaims a custom refresh token claims structure.
type RefreshClaims struct {
	RtID   string `json:"rt_id"`
	UserID string `json:"usr_id"`
	Exp    int64  `json:"exp"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtUuid       string
	RtUuid       string
	AtExpires    int64
	RtExpires    int64
}
