package routes

import (
	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/lib"
)

type GeneratorRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	controller     controllers.GeneratorController
	authMiddleware middlewares.JWTAuthMiddleware
}

func (s GeneratorRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1").Use(s.authMiddleware.Handler())
	{
		root.GET("/trx/generator/:id", s.controller.Get)
		root.GET("/trx/generator", s.controller.List)
		root.POST("/trx/generator", s.controller.Store)
		root.DELETE("/trx/generator/:id", s.controller.Delete)
		root.PATCH("/trx/generator/:id", s.controller.Update)
	}
}

func NewGeneratorRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	controller controllers.GeneratorController,
	authMiddleware middlewares.JWTAuthMiddleware,
) GeneratorRoutes {
	return GeneratorRoutes{
		logger:         logger,
		handler:        handler,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}
