package routes

import (
	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/lib"
)

type GoalRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	controller     controllers.GoalController
	authMiddleware middlewares.JWTAuthMiddleware
}

func (s GoalRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1").Use(s.authMiddleware.Handler())
	{
		root.GET("/goal", s.controller.List)
		root.POST("/goal", s.controller.Create)
	}
}

func NewGoalRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	controller controllers.GoalController,
	authMiddleware middlewares.JWTAuthMiddleware,
) GoalRoutes {
	return GoalRoutes{
		logger:         logger,
		handler:        handler,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}
