package goods

import "github.com/jackc/pgx/v5/pgxpool"

type Goods struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Goods {
	return &Goods{
		db: db,
	}
}
