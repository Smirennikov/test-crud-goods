package goods

import (
	"context"
	"test-crud-goods/internal/models"
)

func (r *Goods) ListMeta(ctx context.Context) (meta models.GoodsMeta, err error) {

	row := r.db.QueryRow(ctx, `
		SELECT
			COUNT(id) total,
			COUNT(id) FILTER (WHERE removed = true) removed
		FROM goods
	`)

	return meta, row.Scan(&meta.Total, &meta.Removed)
}
