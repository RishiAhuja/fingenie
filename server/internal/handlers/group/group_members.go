package group

import (
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GroupMemberHandler struct {
	db *gorm.DB
}

func NewGroupMemberHandler(db *gorm.DB) *GroupMemberHandler {
	return &GroupMemberHandler{db: db}
}

// Add a new member to a group
func (h *GroupMemberHandler) AddMember(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")

	var req struct {
		UserID       string  `json:"userId"`
		Role         string  `json:"role"`
		SharePercent float64 `json:"sharePercent"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}

	var admin models.GroupMember
	h.db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&admin)
	if admin.UserID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Only admin can add members"})
	}

	member := models.GroupMember{
		GroupID:      groupID,
		UserID:       req.UserID,
		Role:         req.Role,
		JoinedAt:     time.Now(),
		IsActive:     true,
		SharePercent: req.SharePercent,
	}

	if err := h.db.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Could not add member"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Member added successfully", "data": member})
}

// Get all members of a group
func (h *GroupMemberHandler) GetMembers(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	var members []models.GroupMember

	if err := h.db.Where("group_id = ?", groupID).Preload("User").Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Could not retrieve members"})
	}

	return c.JSON(fiber.Map{"success": true, "data": members})
}

// Update a member's role or share percent
func (h *GroupMemberHandler) UpdateMember(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")
	memberID := c.Params("memberId")

	var req struct {
		Role         string  `json:"role"`
		SharePercent float64 `json:"sharePercent"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Invalid request body"})
	}

	var admin models.GroupMember
	h.db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&admin)
	if admin.UserID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Only admin can update members"})
	}

	if err := h.db.Model(&models.GroupMember{}).Where("id = ?", memberID).Updates(req).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Could not update member"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Member updated successfully"})
}

// Remove a member from a group
func (h *GroupMemberHandler) RemoveMember(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")
	memberID := c.Params("memberId")

	var admin models.GroupMember
	h.db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&admin)
	if admin.UserID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Only admin can remove members"})
	}

	if err := h.db.Where("id = ?", memberID).Delete(&models.GroupMember{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Could not remove member"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Member removed successfully"})
}

// Register routes for group member management
func RegisterGroupMemberRoutes(app *fiber.App, db *gorm.DB) {
	handler := NewGroupMemberHandler(db)
	groupRoutes := app.Group("/groups/:groupId/members")

	groupRoutes.Post("/", handler.AddMember)
	groupRoutes.Get("/", handler.GetMembers)
	groupRoutes.Patch("/:memberId", handler.UpdateMember)
	groupRoutes.Delete("/:memberId", handler.RemoveMember)
}
