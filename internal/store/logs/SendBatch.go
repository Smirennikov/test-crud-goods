package logs

import (
	"context"
	"test-crud-goods/internal/models"
)

func (r *logs) SendBatch(ctx context.Context, batchLogs []models.GoodLogEvent) error {
	batch, err := r.db.PrepareBatch(ctx,
		"INSERT INTO goods_logs(Id, ProjectId, Name, Description, Priority, Removed, EventTime)",
	)
	if err != nil {
		return err
	}

	for _, logEvent := range batchLogs {
		if err := batch.Append(logEvent.ID, logEvent.ProjectID, logEvent.Name, logEvent.Description, logEvent.Priority, logEvent.Removed, logEvent.EventTime); err != nil {
			return err
		}
	}
	return batch.Send()
}
