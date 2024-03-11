package main

import (
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App) {
	app.Post("/api/register", registerHandler)
}

func registerHandler(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}

	// Parse request body
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"emessagerror": "Invalid request body",
		})
	}

	// Check if teacher email is provided
	if req.Teacher == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Teacher email is required",
		})
	}

	// Check if students list is provided
	if len(req.Students) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one student email is required",
		})
	}

	// Register teacher and students in the database
	if err := registerTeacherStudent(req.Teacher, req.Students); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register teacher and students",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
