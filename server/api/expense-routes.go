package api

import (
	"github.com/davinder1436/fingenie/internal/handlers/expense"
	"github.com/davinder1436/fingenie/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupExpenseRoutes(app *fiber.App, db *gorm.DB) {
	expenseHandler := expense.NewHandler(db)

	// Get the v1 group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Expense routes
	expenses := v1.Group("/expenses")
	expenses.Use(middleware.AuthMiddleware())
	expenses.Post("/", expenseHandler.CreateExpense)
	expenses.Get("/", expenseHandler.ListExpenses) // Added list endpoint
	expenses.Get("/:id", expenseHandler.GetExpense)
	expenses.Put("/:id", expenseHandler.UpdateExpense)
	expenses.Delete("/:id", expenseHandler.DeleteExpense)

	// Split Expense routes
	splitExpenses := v1.Group("/split-expenses")
	splitExpenses.Use(middleware.AuthMiddleware())
	splitExpenses.Post("/", expenseHandler.CreateSplitExpense)
	splitExpenses.Get("/", expenseHandler.ListSplitExpenses) // Added list endpoint
	splitExpenses.Get("/:id", expenseHandler.GetSplitExpense)
	splitExpenses.Put("/:id", expenseHandler.UpdateSplitExpense)
	splitExpenses.Delete("/:id", expenseHandler.DeleteSplitExpense)

	// Split Share routes
	splitShares := v1.Group("/split-shares")
	splitShares.Use(middleware.AuthMiddleware())
	splitShares.Post("/", expenseHandler.CreateSplitShare)
	splitShares.Get("/", expenseHandler.ListSplitShares)
	splitShares.Get("/:id", expenseHandler.GetSplitShare)
	splitShares.Put("/:id", expenseHandler.UpdateSplitShare)
	splitShares.Delete("/:id", expenseHandler.DeleteSplitShare)
}
