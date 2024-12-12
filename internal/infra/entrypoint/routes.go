package entrypoint

import "github.com/gofiber/fiber/v3"

func CreateRoutes(router fiber.Router) {
	router.Group("/api")
}
