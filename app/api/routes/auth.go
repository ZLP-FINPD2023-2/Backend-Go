package routes

import (
	"finapp/api/controllers"
	"finapp/lib"
)

// AuthRoutes struct
type AuthRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	authController controllers.JWTAuthController
}

// Setup auth routes
func (s AuthRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1")
	{
		auth := root.Group("/auth")
		{
			auth.POST("/login", s.authController.Login)
			auth.POST("/register", s.authController.Register)
		}
	}
}

// NewAuthRoutes creates new auth controller
func NewAuthRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	authController controllers.JWTAuthController,
) AuthRoutes {
	return AuthRoutes{
		logger:         logger,
		handler:        handler,
		authController: authController,
	}
}
