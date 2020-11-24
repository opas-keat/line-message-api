package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(recover.New())

	// Routes
	app.Get("/", hello)

	app.Post("/webhook", webhook)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler webhook
func webhook(c *fiber.Ctx) error {
	return c.SendString("Hi👋!, This Line Message Api")
}

// Handler hello
func hello(c *fiber.Ctx) error {
	return c.SendString("Hi👋!, This Line Message Api")
}
