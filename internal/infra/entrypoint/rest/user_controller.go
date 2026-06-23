package rest

import (
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserController struct {
	userService       application.UserService
	identityGenerator domain.IdentityGenerator
	passwordManager   domain.PasswordManager
	authTokenManager  domain.AuthTokenManager
}

func NewUserController(userService application.UserService, identityGenerator domain.IdentityGenerator, passwordManager domain.PasswordManager, authTokenManager domain.AuthTokenManager) *UserController {
	return &UserController{
		userService:       userService,
		identityGenerator: identityGenerator,
		passwordManager:   passwordManager,
		authTokenManager:  authTokenManager,
	}
}

func (c *UserController) GetMe(ctx fiber.Ctx) error {
	authUserID, err := c.authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))
	if err != nil {
		return err
	}

	user, err := c.userService.GetByID(ctx.Context(), authUserID)
	if err != nil {
		return err
	}

	userDTO, err := mapUserFromDomain(*user)
	if err != nil {
		return err
	}

	return ctx.JSON(userDTO)
}

func (c *UserController) Create(ctx fiber.Ctx) error {
	var createUserDTO CreateUserDTO

	if err := ctx.Bind().Body(&createUserDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	user, err := mapCreateUserDTOToDomain(c.identityGenerator, c.passwordManager, createUserDTO)
	if err != nil {
		return err
	}

	err = c.userService.Create(ctx.Context(), *user)
	if err != nil {
		return err
	}

	userDTO, err := mapUserFromDomain(*user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(userDTO)
}

