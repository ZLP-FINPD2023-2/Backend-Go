package routes

import (
	"log"

	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewTrxRoutes),
	fx.Provide(NewBudgetRoutes),
	fx.Provide(NewGoalRoutes),
	fx.Provide(NewDocsRoutes),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	authRoutes AuthRoutes,
	userRoutes UserRoutes,
	trxRoutes TrxRoutes,
	budgetRoutes BudgetRoutes,
	goalRoutes GoalRoutes,
	docsRoutes DocsRoutes,
) Routes {
	return Routes{
		authRoutes,
		userRoutes,
		trxRoutes,
		budgetRoutes,
		goalRoutes,
		docsRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	log.Println("Setting up routes")
	for _, route := range r {
		route.Setup()
	}
}
