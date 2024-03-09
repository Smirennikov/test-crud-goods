package nats

import (
	"fmt"
	"sync"
	"test-crud-goods/pkg/closer"

	"github.com/caarlos0/env/v10"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type natsCfg struct {
	URL string `env:"NATS_URL"`
}

var (
	cfg  *natsCfg
	once sync.Once
)

func Config() natsCfg {

	once.Do(func() {
		cfg = &natsCfg{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
	})
	return *cfg
}

func Connect(closer closer.Closer, logger *zerolog.Logger, url string) (*nats.Conn, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	closer.Add(func() (err error) {
		conn.Close()
		return
	})

	logger.Info().Msg(fmt.Sprintf("Nats successfully connected %s", url))
	return conn, nil
}
