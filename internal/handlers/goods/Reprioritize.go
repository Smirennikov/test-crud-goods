package goods

import (
	"context"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

func (h *handlers) Reprioritize(ctx *fiber.Ctx) error {

	goodID := ctx.QueryInt("id")
	projectID := ctx.QueryInt("projectId")

	list, err := h.store.Goods.List(context.TODO(), models.GoodsFilter{GoodID: &goodID, ProjectID: &projectID}, models.GoodsOptions{Limit: 1})
	if err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	good, found := getGood(list)
	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.NotFoundGoodErr)
	}

	type Body struct {
		NewPriority int `json:"newPriority"`
	}
	body := new(Body)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.store.Goods.Reprioritize(context.Background(), good.Priority, body.NewPriority); err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods reprioritize")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	goods, err := h.store.Goods.List(context.Background(), models.GoodsFilter{MinPriority: &body.NewPriority}, models.GoodsOptions{})
	if err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods list")
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	priorities := make([]fiber.Map, 0, len(goods))
	for _, g := range goods {
		priorities = append(priorities, fiber.Map{"id": g.ID, "priority": g.Priority})

		if err := h.updateCache(g); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
		}
		if err := h.logEvent(g.GetLogEvent()); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"priorities": priorities,
	})
}
