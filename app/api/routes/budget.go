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
		root.GET("/budget/:id", s.controller.Get)
		root.GET("/budget", s.controller.List)
		root.POST("/budget", s.controller.Post)
		root.DELETE("/budget/:id", s.controller.Delete)
		root.PATCH("/budget", s.controller.Patch)
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
