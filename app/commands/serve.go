package commands

import (
	"github.com/spf13/cobra"

	"finapp/api/middlewares"
	"finapp/api/routes"
	"finapp/docs"
	"finapp/lib"
)

// ServeCommand test command
type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		middleware middlewares.Middlewares,
		env lib.Env,
		router lib.RequestHandler,
		route routes.Routes,
		logger lib.Logger,
		database lib.Database,
	) {
		middleware.Setup()
		route.Setup()

		// Динамический хост в доке
		docs.SwaggerInfo.Host = env.Host
		// Олег, насрано 👇
		/*if env.ServerPort != "" {
			docs.SwaggerInfo.Host += ":" + env.ServerPort
		}*/

		logger.Info("Running server")
		if env.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + env.ServerPort)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
