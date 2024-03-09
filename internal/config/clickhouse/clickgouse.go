package clickhouse

import (
	"context"
	"fmt"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/rs/zerolog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/caarlos0/env/v10"
)

type clickhouseCfg struct {
	Host     string `env:"CLICKHOUSE_HOST"`
	Port     string `env:"CLICKHOUSE_PORT"`
	User     string `env:"CLICKHOUSE_USER"`
	Password string `env:"CLICKHOUSE_PASSWORD"`
	Db       string `env:"CLICKHOUSE_DB"`
}

var (
	cfg  *clickhouseCfg
	once sync.Once
)

func Config() clickhouseCfg {

	once.Do(func() {
		cfg = &clickhouseCfg{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
	})
	return *cfg
}

func Connect(logger *zerolog.Logger, host, port, database, username, password string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{host + ":" + port},
			// Auth: clickhouse.Auth{
			// 	Database: database,
			// 	Username: username,
			// 	Password: password,
			// },
			Debugf: func(format string, v ...interface{}) {
				logger.Debug().Msg(fmt.Sprintf(format, v))
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logger.Info().Msg(fmt.Sprintf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
		}
		return nil, err
	}

	logger.Info().Msg(fmt.Sprintf("Clickhouse successfully connected %s:%s", host, port))
	return conn, nil
}
