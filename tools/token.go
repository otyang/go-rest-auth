package tools

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"

	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"

	"auth-project/src/domain/model"
)

var (
	path, _               = os.Getwd()
	privateKey, publicKey = MustLoadRSA(path+"/rsa_keys/private_key.pem", path+"/rsa_keys/public_key.pem")
)

type JwtConfigurator struct {
	AccessTokenMaxAge  time.Duration
	RefreshTokenMaxAge time.Duration
	SigningMethod      *jwt.SigningMethodRSA
}

func NewJwtConfigurator(atMaxAge, rtMaxAge time.Duration) *JwtConfigurator {

	return &JwtConfigurator{
		AccessTokenMaxAge:  atMaxAge,
		RefreshTokenMaxAge: rtMaxAge,
		SigningMethod:      jwt.SigningMethodRS512,
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

func (jc *JwtConfigurator) GenerateTokenPair(userID string) (*model.TokenDetails, error) {
	var err error

	td := new(model.TokenDetails)
	td.AtExpires = time.Now().UTC().Add(jc.AccessTokenMaxAge).Unix()
	td.AtUuid = uuid.New().String()
	td.RtExpires = time.Now().UTC().Add(jc.RefreshTokenMaxAge).Unix()
	td.RtUuid = uuid.New().String()

	claims := jwt.MapClaims{
		"authorized": true,
		"at_id":      td.AtUuid,
		"usr_id":     userID,
		"exp":        td.AtExpires,
	}

	token := jwt.NewWithClaims(jc.SigningMethod, claims)

	td.AccessToken, err = token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	claims = jwt.MapClaims{
		"rt_id":  td.RtUuid,
		"usr_id": userID,
		"exp":    td.RtExpires,
	}

	refreshToken := jwt.NewWithClaims(jc.SigningMethod, claims)

	td.RefreshToken, err = refreshToken.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (jc *JwtConfigurator) GetAccessTokenClaims(accessToken string) (*model.AccessClaims, error) {
	token, err := jwt.Parse(strings.Split(accessToken, "Bearer ")[1], func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	parts := bytes.Split([]byte(token.Raw), []byte("."))
	if len(parts) != 3 {
		return nil, err
	}

	payload, err := Base64Decode(parts[1])
	if err != nil {
		return nil, err
	}

	refreshClaims := &model.AccessClaims{}
	err = json.Unmarshal(payload, refreshClaims)
	if err != nil {
		return nil, err
	}

	return refreshClaims, nil
}

func (jc *JwtConfigurator) GetRefreshTokenClaims(refreshToken string) (*model.RefreshClaims, error) {
	token, err := jwt.Parse(strings.Split(refreshToken, "Bearer ")[1], func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	parts := bytes.Split([]byte(token.Raw), []byte("."))
	if len(parts) != 3 {
		return nil, err
	}

	payload, err := Base64Decode(parts[1])
	if err != nil {
		return nil, err
	}

	refreshClaims := &model.RefreshClaims{}
	err = json.Unmarshal(payload, refreshClaims)
	if err != nil {
		return nil, err
	}

	return refreshClaims, nil
}

// Base64Decode decodes "src" to jwt base64 url format.
// We could use the base64.RawURLEncoding but the below is a bit faster.
func Base64Decode(src []byte) ([]byte, error) {
	if n := len(src) % 4; n > 0 {
		// JWT: Because of no trailing '=' let's suffix it
		// with the correct number of those '=' before decoding.
		src = append(src, bytes.Repeat([]byte("="), 4-n)...)
	}

	buf := make([]byte, base64.URLEncoding.DecodedLen(len(src)))
	n, err := base64.URLEncoding.Decode(buf, src)
	return buf[:n], err
}
