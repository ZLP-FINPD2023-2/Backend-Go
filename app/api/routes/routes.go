package routes

import (
	"log"

	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewDocsRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewGoalRoutes),
	fx.Provide(NewBudgetRoutes),
	fx.Provide(NewTrxRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	docsRoutes DocsRoutes,
	authRoutes AuthRoutes,
	userRoutes UserRoutes,
	goalRoutes GoalRoutes,
	budgetRoutes BudgetRoutes,
	trxRoutes TrxRoutes,
) Routes {
	return Routes{
		docsRoutes,
		authRoutes,
		userRoutes,
		goalRoutes,
		budgetRoutes,
		trxRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	log.Println("Setting up routes")
	for _, route := range r {
		route.Setup()
	}
}
