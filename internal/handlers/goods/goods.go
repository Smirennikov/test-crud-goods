package goods

import (
	"encoding/json"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/store"
	"test-crud-goods/internal/utils/consts"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type handlers struct {
	logger *zerolog.Logger
	store  *store.Store
	cache  *redis.Client
	nats   *nats.Conn
}

func New(logger *zerolog.Logger, store *store.Store, cache *redis.Client, nats *nats.Conn) *handlers {
	return &handlers{
		logger: logger,
		store:  store,
		cache:  cache,
		nats:   nats,
	}
}

func getGood(list []models.Good) (*models.Good, bool) {
	if len(list) == 0 {
		return nil, false
	}
	return &list[0], true
}

func (h *handlers) logEvent(lg models.GoodLogEvent) (err error) {
	bytes, err := json.Marshal(lg)
	if err != nil {
		return err
	}
	return h.nats.Publish(consts.LogsQueue, bytes)
}
