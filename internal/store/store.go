package store

import "test-crud-goods/internal/store/goods"

type Store struct {
	Goods *goods.Goods
}

func New(goods *goods.Goods) *Store {
	return &Store{
		Goods: goods,
	}
}
