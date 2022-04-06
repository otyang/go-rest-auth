package storage

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"

	"github.com/spf13/viper"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitPostgres() *bun.DB {

	db := bun.NewDB(sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithAddr(viper.GetString("db.host")+viper.GetString("db.port")),
			pgdriver.WithUser(viper.GetString("db.user")),
			pgdriver.WithPassword(viper.GetString("db.pass")),
			pgdriver.WithDatabase(viper.GetString("db.name")),
			pgdriver.WithInsecure(true),
		),
	), pgdialect.New())

	//db.AddQueryHook(bundebug.NewQueryHook(
	//	bundebug.WithVerbose(true),
	//	bundebug.FromEnv("BUNDEBUG"),
	//))

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func InitRedis(env string) *redis.Client {

	//Initializing redis
	dsn := viper.GetString("rdb.host") + viper.GetString("rdb.port")
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn, // use default Addr
		Password: "",  // no password set
		DB:       0,   // use default DB
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	if env != "local" {
		rdb.FlushAll(ctx)
	}

	return rdb
}
