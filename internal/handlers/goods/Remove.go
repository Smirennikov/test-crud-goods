package goods

import (
	"context"
	"encoding/json"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *handlers) Remove(ctx *fiber.Ctx) error {

	goodID := ctx.QueryInt("id")
	projectID := ctx.QueryInt("projectId")

	list, err := h.store.Goods.List(context.TODO(),
		models.GoodsFilter{GoodID: &goodID, ProjectID: &projectID}, models.GoodsOptions{Limit: 1})
	if err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	good, found := getGood(list)
	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.NotFoundGoodErr)
	}

	good.Removed = true

	if err := h.store.Goods.Update(context.TODO(), *good); err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods update")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	jsonGood, err := json.Marshal(good)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}
	if err := h.cache.Set(ctx.Context(), ctx.OriginalURL(), jsonGood, time.Minute).Err(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	if err := h.logEvent(good.GetLogEvent()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         goodID,
		"campaignId": projectID,
		"removed":    true,
	})
}
