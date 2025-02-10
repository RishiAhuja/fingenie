package group

import (
	"fmt"
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type groupRequest struct {
	Name                    string  `json:"name"`
	Description             string  `json:"description"`
	DefaultCurrency         string  `json:"defaultCurrency"`
	GroupType               string  `json:"groupType"`
	IsRecurring             bool    `json:"isRecurring"`
	SecurityDepositRequired float64 `json:"securityDepositRequired"`
	RequiresAdminApproval   bool    `json:"requiresAdminApproval"`

	BudgetStrategy string `json:"budgetStrategy"`

	BillingCycleStart time.Time `json:"billingCycleStart"`
	SplitStrategy     string    `json:"splitStrategy"`
	AutoSettlement    bool      `json:"autoSettlement"`
}

func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	// Extract the userID from the request context, ensure it's a valid string
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authorized or userId not found",
		})
	}

	var req groupRequest
	// Parse the request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.GroupType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Name and group type are required",
		})
	}

	// Construct the Group object
	group := models.Group{
		Name:                    req.Name,
		CreatedBy:               userID,
		Description:             req.Description,
		DefaultCurrency:         req.DefaultCurrency,
		GroupType:               req.GroupType,
		IsRecurring:             req.IsRecurring,
		SecurityDepositRequired: req.SecurityDepositRequired,
		RequiresAdminApproval:   req.RequiresAdminApproval,
		BudgetStrategy:          req.BudgetStrategy,
		BillingCycleStart:       req.BillingCycleStart,
		SplitStrategy:           req.SplitStrategy,
		AutoSettlement:          req.AutoSettlement,
	}

	// Begin a database transaction to create the group and its first member
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Create the group record in the database
		if err := tx.Create(&group).Error; err != nil {
			return fmt.Errorf("could not create group: %v", err)
		}

		// Create the group member (admin) record
		member := models.GroupMember{
			GroupID:  group.ID,
			UserID:   userID,
			Role:     "admin",
			JoinedAt: time.Now(),
			IsActive: true,
		}

		// Add the member to the group
		if err := tx.Create(&member).Error; err != nil {
			return fmt.Errorf("could not add member to group: %v", err)
		}

		return nil
	})

	// Handle any errors that occurred during the transaction
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Return the created group in the response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    group,
	})
}

func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")

	var req groupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	var member models.GroupMember
	err := h.db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&member).Error

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only admin can update group",
		})
	}

	updates := models.Group{
		Name:                    req.Name,
		Description:             req.Description,
		DefaultCurrency:         req.DefaultCurrency,
		GroupType:               req.GroupType,
		IsRecurring:             req.IsRecurring,
		SecurityDepositRequired: req.SecurityDepositRequired,
		RequiresAdminApproval:   req.RequiresAdminApproval,

		BudgetStrategy: req.BudgetStrategy,

		BillingCycleStart: req.BillingCycleStart,
		SplitStrategy:     req.SplitStrategy,
		AutoSettlement:    req.AutoSettlement,
	}

	if result := h.db.Model(&models.Group{}).Where("id = ?", groupID).Updates(updates); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not update group",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Group updated successfully",
	})
}

func (h *Handler) GetGroup(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")

	var group models.Group
	err := h.db.Preload("Members").
		Preload("Members.User").
		Preload("Expenses").
		Preload("RecurringExpenses").
		Preload("Budgets").
		Where("id = ?", groupID).
		First(&group).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Group not found",
		})
	}

	isMember := false
	for _, member := range group.Members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Not authorized to view this group",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    group,
	})
}

func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	groupID := c.Params("groupId")

	var member models.GroupMember
	err := h.db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&member).Error

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only admin can delete group",
		})
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&models.GroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Where("group_id = ?", groupID).Delete(&models.Expense{}).Error; err != nil {
			return err
		}

		if err := tx.Where("group_id = ?", groupID).Delete(&models.RecurringExpense{}).Error; err != nil {
			return err
		}

		if err := tx.Where("group_id = ?", groupID).Delete(&models.Budget{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Group{}, groupID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not delete group",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Group deleted successfully",
	})
}

func (h *Handler) ListUserGroups(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var groups []models.Group
	err := h.db.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Preload("Members").
		Preload("Members.User").
		Find(&groups).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not fetch groups",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    groups,
	})
}
