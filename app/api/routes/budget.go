package routes

import (
	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/lib"
)

type BudgetRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	controller     controllers.BudgetController
	authMiddleware middlewares.JWTAuthMiddleware
}

func (s BudgetRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1").Use(s.authMiddleware.Handler())
	{
		root.GET("/budget", s.controller.Get)
		root.POST("/budget", s.controller.Post)
	}
}

func NewBudgetRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	controller controllers.BudgetController,
	authMiddleware middlewares.JWTAuthMiddleware,
) BudgetRoutes {
	return BudgetRoutes{
		logger:         logger,
		handler:        handler,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}
