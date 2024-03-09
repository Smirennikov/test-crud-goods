package goods

import (
	"context"
	"test-crud-goods/internal/models"
	"test-crud-goods/internal/utils/errors"

	"github.com/gofiber/fiber/v2"
)

func (h *handlers) Update(ctx *fiber.Ctx) error {

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

	body := new(models.UpdateGoodBody)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validate(*body, *good); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	good.Name = body.Name
	good.Description = body.Description

	if err := h.store.Goods.Update(context.TODO(), *good); err != nil {
		h.logger.Error().Err(err).Str("url", ctx.OriginalURL()).Msg("goods update")

		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	if err := h.updateCache(*good); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}
	if err := h.logEvent(good.GetLogEvent()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.TryAgainErr)
	}

	return ctx.Status(fiber.StatusOK).JSON(good)
}

func validate(body models.UpdateGoodBody, good models.Good) error {
	if body.Name == "" {
		return errors.NotUpdatedGoodErr.SetDetails(map[string]interface{}{
			"reason": "name - обязательное поле",
		})
	}
	if good.Name == body.Name && good.Description == body.Description {
		return errors.NotUpdatedGoodErr.SetDetails(map[string]interface{}{
			"reason": "изменяемые данные должны отличаться",
		})
	}
	return nil
}
