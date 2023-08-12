package routes

import (
	"github.com/gofiber/fiber/v2"
)

type RouteHandler interface {
	SetupRoutes(app *fiber.App)
}
