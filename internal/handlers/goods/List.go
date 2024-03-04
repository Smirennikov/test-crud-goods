package goods

import (
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

const (
	defaultLimit  = 10
	defaultOffset = 1
)

func (h *handlers) List(ctx *fiber.Ctx) error {

	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	if _, ok := ctx.Queries()["limit"]; !ok {
		limit = defaultLimit
	}
	if _, ok := ctx.Queries()["offset"]; !ok {
		offset = defaultOffset
	}

	g, gCtx := errgroup.WithContext(ctx.Context())

	list := make([]models.Good, 0)
	var meta models.GoodsMeta

	g.Go(func() (err error) {
		meta, err = h.store.Goods.ListMeta(gCtx)
		if err != nil {
			h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list meta")
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
		}

		return
	})
	g.Go(func() (err error) {
		list, err = h.store.Goods.List(gCtx, models.GoodsFilter{}, models.GoodsOptions{Limit: limit, Offset: offset})
		if err != nil {
			h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list")
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
		}
		return
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"total":   meta.Total,
			"removed": meta.Removed,
			"limit":   limit,
			"offset":  offset,
		},
		"goods": list,
	})
}
