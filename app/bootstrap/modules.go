package bootstrap

import (
	"go.uber.org/fx"

	"finapp/api/controllers"
	"finapp/api/middlewares"
	"finapp/api/routes"
	"finapp/lib"
	"finapp/repository"
	"finapp/services"
)

var CommonModules = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
)
