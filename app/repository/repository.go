package repository

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewTrxRepository),
	fx.Provide(NewBudgetRepository),
	fx.Provide(NewGoalRepository),
	fx.Provide(NewGeneratorRepository),
)
