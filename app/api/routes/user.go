package routes

import (
	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/lib"
)

// UserRoutes struct
type UserRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	authMiddleware middlewares.JWTAuthMiddleware
	userController controllers.UserController
}

// Setup user routes
func (s UserRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1").Use(s.authMiddleware.Handler())
	{
		root.GET("/user", s.userController.Get)
		root.PATCH("/user", s.userController.Update)
		root.DELETE("/user", s.userController.Delete)
	}
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	authMiddleware middlewares.JWTAuthMiddleware,
	userController controllers.UserController,
) UserRoutes {
	return UserRoutes{
		logger:         logger,
		handler:        handler,
		authMiddleware: authMiddleware,
		userController: userController,
	}
}
