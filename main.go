package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectDB()
	app := fiber.New()

	// Register API routes
	registerRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
