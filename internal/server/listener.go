package server

import (
	"context"
	"encoding/json"
	"log"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils"
	"test-crud-goods/internal/utils/consts"
	"time"

	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
)

func (s *server) listener() {
	sub, err := s.nats.SubscribeSync(consts.LogsQueue)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(color.GreenString("Listen nats:%s", consts.LogsQueue))

	eventChan := make(chan string)
	var batchLogs []models.GoodLogEvent

	go utils.Debounce(time.Second*15, eventChan, func(name string) {

		batch := make([]models.GoodLogEvent, 0, len(batchLogs))
		copy(batch, batchLogs)
		batchLogs = nil

		if err := s.sendBatchToClickHouse(batchLogs); err != nil {
			s.logger.Error().Err(err).Msg("batch unsuccessfully send")
			return
		}
		s.logger.Info().Msg("batch successfully sent")
	})

	for {
		msg, err := sub.NextMsg(time.Second * 10)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			s.logger.Fatal().Err(err).Msg("batch send")
		}
		eventChan <- "wait"

		var logEvent models.GoodLogEvent
		if err := json.Unmarshal(msg.Data, &logEvent); err != nil {
			s.logger.Error().Err(err).Msg("decoding json")
			continue
		}

		batchLogs = append(batchLogs, logEvent)
	}
}

func (s *server) sendBatchToClickHouse(batchLogs []models.GoodLogEvent) error {
	batch, err := s.clickhouse.PrepareBatch(context.Background(),
		"INSERT INTO goods_logs(Id, ProjectId, Name, Description, Priority, Removed, EventTime)",
	)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("batch prepare")
	}

	for _, logEvent := range batchLogs {
		if err := batch.Append(logEvent.ID, logEvent.ProjectID, logEvent.Name, logEvent.Description, logEvent.Priority, logEvent.Removed, logEvent.EventTime); err != nil {
			s.logger.Error().Err(err).Msg("batch append")
			continue
		}
	}
	return batch.Send()
}
