package expense

import (
	"errors"
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateSplitExpense handles creating a split expense
type CreateSplitExpenseRequest struct {
	GroupID            string         `json:"groupId"`
	ExpenseID          string         `json:"expenseId"`
	TotalAmount        float64        `json:"totalAmount"`
	SplitType          string         `json:"splitType"`
	SettlementPriority int            `json:"settlementPriority"`
	GraceEndDate       time.Time      `json:"graceEndDate"`
	NeedsApproval      bool           `json:"needsApproval"`
	DueDate            time.Time      `json:"dueDate"`
	Shares             []ShareRequest `json:"shares"`
}

type ShareRequest struct {
	UserID string  `json:"userId"`
	Amount float64 `json:"amount"`
}

func (h *Handler) CreateSplitExpense(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	var req CreateSplitExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Set default values if not provided
	if req.SplitType == "" {
		req.SplitType = "EQUAL"
	}
	if req.GraceEndDate.IsZero() {
		req.GraceEndDate = time.Now().Add(24 * time.Hour * 7) // 7 days default
	}

	splitExpense := models.SplitExpense{
		GroupID:            req.GroupID,
		ExpenseID:          req.ExpenseID,
		CreatedBy:          userID,
		TotalAmount:        req.TotalAmount,
		SplitType:          req.SplitType,
		SettlementPriority: req.SettlementPriority,
		GraceEndDate:       req.GraceEndDate,
		NeedsApproval:      req.NeedsApproval,
		DueDate:            req.DueDate,
	}

	if err := h.db.Create(&splitExpense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create split expense",
		})
	}

	// Create shares if provided
	if len(req.Shares) > 0 {
		for _, share := range req.Shares {
			splitShare := models.SplitShare{
				SplitExpenseID: splitExpense.ID,
				UserID:         share.UserID,
				Amount:         share.Amount,
			}

			if err := h.db.Create(&splitShare).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create split shares",
				})
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(splitExpense)
}

// Request structs
type UpdateSplitExpenseRequest struct {
	TotalAmount        float64   `json:"totalAmount"`
	SplitType          string    `json:"splitType"`
	SettlementPriority int       `json:"settlementPriority"`
	GraceEndDate       time.Time `json:"graceEndDate"`
	DueDate            time.Time `json:"dueDate"`
	NeedsApproval      bool      `json:"needsApproval"`
}

type CreateSplitShareRequest struct {
	SplitExpenseID    string  `json:"splitExpenseId"`
	UserID            string  `json:"userId"`
	Amount            float64 `json:"amount"`
	InterestRate      float64 `json:"interestRate"`
	ReminderFrequency string  `json:"reminderFrequency"`
}

type UpdateSplitShareRequest struct {
	Amount            float64 `json:"amount"`
	IsPaid            bool    `json:"isPaid"`
	InterestRate      float64 `json:"interestRate"`
	InterestAccrued   float64 `json:"interestAccrued"`
	ReminderFrequency string  `json:"reminderFrequency"`
}

// UpdateSplitExpense updates an existing split expense
func (h *Handler) UpdateSplitExpense(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitExpenseID := c.Params("id")
	if splitExpenseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split expense ID is required",
		})
	}

	var req UpdateSplitExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	var splitExpense models.SplitExpense
	if err := h.db.First(&splitExpense, "id = ? AND created_by = ?", splitExpenseID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Split expense not found or unauthorized",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve split expense",
		})
	}

	splitExpense.TotalAmount = req.TotalAmount
	splitExpense.SplitType = req.SplitType
	splitExpense.SettlementPriority = req.SettlementPriority
	splitExpense.GraceEndDate = req.GraceEndDate
	splitExpense.DueDate = req.DueDate
	splitExpense.NeedsApproval = req.NeedsApproval

	if err := h.db.Save(&splitExpense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update split expense",
		})
	}

	return c.JSON(splitExpense)
}

// DeleteSplitExpense removes a split expense
func (h *Handler) DeleteSplitExpense(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitExpenseID := c.Params("id")
	if splitExpenseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split expense ID is required",
		})
	}

	result := h.db.Delete(&models.SplitExpense{}, "id = ? AND created_by = ?", splitExpenseID, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete split expense",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Split expense not found or unauthorized",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSplitExpense retrieves a single split expense
func (h *Handler) GetSplitExpense(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitExpenseID := c.Params("id")
	if splitExpenseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split expense ID is required",
		})
	}

	var splitExpense models.SplitExpense
	if err := h.db.Preload("Shares").Preload("Expense").
		First(&splitExpense, "id = ?", splitExpenseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Split expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve split expense",
		})
	}

	// Check if user is either the creator or part of the shares
	isAuthorized := splitExpense.CreatedBy == userID
	if !isAuthorized {
		for _, share := range splitExpense.Shares {
			if share.UserID == userID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Not authorized to view this split expense",
		})
	}

	return c.JSON(splitExpense)
}

// ListSplitExpenses retrieves all split expenses for a user
func (h *Handler) ListSplitExpenses(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	var splitExpenses []models.SplitExpense
	if err := h.db.Preload("Shares").Preload("Expense").
		Joins("LEFT JOIN split_shares ON split_shares.split_expense_id = split_expenses.id").
		Where("split_expenses.created_by = ? OR split_shares.user_id = ?", userID, userID).
		Group("split_expenses.id").
		Find(&splitExpenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve split expenses",
		})
	}

	return c.JSON(splitExpenses)
}

// CreateSplitShare handles creating a split share
func (h *Handler) CreateSplitShare(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	var req CreateSplitShareRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Verify split expense exists and user is authorized
	var splitExpense models.SplitExpense
	if err := h.db.First(&splitExpense, "id = ? AND created_by = ?", req.SplitExpenseID, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Split expense not found or unauthorized",
		})
	}

	splitShare := models.SplitShare{
		SplitExpenseID:    req.SplitExpenseID,
		UserID:            req.UserID,
		Amount:            req.Amount,
		InterestRate:      req.InterestRate,
		ReminderFrequency: req.ReminderFrequency,
	}

	if err := h.db.Create(&splitShare).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create split share",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(splitShare)
}

// UpdateSplitShare updates an existing split share
func (h *Handler) UpdateSplitShare(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitShareID := c.Params("id")
	if splitShareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split share ID is required",
		})
	}

	var req UpdateSplitShareRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	var splitShare models.SplitShare
	if err := h.db.Preload("SplitExpense").First(&splitShare, "id = ?", splitShareID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Split share not found",
		})
	}

	// Check if user is either the expense creator or the share owner
	if splitShare.SplitExpense.CreatedBy != userID && splitShare.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Not authorized to update this split share",
		})
	}

	splitShare.Amount = req.Amount
	splitShare.IsPaid = req.IsPaid
	splitShare.InterestRate = req.InterestRate
	splitShare.InterestAccrued = req.InterestAccrued
	splitShare.ReminderFrequency = req.ReminderFrequency

	if err := h.db.Save(&splitShare).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update split share",
		})
	}

	return c.JSON(splitShare)
}

// DeleteSplitShare removes a split share
func (h *Handler) DeleteSplitShare(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitShareID := c.Params("id")
	if splitShareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split share ID is required",
		})
	}

	var splitShare models.SplitShare
	if err := h.db.Preload("SplitExpense").First(&splitShare, "id = ?", splitShareID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Split share not found",
		})
	}

	// Only expense creator can delete shares
	if splitShare.SplitExpense.CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Not authorized to delete this split share",
		})
	}

	if err := h.db.Delete(&splitShare).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete split share",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSplitShare retrieves a single split share
func (h *Handler) GetSplitShare(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	splitShareID := c.Params("id")
	if splitShareID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Split share ID is required",
		})
	}

	var splitShare models.SplitShare
	if err := h.db.Preload("SplitExpense").First(&splitShare, "id = ?", splitShareID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Split share not found",
		})
	}

	// Check if user is either the expense creator or the share owner
	if splitShare.SplitExpense.CreatedBy != userID && splitShare.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Not authorized to view this split share",
		})
	}

	return c.JSON(splitShare)
}

// ListSplitShares retrieves all split shares for a user
func (h *Handler) ListSplitShares(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	var splitShares []models.SplitShare
	if err := h.db.Preload("SplitExpense").
		Joins("JOIN split_expenses ON split_expenses.id = split_shares.split_expense_id").
		Where("split_shares.user_id = ? OR split_expenses.created_by = ?", userID, userID).
		Find(&splitShares).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve split shares",
		})
	}

	return c.JSON(splitShares)
}
