package main

import (
	"flag"
	"log"
	"os"
	"test-crud-goods/internal/config/clickhouse"
	"test-crud-goods/internal/config/nats"
	"test-crud-goods/internal/config/postgres"
	"test-crud-goods/internal/config/redis"
	"test-crud-goods/pkg/closer"
	"time"

	"test-crud-goods/internal/listener"
	"test-crud-goods/internal/server"
	"test-crud-goods/internal/store"
	"test-crud-goods/internal/store/goods"
	"test-crud-goods/internal/store/logs"
	"test-crud-goods/internal/utils/consts"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var (
	envFile string
)

func init() {

	args := os.Args
	if len(args) == 1 {
		log.Fatal(color.RedString("Didn`t pass app mode: dev or prod"))
		return
	}
	mode := args[1]

	var env_mode_file string
	if mode == consts.DEV_MODE {
		log.Println(color.BlueString("run in development mode"))
		os.Setenv("APP_MODE", consts.DEV_MODE)
		env_mode_file = consts.DEV_MODE + ".env"
	}
	if mode == consts.PROD_MODE {
		log.Println(color.BlueString("run in production mode"))
		os.Setenv("APP_MODE", consts.PROD_MODE)
		env_mode_file = consts.PROD_MODE + ".env"
	}

	flag.StringVar(&envFile, "env", env_mode_file, "path to env file")
}

func main() {
	flag.Parse()

	c := closer.New(10 * time.Second)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if err := godotenv.Load(envFile); err != nil {
		logger.Panic().Err(err).Msg("Error loading .env file")
	}

	go func() {
		if err := runLogsListener(c, logger.With().Str("module", "logs-listener").Logger()); err != nil {
			logger.Panic().Err(err).Msg("Error run logs listener")
		}
	}()

	go func() {
		if err := runServer(c, logger.With().Str("module", "server").Logger()); err != nil {
			logger.Panic().Err(err).Msg("Error run server")
		}
	}()

	c.GracefulShutdown()
}

func runServer(closer closer.Closer, logger zerolog.Logger) error {

	pg_pool, err := postgres.Connect(
		closer,
		&logger,
		postgres.Config().Host,
		postgres.Config().Port,
		postgres.Config().User,
		postgres.Config().Password,
		postgres.Config().Db,
	)
	if err != nil {
		return err
	}

	clickhouse_conn, err := clickhouse.Connect(
		closer,
		&logger,
		clickhouse.Config().Host,
		clickhouse.Config().Port,
		clickhouse.Config().Db,
		clickhouse.Config().User,
		clickhouse.Config().Password,
	)
	if err != nil {
		return err
	}

	redis, err := redis.Connect(
		closer,
		&logger,
		redis.Config().Db,
		redis.Config().Host,
		redis.Config().Port,
		redis.Config().Password,
	)
	if err != nil {
		return err
	}

	nats, err := nats.Connect(
		closer,
		&logger,
		nats.Config().URL,
	)
	if err != nil {
		return err
	}

	s := store.New(goods.New(pg_pool), logs.New(clickhouse_conn))

	return server.New(closer, &logger, s, redis, nats).Start()
}

func runLogsListener(closer closer.Closer, logger zerolog.Logger) error {

	pg_pool, err := postgres.Connect(
		closer,
		&logger,
		postgres.Config().Host,
		postgres.Config().Port,
		postgres.Config().User,
		postgres.Config().Password,
		postgres.Config().Db,
	)
	if err != nil {
		return err
	}

	clickhouse_conn, err := clickhouse.Connect(
		closer,
		&logger,
		clickhouse.Config().Host,
		clickhouse.Config().Port,
		clickhouse.Config().Db,
		clickhouse.Config().User,
		clickhouse.Config().Password,
	)
	if err != nil {
		return err
	}

	nats, err := nats.Connect(
		closer,
		&logger,
		nats.Config().URL,
	)
	if err != nil {
		return err
	}

	s := store.New(goods.New(pg_pool), logs.New(clickhouse_conn))

	return listener.New(closer, &logger, s, nats).ListenLogs()
}
