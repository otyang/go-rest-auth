package tools

import (

	"errors"
	"github.com/rs/xid"
	"math"
	"time"
)


func GenerateCode() (string, error) {

	guid := xid.New()
	if len(guid.String()) != 20 {
		err := errors.New("unique code generation error")
		return "", err
	}

	code := guid.String()[4:8] + guid.String()[16:20]
	return code, nil
}

func GenerateLink() string {
	guid := xid.New()
	return guid.String()
}

func AddDaysToCurrentDate(days int) int64 {
	return time.Now().UTC().AddDate(0, 0, days).Unix()
}

func GetDaysDifferenceFromCurrentDate(expirationTime int64) int {
	expirationDate := time.Unix(expirationTime, 0)
	diff := expirationDate.Sub(time.Now().UTC())
	return int(math.Round(diff.Hours() / 24))
}
