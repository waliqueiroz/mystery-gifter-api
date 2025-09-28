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

	userDTO, err := mapUserFromDomain(*user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(userDTO)
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

func (c *UserController) Search(ctx *fiber.Ctx) error {
	var userFiltersDTO UserFiltersDTO

	if err := ctx.QueryParser(&userFiltersDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	filters, err := mapUserFiltersDTOToDomain(userFiltersDTO)
	if err != nil {
		return err
	}

	searchResult, err := c.userService.Search(ctx.Context(), *filters)
	if err != nil {
		return err
	}

	searchResultDTO, err := mapSearchResultFromDomain(searchResult, mapUserFromDomain)
	if err != nil {
		return err
	}

	return ctx.JSON(searchResultDTO)
}
