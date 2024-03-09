package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type postgresCfg struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Db       string `env:"POSTGRES_DB"`
}

var (
	cfg  *postgresCfg
	once sync.Once
)

func Config() postgresCfg {

	once.Do(func() {
		cfg = &postgresCfg{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
	})
	return *cfg
}

func Connect(logger *zerolog.Logger, host, port, user, password, dbname string) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	logger.Info().Msg(fmt.Sprintf("Postgres pool successfully connected %s:%s", host, port))
	return conn, nil
}
