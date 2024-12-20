package infra

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/config"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/identity"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/security"
)

func Run() error {
	time.Local = time.UTC

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

	app := fiber.New(fiber.Config{
		ErrorHandler: entrypoint.CustomErrorHandler,
	})
	app.Use(cors.New())
	app.Use(recover.New())

	uuidIdentityGenerator := identity.NewUUIDIdentityGenerator(uuid.NewV7)
	bcryptPasswordManager := security.NewBcryptPasswordManager()

	userRepository := postgres.NewUserRepository(db)
	userService := application.NewUserService(userRepository)
	userController := rest.NewUserController(userService, uuidIdentityGenerator, bcryptPasswordManager)

	entrypoint.CreateRoutes(app, userController)

	return app.Listen(fmt.Sprintf(":%d", 8080))
}
