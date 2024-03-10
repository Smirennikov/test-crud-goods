package goods

import (
	"encoding/json"
	"fmt"
	"strings"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

const (
	defaultLimit             = 10
	defaultOffset            = 1
	cacheGoodsListKey        = "/goods/list"
	cacheGoodsListTotalKey   = cacheGoodsListKey + "/meta/total"
	cacheGoodsListRemovedKey = cacheGoodsListKey + "/meta/removed"
)

func (h *handlers) List(ctx *fiber.Ctx) error {

	opts := getListOptions(ctx)

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
		list, err = h.store.Goods.List(gCtx, models.GoodsFilter{}, opts)
		if err != nil {
			h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list")
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
		}
		return
	})

	if err := g.Wait(); err != nil {
		return err
	}

	goodsIDs := make([]string, 0, len(list))
	for _, g := range list {
		goodsIDs = append(goodsIDs, g.Key())
		bytes, err := json.Marshal(g)
		if err != nil {
			return err
		}
		h.cache.Set(ctx.Context(), cacheGoodsListKey+"?"+g.Key(), bytes, time.Minute)
	}

	h.cache.Set(
		ctx.Context(),
		cacheGoodsListTotalKey,
		meta.Total,
		time.Minute,
	)
	h.cache.Set(
		ctx.Context(),
		cacheGoodsListRemovedKey,
		meta.Removed,
		time.Minute,
	)
	h.cache.Set(
		ctx.Context(),
		fmt.Sprintf("%s?limit=%d&offset=%d", cacheGoodsListKey, opts.Limit, opts.Offset),
		strings.Join(goodsIDs, ","),
		time.Minute,
	)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"total":   meta.Total,
			"removed": meta.Removed,
			"limit":   opts.Limit,
			"offset":  opts.Offset,
		},
		"goods": list,
	})
}

func (h *handlers) CachedList(ctx *fiber.Ctx) error {

	opts := getListOptions(ctx)

	total, err := h.cache.Get(ctx.Context(), cacheGoodsListTotalKey).Int()
	if err != nil {
		return ctx.Next()
	}
	removed, err := h.cache.Get(ctx.Context(), cacheGoodsListRemovedKey).Int()
	if err != nil {
		return ctx.Next()
	}
	goodIDs, err := h.cache.Get(ctx.Context(), fmt.Sprintf("%s?limit=%d&offset=%d", cacheGoodsListKey, opts.Limit, opts.Offset)).Result()
	if err != nil {
		return ctx.Next()
	}

	var list []models.Good
	for _, goodKey := range strings.Split(goodIDs, ",") {
		result, err := h.cache.Get(ctx.Context(), cacheGoodsListKey+"?"+goodKey).Bytes()
		if err != nil {
			return ctx.Next()
		}
		var good models.Good
		if err := json.Unmarshal(result, &good); err != nil {
			return ctx.Next()
		}
		list = append(list, good)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"meta": fiber.Map{
			"total":   total,
			"removed": removed,
			"limit":   opts.Limit,
			"offset":  opts.Offset,
		},
		"goods": list,
	})
}

func getListOptions(ctx *fiber.Ctx) models.GoodsOptions {

	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	if _, ok := ctx.Queries()["limit"]; !ok {
		limit = defaultLimit
	}
	if _, ok := ctx.Queries()["offset"]; !ok {
		offset = defaultOffset
	}

	return models.GoodsOptions{
		Limit:  limit,
		Offset: offset,
	}
}
