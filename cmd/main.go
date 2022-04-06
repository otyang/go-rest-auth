package main

import (
	"auth-project/conf"
	"auth-project/src/infrastructure/delivery/http"
	"auth-project/src/registry"
	"auth-project/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"auth-project/src/infrastructure/storage"
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
	jwtConf := tools.NewJwtConfigurator(
		viper.GetDuration("jwt.at_min_lifetime"),
		viper.GetDuration("jwt.rt_min_lifetime"))

	// Init a new fiber application
	app := fiber.New()

	// Init a new registry
	r := registry.NewRegistry(db, rdb, jwtConf)

	app = http.NewRouter(app, r.NewAPIController())

	app.Name("auth-project")

	err := app.Listen(viper.GetString("http.port"))
	if err != nil {
		panic(err)
	}
}
