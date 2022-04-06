package authentication

import (
	"auth-project/tools"
	"encoding/base32"
	"errors"
	"github.com/dgryski/dgoogauth"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"strings"
)

func VerifyGoogleTwoFactorAuthCode(token string, secret string) error {

	// setup the one-time-password configuration.
	otpConfig := &dgoogauth.OTPConfig{
		Secret:      strings.TrimSpace(secret),
		WindowSize:  3,
		HotpCounter: 0,
	}

	trimmedToken := strings.TrimSpace(token)

	// Validate token
	ok, err := otpConfig.Authenticate(trimmedToken)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("invalid token")
	}

	return nil
}

func GenerateGoogleTwoFactorAuthQrCode(usrEmail string) ([]byte, string, error) {

	secret := base32.StdEncoding.EncodeToString([]byte(tools.RandStr(16, "alphanum")))

	authLink := "otpauth://totp/" + usrEmail + "?secret=" + secret + "&issuer=" + viper.GetString("project_name")
	qrCodeByte, err := qrcode.Encode(authLink, qrcode.Medium, 256)
	if err != nil {
		return nil, "", err
	}

	return qrCodeByte, secret, nil
}
