package main

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App) {
	app.Post("/api/register", registerHandler)
	app.Get("/api/commonstudents", retrieveStudents)
	app.Post("/api/suspend", suspendStudent)
	app.Post("/api/retrievefornotifications", notifyStudents)
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
			"message": "Invalid request body",
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

func suspendStudent(c *fiber.Ctx) error {
	type SuspendRequest struct {
		Student string `json:"student"`
	}

	// Parse request body
	var req SuspendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if teacher email is provided
	if req.Student == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Student email is required",
		})
	}

	// Register teacher and students in the database
	result, err := suspendSpecificStudent(req.Student)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to suepend student",
		})
	}

	if result <= 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Student not found in database",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func notifyStudents(c *fiber.Ctx) error {
	type NotifyRequest struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}

	// Parse request body
	var req NotifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if teacher email is provided
	if req.Teacher == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Teacher email is required",
		})
	}

	// Check if Notification is provided
	if req.Notification == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Notification message is required",
		})
	}

	// Regular expression to match email addresses after the "@"
	re := regexp.MustCompile(`@([\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,})`)
	matches := re.FindAllStringSubmatch(req.Notification, -1)

	// Extract email addresses from the matches
	var emails []string
	for _, match := range matches {
		emails = append(emails, match[1])
	}

	recipients, err := returnRecipients(req.Teacher, emails)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to notify students",
		})
	}

	if len(recipients) <= 0 {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "No recipients of this notification",
		})
	}

	// Return the list of recipients
	return c.JSON(fiber.Map{
		"recipients": recipients,
	})
}
