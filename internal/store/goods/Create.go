package goods

import (
	"context"
	"test-crud-goods/internal/models"
)

func (r *Goods) Create(ctx context.Context, good models.Good) (id *int, err error) {

	row := r.db.QueryRow(ctx, `
		WITH goods_priority AS (
			SELECT COALESCE(MAX(priority), 0) AS max FROM goods
		)
		INSERT INTO goods(project_id, name, priority) 
		VALUES($1, $2, (SELECT goods_priority.max+1 FROM goods_priority))
		RETURNING id
	`, good.ProjectID, good.Name)

	return id, row.Scan(&id)
}
