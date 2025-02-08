package profile

import (
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
)

type budgetRequest struct {
	Category          string      `json:"category"`
	Tags              interface{} `json:"tags"`
	Amount            float64     `json:"amount"`
	Period            string      `json:"period"`
	StartDate         time.Time   `json:"startDate"`
	EndDate           time.Time   `json:"endDate"`
	AISuggestedAmount float64     `json:"aiSuggestedAmount"`
	IsAutoAdjusting   bool        `json:"isAutoAdjusting"`
	GroupID           *string     `json:"groupId"`
}

func (h *Handler) CreateBudget(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var req budgetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Category == "" || req.Amount == 0 || req.Period == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Category, amount, and period are required",
		})
	}

	var tags []string
	switch v := req.Tags.(type) {
	case string:
		tags = []string{v}
	case []interface{}:
		tags = make([]string, len(v))
		for i, tag := range v {
			if str, ok := tag.(string); ok {
				tags[i] = str
			}
		}
	case []string:
		tags = v
	}

	budget := models.Budget{
		UserID:            userID,
		GroupID:           req.GroupID,
		Category:          req.Category,
		Tags:              tags,
		Amount:            req.Amount,
		Period:            req.Period,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		AISuggestedAmount: req.AISuggestedAmount,
		IsAutoAdjusting:   req.IsAutoAdjusting,
		CurrentSpent:      0,
	}

	if result := h.db.Create(&budget); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not create budget",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    budget,
	})
}
