package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"go-todo/database"
	"go-todo/models"
)

type todoRouteHandler struct {
	todoRepo database.TodoRepo
}

func NewTodoRouteHandler(todoRepository database.TodoRepo) *todoRouteHandler {
	return &todoRouteHandler{
		todoRepo: todoRepository,
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

	result, err := t.todoRepo.AddOne(c.Context(), models.Todo{Text: reqBody.Text})
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (t *todoRouteHandler) getTodoById(c *fiber.Ctx) error {
	result, err := t.todoRepo.FindById(c.Context(), c.Params("todoId"))

	if err == mongo.ErrNoDocuments {
		return fiber.NewError(fiber.StatusNotFound, "No todo was found with the given ID")
	}

	if err != nil {
		return err
	}

	return c.JSON(result)
}
