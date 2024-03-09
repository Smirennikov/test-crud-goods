package server

import (
	"os"
	"test-crud-goods/internal/store"
	"test-crud-goods/pkg/closer"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type server struct {
	logger *zerolog.Logger
	router fiber.Router
	store  *store.Store
	cache  *redis.Client
	nats   *nats.Conn
	closer closer.Closer
}

func New(closer closer.Closer, logger *zerolog.Logger, store *store.Store, cache *redis.Client, nats *nats.Conn) *server {
	return &server{
		logger: logger,
		store:  store,
		cache:  cache,
		nats:   nats,
		closer: closer,
	}
}

func (s *server) Start() error {
	app := fiber.New()

	s.closer.Add(app.Shutdown)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
	}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: s.logger,
	}))

	s.router = app.Group("/api")

	s.configureHandlers()

	return app.Listen(":" + os.Getenv("APP_PORT"))
}
