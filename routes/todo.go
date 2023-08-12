package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/models"
)

type todoRouteHandler struct {
	todoCollection *mongo.Collection
}

func NewTodoRouteHandler(todoCollection *mongo.Collection) *todoRouteHandler {
	return &todoRouteHandler{
		todoCollection: todoCollection,
	}
}

func (t *todoRouteHandler) SetupRoutes(app *fiber.App) {
	app.Post("/todos", t.createTodo)
	app.Get("/todos/:todoId", t.getTodoById)
}

// HANDLERS

func (t *todoRouteHandler) createTodo(c *fiber.Ctx) error {
	reqBody := struct {
		Text string
	}{}
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	newTodo := models.Todo{Text: reqBody.Text}
	result, err := t.todoCollection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (t *todoRouteHandler) getTodoById(c *fiber.Ctx) error {
	var result bson.M

	objId, err := primitive.ObjectIDFromHex(c.Params("todoId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not convert todoId string to ObjectID")
	}

	err = t.todoCollection.FindOne(
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
}
