package sms

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(sendTo, smsBody string) error {
	accountSid := viper.GetString("twilio_sms.accountSid")
	authToken := viper.GetString("twilio_sms.authToken")

	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(sendTo)
	params.SetFrom(viper.GetString("twilio_sms.phone"))
	params.SetBody(smsBody)

	_, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		return errors.New("error sending sms")
	}

	return nil
}

func CreateSmsBodyVerificationCode(code string) string {
	return "Your " + viper.GetString("project_name") + " verification code is: " + code
}
