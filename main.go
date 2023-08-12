package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/database"
)

type Todo struct {
	Text string
}

func main() {
	app := fiber.New()

	mongoClient := database.InitMongoClient()

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := mongoClient.Database("go-todo").Collection("todos")

	app.Post("/todos", func(c *fiber.Ctx) error {
		reqBody := struct {
			Text string
		}{}
		if err := c.BodyParser(&reqBody); err != nil {
			return err
		}

		newTodo := Todo{Text: reqBody.Text}
		result, err := coll.InsertOne(context.TODO(), newTodo)
		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	app.Get("/todos/:todoId", func(c *fiber.Ctx) error {
		var result bson.M

		objId, err := primitive.ObjectIDFromHex(c.Params("todoId"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Could not convert todoId string to ObjectID")
		}

		err = coll.FindOne(
			context.Background(),
			bson.D{{Key: "_id", Value: objId}},
		).Decode(&result)

		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "No todo was found with the given ID")
		}

		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	log.Fatal(app.Listen(":3001"))
}
