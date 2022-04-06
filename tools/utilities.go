package tools

import (
	"crypto/rand"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/xid"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func VerifyPassword(password string) (string, error) {

	var number, upper bool
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return "", errors.New("bad password")
		}
	}

	if letters < 7 || !number || !upper {
		return "", errors.New("password must be 7 characters or more, one capital letter and one number")
	}

	return password, nil
}

func VerifyPhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsPossibleNumber(num) {
		return "", errors.New("invalid phone number")
	}
	return "+" + strconv.FormatInt(int64(*num.CountryCode), 10) + strconv.FormatUint(*num.NationalNumber, 10), nil
}

func ParseAndCheckToken(ctx *fiber.Ctx) (token string, err error) {
	// Parse and check token
	authHeader := ctx.Get("Authorization")

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", errors.New("authorization token missing")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("authorization token missing")
	}

	return headerParts[1], nil
}

func GenerateLink() string {
	guid := xid.New()
	return guid.String()
}

func AddTimeToCurrentDate(duration time.Duration) time.Time {
	return time.Now().UTC().Add(duration)
}

func RandStr(strSize int, randType string) string {
	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
