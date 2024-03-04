package server

import (
	"test-crud-goods/internal/store"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type server struct {
	logger     *zerolog.Logger
	router     fiber.Router
	store      *store.Store
	cache      *redis.Client
	clickhouse driver.Conn
	nats       *nats.Conn
}

func New(logger *zerolog.Logger, router fiber.Router, store *store.Store, cache *redis.Client, clickhouse driver.Conn, nats *nats.Conn) *server {
	server := server{
		logger:     logger,
		router:     router,
		store:      store,
		cache:      cache,
		clickhouse: clickhouse,
		nats:       nats,
	}
	server.configureHandlers()
	go server.listener()

	return &server
}
