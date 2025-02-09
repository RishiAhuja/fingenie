package profile

import (
	"errors"
	"time"

	"github.com/davinder1436/fingenie/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (h *Handler) GetUserSocialScoreHistory(c *fiber.Ctx) error {
	// Get userId from context (assuming it was set by auth middleware)
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var socialScores []models.SocialScoreHistory
	result := h.db.Where("user_id = ?", userId).
		Order("timestamp DESC").
		Find(&socialScores)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error fetching social score history",
		})
	}

	// If no scores found, return empty array with success true
	if len(socialScores) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []models.SocialScoreHistory{},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    socialScores,
	})
}

func (h *Handler) SetupUserSocialScore(c *fiber.Ctx) error {
	// Get userId from context
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	// Check if user already has a social score
	var existingScore models.SocialScoreHistory
	result := h.db.Where("user_id = ?", userId).Order("timestamp DESC").First(&existingScore)

	// If no score exists, create initial score
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			initialScore := models.SocialScoreHistory{
				UserID:    userId.(string),
				OldScore:  0,
				NewScore:  50, // Default starting score
				Reason:    "Initial score setup",
				Timestamp: time.Now(),
			}

			if err := h.db.Create(&initialScore).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error":   "Failed to create initial social score",
				})
			}

			return c.JSON(fiber.Map{
				"success": true,
				"data":    initialScore,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error checking existing score",
		})
	}

	// If score already exists, return it
	return c.JSON(fiber.Map{
		"success": true,
		"data":    existingScore,
	})
}

// Request body struct
type UpdateSocialScoreRequest struct {
	NewScore float64 `json:"newScore"`
	Reason   string  `json:"reason"`
}

func (h *Handler) UpdateUserSocialScore(c *fiber.Ctx) error {
	// Get userId from context
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	// Parse request body
	var req UpdateSocialScoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate request
	if req.NewScore < 0 || req.NewScore > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Score must be between 0 and 100",
		})
	}

	if req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Reason is required",
		})
	}

	// Get current score
	var currentScore models.SocialScoreHistory
	result := h.db.Where("user_id = ?", userId).
		Order("timestamp DESC").
		First(&currentScore)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "No existing social score found. Please setup initial score first",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error fetching current score",
		})
	}

	// Create new score history entry
	newScoreEntry := models.SocialScoreHistory{
		UserID:    userId.(string),
		OldScore:  currentScore.NewScore, // Previous score becomes old score
		NewScore:  req.NewScore,
		Reason:    req.Reason,
		Timestamp: time.Now(),
	}

	if err := h.db.Create(&newScoreEntry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update social score",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    newScoreEntry,
	})
}
