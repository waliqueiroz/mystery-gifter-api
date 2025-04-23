package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserController struct {
	userService       application.UserService
	identityGenerator domain.IdentityGenerator
	passwordManager   domain.PasswordManager
}

func NewUserController(userService application.UserService, identityGenerator domain.IdentityGenerator, passwordManager domain.PasswordManager) *UserController {
	return &UserController{
		userService:       userService,
		identityGenerator: identityGenerator,
		passwordManager:   passwordManager,
	}
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var createUserDTO CreateUserDTO

	if err := ctx.BodyParser(&createUserDTO); err != nil {
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

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID})
}

func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	user, err := c.userService.GetByID(ctx.Context(), userID)
	if err != nil {
		return err
	}

	userDTO, err := mapUserFromDomain(*user)
	if err != nil {
		return err
	}

	return ctx.JSON(userDTO)
}
