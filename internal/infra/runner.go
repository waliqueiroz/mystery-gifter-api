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
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/identity"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/security"
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

	err = postgres.Migrate(db.GetDB())
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
	jwtAuthTokenManager := security.NewJWTAuthTokenManager(cfg.Auth.SecretKey)

	userRepository := postgres.NewUserRepository(db)
	userService := application.NewUserService(userRepository)
	userController := rest.NewUserController(userService, uuidIdentityGenerator, bcryptPasswordManager)

	groupRepository := postgres.NewGroupRepository(db)
	groupService := application.NewGroupService(groupRepository, userService, uuidIdentityGenerator)
	groupController := rest.NewGroupController(groupService, jwtAuthTokenManager)

	groupInviteRepository := postgres.NewGroupInviteRepository(db)
	groupInviteService := application.NewGroupInviteService(groupInviteRepository, groupRepository, userService, uuidIdentityGenerator, cfg.Invite.LinkExpiration)
	groupInviteController := rest.NewGroupInviteController(groupInviteService, jwtAuthTokenManager)

	authService := application.NewAuthService(cfg.Auth.SessionDuration, userRepository, bcryptPasswordManager, jwtAuthTokenManager)
	authController := rest.NewAuthController(authService, cfg.Auth.CookieSecure)

	authMiddleware := entrypoint.NewAuthMiddleware(cfg.Auth.SecretKey)

	entrypoint.CreateRoutes(app, authMiddleware, userController, authController, groupController, groupInviteController)

	return app.Listen(fmt.Sprintf(":%d", 8080))
}
