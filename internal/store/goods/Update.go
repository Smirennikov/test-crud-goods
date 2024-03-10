package goods

import (
	"context"
	"test-crud-goods/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *goods) Update(ctx context.Context, good models.Good) error {

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT id FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE", good.ID, good.ProjectID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		"UPDATE goods SET name = $3, description= $4, removed = $5  WHERE id = $1 AND project_id = $2",
		good.ID, good.ProjectID, good.Name, good.Description, good.Removed); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
