package routes

import (
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "finapp/docs"
	"finapp/lib"
)

// DocsRoutes struct
type DocsRoutes struct {
	logger  lib.Logger
	handler lib.RequestHandler
}

// Setup user routes
func (s DocsRoutes) Setup() {
	root := s.handler.Gin.Group("/api/v1")
	{
		root.GET("/schemes/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// NewDocsRoutes creates new user controller
func NewDocsRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
) DocsRoutes {
	return DocsRoutes{
		logger:  logger,
		handler: handler,
	}
}
