package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
)

type AuthController struct {
	authService application.AuthService
}

func NewAuthController(authService application.AuthService) *AuthController {
	return &AuthController{
		authService,
	}
}

func (c *AuthController) Login(ctx fiber.Ctx) error {
	var credentialsDTO CredentialsDTO

	if err := ctx.Bind().Body(&credentialsDTO); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	credentials, err := mapCredentialsToDomain(credentialsDTO)
	if err != nil {
		return err
	}

	authSession, err := c.authService.Login(ctx.Context(), *credentials)
	if err != nil {
		return err
	}

	authSessionDTO, err := mapAuthSessionFromDomain(*authSession)
	if err != nil {
		return err
	}

	return ctx.JSON(authSessionDTO)
}
