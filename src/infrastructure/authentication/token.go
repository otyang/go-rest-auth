package authentication

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/matoous/go-nanoid/v2"
	"io/ioutil"
	"os"
	"time"

	"auth-project/src/domain/model"
	"github.com/go-playground/validator/v10"
)

var (
	path, _               = os.Getwd()
	privateKey, publicKey = MustLoadRSA(path+"/rsa_keys/private_key.pem", path+"/rsa_keys/public_key.pem")
)

type JwtConfigurator struct {
	AccessTokenMaxAge        time.Duration
	RefreshTokenMaxAge       time.Duration
	TwoFactorAuthTokenMaxAge time.Duration
	SigningMethod            *jwt.SigningMethodRSA
}

func NewJwtConfigurator(atMaxAge, rtMaxAge, TwoFactorAuthTokenMaxAge time.Duration) *JwtConfigurator {

	return &JwtConfigurator{
		AccessTokenMaxAge:        atMaxAge,
		RefreshTokenMaxAge:       rtMaxAge,
		TwoFactorAuthTokenMaxAge: TwoFactorAuthTokenMaxAge,
		SigningMethod:            jwt.SigningMethodRS512,
	}
}

func MustLoadRSA(privateKeyFilename, publicKeyFilename string) (*rsa.PrivateKey, *rsa.PublicKey) {
	b, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		panic(err)
	}

	private, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadFile(publicKeyFilename)
	if err != nil {
		panic(err)
	}

	public, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		panic(err)
	}

	return private, public
}

func (jc *JwtConfigurator) GenerateTokenPair(userID, sessionID string) (*model.TokenDetails, error) {
	var err error

	td := new(model.TokenDetails)
	td.AtExpires = time.Now().UTC().Add(jc.AccessTokenMaxAge).Unix()
	td.AtID, err = gonanoid.New()
	if err != nil {
		return nil, err
	}
	td.RtExpires = time.Now().UTC().Add(jc.RefreshTokenMaxAge).Unix()
	td.RtID, err = gonanoid.New()
	if err != nil {
		return nil, err
	}

	accessClaims := model.AccessClaims{
		Authorized: true,
		SessionID:  sessionID,
		AtID:       td.AtID,
		UserID:     userID,
		Exp:        td.AtExpires,
		Type:       model.AccessTokenTypeAuth,
	}

	token := jwt.NewWithClaims(jc.SigningMethod, accessClaims)

	td.AccessToken, err = token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := model.RefreshClaims{
		RtID:      td.RtID,
		SessionID: sessionID,
		UserID:    userID,
		Exp:       td.RtExpires,
	}

	refreshToken := jwt.NewWithClaims(jc.SigningMethod, refreshClaims)

	td.RefreshToken, err = refreshToken.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (jc *JwtConfigurator) GetAccessTokenClaims(accessToken string) (*model.AccessClaims, error) {
	var claims model.AccessClaims
	tkn, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("unauthorized")
		}

		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	err = validator.New().Struct(&claims)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}

func (jc *JwtConfigurator) GetRefreshTokenClaims(refreshToken string) (*model.RefreshClaims, error) {

	var claims model.RefreshClaims
	tkn, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("unauthorized")
		}

		return nil, errors.New("unauthorized")
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	err = validator.New().Struct(&claims)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}

func (jc *JwtConfigurator) GenerateTwoFactorAuthToken(userID, sessionID string) (*model.AccessTokenDetails, error) {
	var err error

	td := new(model.AccessTokenDetails)
	td.AtExpires = time.Now().UTC().Add(jc.TwoFactorAuthTokenMaxAge).Unix()
	td.AtID, err = gonanoid.New()
	if err != nil {
		return nil, err
	}

	claims := model.AccessClaims{
		Authorized: false,
		AtID:       td.AtID,
		SessionID:  sessionID,
		UserID:     userID,
		Exp:        td.AtExpires,
		Type:       model.AccessTokenTypeTwoFactorAuth,
	}

	token := jwt.NewWithClaims(jc.SigningMethod, claims)

	td.AccessToken, err = token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return td, nil
}
