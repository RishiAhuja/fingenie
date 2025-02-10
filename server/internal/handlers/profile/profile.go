package profile

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
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var user models.User
	result := h.db.Preload("IncomeStreams").
		Preload("BehavioralPatterns").
		Preload("Budgets").
		Where("id = ?", userID).
		First(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var updateData struct {
		DisplayName       string     `json:"displayName"`
		PhoneNumber       string     `json:"phoneNumber"`
		PreferredCurrency string     `json:"preferredCurrency"`
		TelegramID        string     `json:"telegramId"`
		WhatsappNumber    string     `json:"whatsappNumber"`
		NextSalaryDate    *time.Time `json:"nextSalaryDate"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	result := h.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updateData)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not update profile",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
	})
}

func (h *Handler) AddIncomeStream(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var incomeStream models.IncomeStream
	if err := c.BodyParser(&incomeStream); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	incomeStream.UserID = userID

	if result := h.db.Create(&incomeStream); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not create income stream",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    incomeStream,
	})
}

func (h *Handler) UpdateIncomeStream(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	streamID := c.Params("streamId")

	var updateData models.IncomeStream
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	result := h.db.Model(&models.IncomeStream{}).
		Where("id = ? AND user_id = ?", streamID, userID).
		Updates(updateData)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not update income stream",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Income stream updated successfully",
	})
}

func (h *Handler) DeleteIncomeStream(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	streamID := c.Params("streamId")

	result := h.db.Where("id = ? AND user_id = ?", streamID, userID).
		Delete(&models.IncomeStream{})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Could not delete income stream",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Income stream deleted successfully",
	})
}
// Add this function to the existing profile.go file

func (h *Handler) GetUsersByPhoneNumber(c *fiber.Ctx) error {
    phoneNumber := c.Params("phoneNumber")
    
    // Basic validation for phone number
    if phoneNumber == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Phone number is required",
        })
    }

    var users []models.User
    result := h.db.Select("id, display_name, phone_number, telegram_id, whatsapp_number").
        Where("phone_number = ?", phoneNumber).
        Find(&users)

    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Error fetching users",
        })
    }

    if len(users) == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "error":   "No users found with this phone number",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    users,
    })
}