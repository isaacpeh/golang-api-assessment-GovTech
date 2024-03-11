package main

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App) {
	app.Post("/api/register", registerHandler)
	app.Get("/api/commonstudents", retrieveStudents)
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

func retrieveStudents(c *fiber.Ctx) error {
	// Parse query parameters
	originalURL := c.OriginalURL()

	// Parse the URL
	url, err := url.Parse(originalURL)
	if err != nil {
		return err // handle error
	}

	// Extract the query parameters
	query := url.Query()

	// Get the values for the "teacher" parameter
	teacherList := query["teacher"]

	// Check if the teacher list is empty
	if len(teacherList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one teacher email is required",
		})
	}

	// Retrieve common students from the database
	commonStudents, err := getCommonStudents(teacherList)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve common students",
		})
	}

	// Return message if there are no common students for listed teachers
	if len(commonStudents) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No common students found for the provided list of teachers.",
		})
	}

	// Return the list of common students
	return c.JSON(fiber.Map{
		"students": commonStudents,
	})
}
