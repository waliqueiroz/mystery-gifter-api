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

	authUserID, err := c.tokenManager.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return err
	}

	group, err := c.groupService.Create(ctx.Context(), createGroupDTO.Name, authUserID)
	if err != nil {
		return err
	}

	groupDTO, err := mapGroupFromDomain(*group)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(groupDTO)
}

func (c *GroupController) GetByID(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupID")

	group, err := c.groupService.GetByID(ctx.Context(), groupID)
	if err != nil {
		return err
	}

	groupDTO, err := mapGroupFromDomain(*group)
	if err != nil {
		return err
	}

	return ctx.JSON(groupDTO)
}

func (c *GroupController) AddUser(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupID")

	var addUserDTO AddUserDTO

	if err := ctx.BodyParser(&addUserDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	if err := addUserDTO.Validate(); err != nil {
		return err
	}

	authUserID, err := c.tokenManager.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return err
	}

	group, err := c.groupService.AddUser(ctx.Context(), groupID, authUserID, addUserDTO.UserID)
	if err != nil {
		return err
	}

	groupDTO, err := mapGroupFromDomain(*group)
	if err != nil {
		return err
	}

	return ctx.JSON(groupDTO)
}

func (c *GroupController) RemoveUser(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupID")
	targetUserID := ctx.Params("userID")

	authUserID, err := c.tokenManager.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return err
	}

	group, err := c.groupService.RemoveUser(ctx.Context(), groupID, authUserID, targetUserID)
	if err != nil {
		return err
	}

	groupDTO, err := mapGroupFromDomain(*group)
	if err != nil {
		return err
	}

	return ctx.JSON(groupDTO)
}
