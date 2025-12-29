package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/dto"
	"github.com/arieffadhlan/go-fitbyte/internal/usecases/profile"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	profileUseCase profile.ProfileUseCaseInterface
}

func NewProfileHandler(profileUseCase profile.ProfileUseCaseInterface) *ProfileHandler {
	return &ProfileHandler{
		profileUseCase: profileUseCase,
	}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	// TODO: Replace with JWT middleware authentication
	userID := c.Locals("id").(int)

	// // Get user ID from query parameter for now
	// userIDStr := c.Query("user_id")
	// if userIDStr == "" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "user_id query parameter is required",
	// 	})
	// }

	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid user_id format",
	// 	})
	// }

	profile, err := h.profileUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	// TODO: Replace with JWT middleware authentication
	userID := c.Locals("id").(int)

	// Get user ID from query parameter for now
	// userIDStr := c.Query("user_id")
	// if userIDStr == "" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "user_id query parameter is required",
	// 	})
	// }

	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid user_id format",
	// 	})
	// }

	var req dto.ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	profile, err := h.profileUseCase.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		// Handle validation errors
		if err.Error() == "validation error" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}
