package entrypoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

func CreateRoutes(router fiber.Router, authMiddleware fiber.Handler, userController *rest.UserController, authController *rest.AuthController, groupController *rest.GroupController) {
	api := router.Group("/api")

	api.Post("/login", authController.Login)
	api.Post("/users", userController.Create)

	api.Use(authMiddleware) // from now on, all routes will require authentication

	api.Get("/users", userController.Search)
	api.Get("/users/:userID", userController.GetByID)
	api.Post("/groups", groupController.Create)
	api.Get("/groups/:groupID", groupController.GetByID)
	api.Post("/groups/:groupID/users", groupController.AddUser)
	api.Delete("/groups/:groupID/users/:userID", groupController.RemoveUser)
	api.Post("/groups/:groupID/matches", groupController.GenerateMatches)
	api.Post("/groups/:groupID/reopen", groupController.Reopen)
	api.Post("/groups/:groupID/archive", groupController.Archive)
	api.Get("/groups/:groupID/matches/user", groupController.GetUserMatch)
}
