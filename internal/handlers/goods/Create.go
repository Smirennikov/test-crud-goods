package goods

import (
	"context"
	"encoding/json"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *handlers) Create(ctx *fiber.Ctx) error {

	projectID := ctx.QueryInt("projectId")

	type Body struct {
		Name string `json:"name"`
	}
	body := new(Body)

	if err := ctx.BodyParser(body); err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Bytes("body", ctx.Body()).Msg("parse body")

		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	newGood := models.Good{ProjectID: projectID, Name: body.Name}
	id, err := h.store.Goods.Create(context.TODO(), newGood)
	if err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Any("newGood", newGood).Msg("goods create")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	list, err := h.store.Goods.List(context.TODO(), models.GoodsFilter{GoodID: id, ProjectID: &projectID}, models.GoodsOptions{Limit: 1})
	if err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Any("goodID", id).Msg("goods list")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}
	good, _ := getGood(list)

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

	return ctx.Status(fiber.StatusOK).JSON(good)
}
