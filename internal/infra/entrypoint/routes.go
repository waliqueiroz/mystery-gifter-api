package entrypoint

import (
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

func CreateRoutes(router fiber.Router, userController *rest.UserController) {
	api := router.Group("/api")

	api.Post("/users", userController.Create)
	api.Get("/users/:userID", userController.GetByID)
}
