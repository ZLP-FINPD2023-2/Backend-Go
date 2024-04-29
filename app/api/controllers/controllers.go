package controllers

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewJWTAuthController),
	fx.Provide(NewUserController),
	fx.Provide(NewGoalController),
	fx.Provide(NewBudgetController),
	fx.Provide(NewTrxController),
	fx.Provide(NewGeneratorController),
)
