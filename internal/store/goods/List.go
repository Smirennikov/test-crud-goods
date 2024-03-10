package goods

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"test-crud-goods/internal/models"
)

func (r *goods) List(ctx context.Context, filter models.GoodsFilter, options models.GoodsOptions) (goods []models.Good, err error) {

	query := `
		SELECT g.id, g.project_id, g.name, g.description, g.priority, g.removed, g.created_at FROM goods g
	`

	var args []interface{}
	var whereAnd []string

	if filter.GoodID != nil {
		args = append(args, filter.GoodID)
		whereAnd = append(whereAnd, fmt.Sprintf("g.id = $%d", len(args)))
	}
	if filter.ProjectID != nil {
		args = append(args, filter.ProjectID)
		whereAnd = append(whereAnd, fmt.Sprintf("g.project_id = $%d", len(args)))
	}
	if filter.MinPriority != nil {
		args = append(args, filter.MinPriority)
		whereAnd = append(whereAnd, fmt.Sprintf("g.priority >= $%d", len(args)))
	}

	if len(whereAnd) != 0 {
		query += "WHERE " + strings.Join(whereAnd, " AND ")
	}

	if options.Offset != 0 {
		query += "OFFSET " + strconv.Itoa(options.Offset)
	}
	if options.Limit != 0 {
		query += "LIMIT " + strconv.Itoa(options.Limit)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var good models.Good
		var description sql.NullString

		if err := rows.Scan(&good.ID, &good.ProjectID, &good.Name, &description, &good.Priority, &good.Removed, &good.CreatedAt); err != nil {
			return nil, err
		}

		if description.Valid {
			good.Description = description.String
		}
		goods = append(goods, good)
	}

	return goods, rows.Err()
}
