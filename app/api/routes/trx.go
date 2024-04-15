package routes

import (
	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/lib"
)

type TrxRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	authMiddleware middlewares.JWTAuthMiddleware
	trxController  controllers.TrxController
}

func (s TrxRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1").Use(s.authMiddleware.Handler())
	{
		root.GET("/trx", s.trxController.List)
		root.POST("/trx", s.trxController.Post)
		root.PATCH("/trx/:id", s.trxController.Patch)
		root.DELETE("/trx/:id", s.trxController.Delete)
	}
}

func NewTrxRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	authMiddleware middlewares.JWTAuthMiddleware,
	trxController controllers.TrxController,
) TrxRoutes {
	return TrxRoutes{
		logger:         logger,
		handler:        handler,
		authMiddleware: authMiddleware,
		trxController:  trxController,
	}
}
