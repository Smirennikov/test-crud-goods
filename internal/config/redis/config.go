package redis

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type redisCfg struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	Db       int    `env:"REDIS_DB"`
}

var (
	cfg  *redisCfg
	once sync.Once
)

func Config() redisCfg {

	once.Do(func() {
		cfg = &redisCfg{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
	})
	return *cfg
}

func Connect(logger *zerolog.Logger, db int, host, port, password string) (*redis.Client, error) {
	addr := net.JoinHostPort(host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	logger.Info().Msg(fmt.Sprintf("Redis successfully connected %s", addr))
	return rdb, nil
}
