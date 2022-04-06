package email

import (
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func SendEmail(name, address, subject, plainTextContent, htmlContent string) error {
	apiKey := viper.GetString("sendgrid.api_key")
	fromName := viper.GetString("sendgrid.from_name")
	fromAddress := viper.GetString("sendgrid.from_address")

	from := mail.NewEmail(fromName, fromAddress)
	to := mail.NewEmail(name, address)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	_, err := sendgrid.NewSendClient(apiKey).Send(message)
	if err != nil {
		return errors.New("error sending email")
	}

	return nil
}

func CreateEmailBodyVerificationCode(code string) (plainTextContent, htmlContent string) {
	plainTextContent = "Your " + viper.GetString("project_name") + " verification code is: " + code
	htmlContent = "Your " + viper.GetString("project_name") + " verification code is: " + code
	return
}
