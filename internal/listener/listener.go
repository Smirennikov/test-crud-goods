package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/store"
	"test-crud-goods/internal/utils"
	"test-crud-goods/internal/utils/consts"
	"test-crud-goods/pkg/closer"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type listener struct {
	logger *zerolog.Logger
	store  *store.Store
	nats   *nats.Conn
	closer closer.Closer
}

func New(closer closer.Closer, logger *zerolog.Logger, store *store.Store, nats *nats.Conn) *listener {
	return &listener{
		logger: logger,
		store:  store,
		nats:   nats,
		closer: closer,
	}
}

func (srv *listener) ListenLogs() error {
	sub, err := srv.nats.SubscribeSync(consts.LogsQueue)
	if err != nil {
		return err
	}
	srv.logger.Info().Msg(fmt.Sprintf("Listen nats:%s", consts.LogsQueue))

	eventChan := make(chan struct{})
	var batchLogs []models.GoodLogEvent

	sendBatch := func() (err error) {
		if len(batchLogs) == 0 {
			return
		}

		batch := make([]models.GoodLogEvent, 0, len(batchLogs))
		copy(batch, batchLogs)
		batchLogs = nil

		if err = srv.store.Logs.SendBatch(context.TODO(), batchLogs); err != nil {
			srv.logger.Error().Err(err).Msg("batch didn`t send")
			return
		}
		srv.logger.Info().Msg("batch successfully sent")
		return
	}

	srv.closer.Add(sendBatch)
	go utils.Debounce(time.Second*15, eventChan, func() {
		sendBatch()
	})

	for {
		msg, err := sub.NextMsg(time.Second * 10)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			srv.logger.Fatal().Err(err).Msg("batch send")
		}
		eventChan <- struct{}{}

		var logEvent models.GoodLogEvent
		if err := json.Unmarshal(msg.Data, &logEvent); err != nil {
			srv.logger.Error().Err(err).Msg("decoding json")
			continue
		}

		batchLogs = append(batchLogs, logEvent)
	}
}
