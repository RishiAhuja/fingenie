package expense

import (
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateExpenseRequest represents the structure of the expense creation request
type CreateExpenseRequest struct {
	Amount           float64 `json:"amount"`
	Category         string  `json:"category"`
	GroupID          string  `json:"groupId"`
	Description      string  `json:"description"`
	OriginalCurrency string  `json:"originalCurrency"`
	SplitType        string  `json:"splitType"` // EQUAL, PERCENTAGE, CUSTOM
}

// CreateExpense handles the creation of an expense with split expenses
func (h *Handler) CreateExpense(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - valid user ID required",
		})
	}

	var req CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	expense := models.Expense{
		UserID:           userID,
		GroupID:          &req.GroupID,
		Amount:           req.Amount,
		OriginalCurrency: req.OriginalCurrency,
		Category:         req.Category,
		Description:      req.Description,
		Date:             time.Now(),
	}

	if err := h.db.Create(&expense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create expense",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(expense)
}

// UpdateExpense updates an existing expense
func (h *Handler) UpdateExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")
	var req models.Expense

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var expense models.Expense
	if err := h.db.First(&expense, "id = ?", expenseID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expense not found"})
	}

	expense.Amount = req.Amount
	expense.Category = req.Category
	expense.Description = req.Description
	expense.Date = req.Date

	if err := h.db.Save(&expense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update expense"})
	}

	return c.JSON(expense)
}

// DeleteExpense removes an expense from the database
func (h *Handler) DeleteExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")
	if err := h.db.Delete(&models.Expense{}, "id = ?", expenseID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete expense"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetExpense retrieves a single expense by ID
func (h *Handler) GetExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")
	var expense models.Expense
	if err := h.db.First(&expense, "id = ?", expenseID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expense not found"})
	}
	return c.JSON(expense)
}

// ListExpenses retrieves all expenses
func (h *Handler) ListExpenses(c *fiber.Ctx) error {
	var expenses []models.Expense
	if err := h.db.Find(&expenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve expenses"})
	}
	return c.JSON(expenses)
}
