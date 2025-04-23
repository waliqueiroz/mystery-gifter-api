package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupController struct {
	groupService application.GroupService
	tokenManager domain.TokenManager
}

func NewGroupController(
	groupService application.GroupService,
	tokenManager domain.TokenManager,
) *GroupController {
	return &GroupController{
		groupService: groupService,
		tokenManager: tokenManager,
	}
}

func (c *GroupController) Create(ctx *fiber.Ctx) error {
	var createGroupDTO CreateGroupDTO

	if err := ctx.BodyParser(&createGroupDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	if err := createGroupDTO.Validate(); err != nil {
		return err
	}

	userID, err := c.tokenManager.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return err
	}

	group, err := c.groupService.Create(ctx.Context(), createGroupDTO.Name, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": group.ID})
}
