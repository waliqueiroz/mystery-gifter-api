package infra

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/config"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		return err
	}

	defer db.Close()

	err = postgres.Migrate(db.DB)
	if err != nil {
		return err
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())

	entrypoint.CreateRoutes(app)

	return app.Listen(fmt.Sprintf(":%d", 8080))
}
