package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (s *server) verifyCache(c *fiber.Ctx) error {
	val, err := s.cache.Get(context.Background(), c.OriginalURL()).Result()
	if err != nil {
		return c.Next()
	}
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).SendString(val)
}
