package main

import (
	"auth-project/conf"
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/infrastructure/delivery/http"
	"auth-project/src/registry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"

	"auth-project/src/infrastructure/storage"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Init configs
	conf.InitConfig()

	// Get environment
	env := viper.GetString("env")

	// Init Postgres connection
	db := storage.InitPostgres()
	defer db.Close()

	// Init Redis connection
	rdb := storage.InitRedis(env)
	defer rdb.Close()

	// Init a new jwt configurator
	jwtConf := authentication.NewJwtConfigurator(
		viper.GetDuration("jwt.access_token_min_lifetime"),
		viper.GetDuration("jwt.refresh_token_min_lifetime"),
		viper.GetDuration("jwt.two_factor_auth_token_min_lifetime"),
		viper.GetDuration("jwt.reset_password_token_min_lifetime"))

	// Init a new fiber application
	app := fiber.New()

	app.Use(recover.New())

	if env == "local" {
		app.Use(logger.New())
	}

	// Init a new registry
	r := registry.NewRegistry(db, rdb, jwtConf)

	app = http.NewRouter(app, r.NewAPIController())

	app.Name(viper.GetString("project_name"))

	err := app.Listen(viper.GetString("http.port"))
	if err != nil {
		panic(err)
	}
}
