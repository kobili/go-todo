package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/database"
	"go-todo/models"
)

type TodoRouteHandler struct {
	todoRepo database.TodoRepo
}

func NewTodoRouteHandler(todoRepository database.TodoRepo) *TodoRouteHandler {
	return &TodoRouteHandler{
		todoRepo: todoRepository,
	}
}

func (t *TodoRouteHandler) SetupRoutes(app *fiber.App) {
	app.Post("/todos", t.createTodo)
	app.Get("/todos/:todoId", t.getTodoById)
	app.Patch("/todos/:todoId", t.updateTodoById)
}

// HANDLERS

func (t *TodoRouteHandler) createTodo(c *fiber.Ctx) error {
	reqBody := struct {
		Text string
	}{}
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	result, err := t.todoRepo.AddOne(c.Context(), models.Todo{
		Text: reqBody.Text,
		Metadata: models.TodoMetadata{
			Timestamp: time.Now(),
			User:      "Anon",
		},
	})
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (t *TodoRouteHandler) getTodoById(c *fiber.Ctx) error {
	result, err := t.todoRepo.FindById(c.Context(), c.Params("todoId"))

	if err == mongo.ErrNoDocuments {
		return fiber.NewError(fiber.StatusNotFound, "No todo was found with the given ID")
	}

	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (t *TodoRouteHandler) updateTodoById(c *fiber.Ctx) error {
	requestBody := struct {
		Text string
	}{}
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}
	id := c.Params("todoId")

	result, err := t.todoRepo.UpdateById(c.Context(), id, requestBody)
	if err == mongo.ErrNoDocuments {
		return fiber.NewError(fiber.StatusNotFound, "No todo was found with the given ID")
	}
	if err != nil {
		return err
	}

	return c.JSON(result)
}
