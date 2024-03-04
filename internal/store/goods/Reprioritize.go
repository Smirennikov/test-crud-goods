package goods

import (
	"context"
	"test-crud-goods/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *Goods) Reprioritize(ctx context.Context, curPriority, newPriority int) error {

	rows, err := r.db.Query(ctx, "SELECT g.id, g.project_id FROM goods g WHERE g.priority >= $1 ORDER BY priority", curPriority)
	if err != nil {
		return err
	}
	defer rows.Close()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for rows.Next() {
		var good models.Good
		if err := rows.Scan(&good.ID, &good.ProjectID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "UPDATE goods SET priority = $3 WHERE id = $1 AND project_id = $2", good.ID, good.ProjectID, newPriority); err != nil {
			return err
		}

		newPriority++
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
