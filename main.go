package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/database"
	"go-todo/routes"
)

func main() {
	mongoClient := database.InitMongoClient()

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := setUpFiberApp(mongoClient)
	log.Fatal(app.Listen(":3001"))
}

func setUpFiberApp(mongoClient *mongo.Client) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	todoCollection := mongoClient.Database("go-todo").Collection("todos")
	todoRouteHandler := routes.NewTodoRouteHandler(todoCollection)
	todoRouteHandler.SetupRoutes(app)

	return app
}
