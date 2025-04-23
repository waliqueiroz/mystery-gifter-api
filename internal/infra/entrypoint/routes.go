package entrypoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

func CreateRoutes(router fiber.Router, authMiddlware fiber.Handler, userController *rest.UserController, authController *rest.AuthController, groupController *rest.GroupController) {
	api := router.Group("/api")

	api.Post("/login", authController.Login)
	api.Post("/users", userController.Create)

	api.Use(authMiddlware) // from now on, all routes will require authentication

	api.Get("/users/:userID", userController.GetByID)
	api.Post("/groups", groupController.Create)
}
