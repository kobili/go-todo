package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"go-todo/database"
	"go-todo/routes"
)

func main() {
	fx.New(
		fx.Provide(
			database.NewMongoClient,
			fx.Annotate(
				database.NewTodoRepository,
				fx.As(new(database.TodoRepo)),
			),
			routes.NewTodoRouteHandler,
			NewFiberApp,
		),
		fx.Invoke(
			func(*mongo.Client) {},
			func(*fiber.App) {},
		),
	).Run()
}

func NewFiberApp(lc fx.Lifecycle, todoRouteHandler *routes.TodoRouteHandler) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	todoRouteHandler.SetupRoutes(app)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(":3001")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}
