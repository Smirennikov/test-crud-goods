package store

import (
	"context"
	"test-crud-goods/internal/models"
)

type Store struct {
	Goods goodsRepository
	Logs  logsRepository
}

func New(goods goodsRepository, logs logsRepository) *Store {
	return &Store{
		Goods: goods,
		Logs:  logs,
	}
}

type logsRepository interface {
	SendBatch(ctx context.Context, batchLogs []models.GoodLogEvent) error
}

type goodsRepository interface {
	Create(ctx context.Context, good models.Good) (id *int, err error)

	List(ctx context.Context, filter models.GoodsFilter, options models.GoodsOptions) (goods []models.Good, err error)
	ListMeta(ctx context.Context) (meta models.GoodsMeta, err error)

	Update(ctx context.Context, good models.Good) error
	Reprioritize(ctx context.Context, curPriority, newPriority int) error
}
