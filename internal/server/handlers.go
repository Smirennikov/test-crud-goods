package server

import "test-crud-goods/internal/handlers/goods"

func (s *server) configureHandlers() {

	goodsHandlers := goods.New(s.logger, s.store, s.cache, s.nats)

	v1 := s.router.Group("/v1")
	goodsGroup := v1.Group("/goods")

	goodsGroup.Get("/list", s.verifyCache, goodsHandlers.List)
	goodsGroup.Post("/create", goodsHandlers.Create)
	goodsGroup.Patch("/update", goodsHandlers.Update)
	goodsGroup.Patch("/reprioritize", goodsHandlers.Reprioritize)
	goodsGroup.Delete("/remove", goodsHandlers.Remove)
}
