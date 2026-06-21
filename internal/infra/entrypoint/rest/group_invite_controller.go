package rest

import (
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupInviteController struct {
	groupInviteService application.GroupInviteService
	authTokenManager   domain.AuthTokenManager
}

func NewGroupInviteController(
	groupInviteService application.GroupInviteService,
	authTokenManager domain.AuthTokenManager,
) *GroupInviteController {
	return &GroupInviteController{
		groupInviteService: groupInviteService,
		authTokenManager:   authTokenManager,
	}
}

func (c *GroupInviteController) Create(ctx fiber.Ctx) error {
	groupID := ctx.Params("groupID")

	authUserID, err := c.authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))
	if err != nil {
		return err
	}

	groupInvite, err := c.groupInviteService.Create(ctx.Context(), groupID, authUserID)
	if err != nil {
		return err
	}

	groupInviteDTO, err := mapGroupInviteFromDomain(*groupInvite)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(groupInviteDTO)
}

func (c *GroupInviteController) GetActive(ctx fiber.Ctx) error {
	groupID := ctx.Params("groupID")

	authUserID, err := c.authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))
	if err != nil {
		return err
	}

	groupInvite, err := c.groupInviteService.GetActive(ctx.Context(), groupID, authUserID)
	if err != nil {
		return err
	}

	groupInviteDTO, err := mapGroupInviteFromDomain(*groupInvite)
	if err != nil {
		return err
	}

	return ctx.JSON(groupInviteDTO)
}

func (c *GroupInviteController) Join(ctx fiber.Ctx) error {
	inviteID := ctx.Params("inviteID")

	authUserID, err := c.authTokenManager.GetAuthUserID(jwtware.FromContext(ctx))
	if err != nil {
		return err
	}

	group, err := c.groupInviteService.JoinGroup(ctx.Context(), inviteID, authUserID)
	if err != nil {
		return err
	}

	groupDTO, err := mapGroupFromDomain(*group)
	if err != nil {
		return err
	}

	return ctx.JSON(groupDTO)
}
