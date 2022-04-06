package conf

import (
	"auth-project/src/domain/model"
	"log"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	fillErrWithViper()
}

func fillErrWithViper() {
	model.TokenTimeSendErr = "code was sent less than a " + viper.GetDuration("2fa.send_timeout").String() + " ago"
}
