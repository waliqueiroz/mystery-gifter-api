package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application"
)

type AuthController struct {
	authService  application.AuthService
	cookieSecure bool
}

func NewAuthController(authService application.AuthService, cookieSecure bool) *AuthController {
	return &AuthController{
		authService:  authService,
		cookieSecure: cookieSecure,
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

	setCookie(ctx, authSession.AccessToken, authSession.ExpiresIn, c.cookieSecure)

	return ctx.JSON(authSessionDTO)
}
