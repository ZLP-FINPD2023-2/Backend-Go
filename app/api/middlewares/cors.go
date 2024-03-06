package middlewares

import (
	cors "github.com/rs/cors/wrapper/gin"

	"finapp/lib"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	logger  lib.Logger
	handler lib.RequestHandler
	env     lib.Env
}

func NewCorsMiddleware(logger lib.Logger, handler lib.RequestHandler, env lib.Env) CorsMiddleware {
	return CorsMiddleware{
		logger:  logger,
		handler: handler,
		env:     env,
	}
}
func (m CorsMiddleware) Setup() {
	m.logger.Info("Setting up cors middleware")

	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{m.env.CORSAllowedHeaders},
		AllowedOrigins:   []string{m.env.CORSAllowedOrigins},
		AllowedMethods:   []string{m.env.CORSAllowedMethods},
		Debug:            false,
	}))
}
