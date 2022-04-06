package registry

import (
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/interface/controller"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bun"
)

type registry struct {
	db      *bun.DB
	rdb     *redis.Client
	jwtConf *authentication.JwtConfigurator
}

type Registry interface {
	NewAPIController() controller.APIController
}

func NewRegistry(db *bun.DB,
	rdb *redis.Client,
	jwtConf *authentication.JwtConfigurator) Registry {
	return &registry{db, rdb, jwtConf}
}

func (r *registry) NewAPIController() controller.APIController {
	return controller.APIController{
		Auth:          r.NewAuthController(),
		QrCodeAuth:    r.NewQrCodeAuthController(),
		TwoFactorAuth: r.NewTwoFactorAuthController(),
		User:          r.NewUserController(),
		Token:         r.NewTokenController(),
	}
}
