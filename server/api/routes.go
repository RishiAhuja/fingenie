package api

import (
	"github.com/davinder1436/fingenie/internal/handlers/auth"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// API version group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health check route
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is healthy",
		})
	})

	// Initialize handlers
	authHandler := auth.NewHandler(db)

	// Auth routes
	authGroup := v1.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/signup", authHandler.Signup)

	SetupProfileRoutes(app, db)
	SetupGroupRoutes(app, db)
	SetupExpenseRoutes(app, db)
}
