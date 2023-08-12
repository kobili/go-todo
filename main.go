package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"go-todo/database"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	mongoClient := database.InitDB()

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	log.Fatal(app.Listen(":3001"))
}
