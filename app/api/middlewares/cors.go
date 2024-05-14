package middlewares

import (
	"strings"

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
	m.logger.Info("AllowedHeaders: " + strings.Join(strings.Split(m.env.CORSAllowedHeaders, " "), ", "))
	m.logger.Info("AllowedOrigins: " + strings.Join(strings.Split(m.env.CORSAllowedOrigins, " "), ", "))
	m.logger.Info("AllowedMethods: " + strings.Join(strings.Split(m.env.CORSAllowedMethods, " "), ", "))

	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   strings.Split(m.env.CORSAllowedHeaders, " "),
		AllowedOrigins:   strings.Split(m.env.CORSAllowedOrigins, " "),
		AllowedMethods:   strings.Split(m.env.CORSAllowedMethods, " "),
		Debug:            false,
	}))
}
