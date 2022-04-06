package registry

import (
	"auth-project/src/interface/controller"
	"auth-project/tools"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bun"
)

type registry struct {
	db       *bun.DB
	rdb      *redis.Client
	jwtConf  *tools.JwtConfigurator
}

type Registry interface {
	NewAPIController() controller.APIController
}

func NewRegistry(db *bun.DB,
	rdb *redis.Client,
	jwtConf *tools.JwtConfigurator) Registry {
	return &registry{db, rdb,  jwtConf}
}

func (r *registry) NewAPIController() controller.APIController {
	return controller.APIController{
		Auth:            r.NewAuthController(),
		User:            r.NewUserController(),
	}
}
