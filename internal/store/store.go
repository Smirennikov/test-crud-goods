package store

import (
	"test-crud-goods/internal/store/goods"
	"test-crud-goods/internal/store/logs"
)

type Store struct {
	Goods *goods.Goods
	Logs  *logs.Logs
}

func New(goods *goods.Goods, logs *logs.Logs) *Store {
	return &Store{
		Goods: goods,
		Logs:  logs,
	}
}
