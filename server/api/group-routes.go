package api

import (
	"github.com/davinder1436/fingenie/internal/handlers/group"
	"github.com/davinder1436/fingenie/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupGroupRoutes(app *fiber.App, db *gorm.DB) {
	h := group.NewHandler(db)
	mh := group.NewGroupMemberHandler(db) // Group Member Handler

	groupRoutes := app.Group("/api/v1/groups")
	groupRoutes.Use(middleware.AuthMiddleware())
	groupRoutes.Post("/", h.CreateGroup)        // Create a new group
	groupRoutes.Put("/:groupId", h.UpdateGroup) // Update an existing group
	groupRoutes.Get("/:groupId", h.GetGroup)    // Fetch group details
	groupRoutes.Delete("/:groupId", h.DeleteGroup)
	groupRoutes.Get("/user-groups", h.ListUserGroups)

	// Group Member Routes
	memberRoutes := groupRoutes.Group("/:groupId/members")
	memberRoutes.Use(middleware.AuthMiddleware())
	memberRoutes.Post("/", mh.AddMember)               // Add a member
	memberRoutes.Get("/", mh.GetMembers)               // List group members
	memberRoutes.Patch("/:memberId", mh.UpdateMember)  // Update member role/share
	memberRoutes.Delete("/:memberId", mh.RemoveMember) // Remove a member
}
