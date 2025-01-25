package entrypoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

func CreateRoutes(router fiber.Router, userController *rest.UserController, authController *rest.AuthController) {
	api := router.Group("/api")

	api.Post("/login", authController.Login)

	api.Post("/users", userController.Create)
	api.Get("/users/:userID", userController.GetByID)
}
