package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/database"
)

func main() {
	app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	mongoClient := database.InitMongoClient()

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := mongoClient.Database("go-todo").Collection("todos")

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"text", "hello from mongodb"}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found")
	}
	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(result)
	})

	log.Fatal(app.Listen(":3001"))
}
