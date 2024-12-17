package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
)

type UserController struct {
	userService application.UserService
}

func NewUserController(userService application.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Create(ctx fiber.Ctx) error {
	var createUserDTO CreateUserDTO

	if err := ctx.Bind().Body(&createUserDTO); err != nil {
		return err
	}

	user, err := mapCreateUserDTOToUser(createUserDTO)
	if err != nil {
		return err
	}

	err = c.userService.Create(ctx.Context(), *user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": user.ID})
}

func (c *UserController) GetByID(ctx fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "you should pass a valid user ID")
	}

	user, err := c.userService.GetByID(ctx.Context(), userID)
	if err != nil {
		return err
	}

	userDTO, err := mapUserToUserDTO(*user)
	if err != nil {
		return err
	}

	return ctx.JSON(userDTO)
}
