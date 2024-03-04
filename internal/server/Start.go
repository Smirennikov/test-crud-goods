package server

import (
	"os"
	"test-crud-goods/internal/store"
	"test-crud-goods/internal/store/goods"

	"github.com/gofiber/contrib/fiberzerolog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

func Start() error {

	pg_pool, err := pg_connect(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		return err
	}
	defer pg_pool.Close()

	clickhouse, err := clickhouse_connect(
		os.Getenv("CLICKHOUSE_HOST"),
		os.Getenv("CLICKHOUSE_PORT"),
		os.Getenv("CLICKHOUSE_DB"),
		os.Getenv("CLICKHOUSE_USER"),
		os.Getenv("CLICKHOUSE_PASSWORD"),
	)
	if err != nil {
		return err
	}
	defer clickhouse.Close()

	redis, err := redis_connect(
		0,
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASSWORD"),
	)
	if err != nil {
		return err
	}
	defer redis.Close()

	nats, err := nats_connect()
	if err != nil {
		return err
	}
	defer nats.Close()

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
	}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	s := store.New(goods.New(pg_pool))
	New(&logger, app.Group("/api"), s, redis, clickhouse, nats)

	return app.Listen(":" + os.Getenv("APP_PORT"))
}
