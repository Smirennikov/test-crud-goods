package goods

import "github.com/jackc/pgx/v5/pgxpool"

type goods struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *goods {
	return &goods{
		db: db,
	}
}
