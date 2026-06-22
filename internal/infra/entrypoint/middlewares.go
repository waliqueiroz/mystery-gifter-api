package entrypoint

import (
	"errors"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func NewAuthMiddleware(secretKey string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secretKey)},
		Extractor: extractors.Chain(
			extractors.FromCookie(authCookieName),
			extractors.FromAuthHeader("Bearer"),
		),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if errors.Is(err, extractors.ErrNotFound) {
				return fiber.NewError(fiber.StatusBadRequest, jwtware.ErrMissingToken.Error())
			}
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired JWT")
		},
	})
}
