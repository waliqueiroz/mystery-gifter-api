package entrypoint

import (
	"errors"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
)

func NewAuthMiddleware(secretKey string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secretKey)},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if errors.Is(err, jwtware.ErrMissingToken) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired JWT")
		},
	})
}
