package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

const authCookieName = "access_token"

func setCookie(ctx fiber.Ctx, token string, expiresIn int64, secure bool) {
	ctx.Cookie(&fiber.Cookie{
		Name:     authCookieName,
		Value:    token,
		Expires:  time.Unix(expiresIn, 0),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	})
}

func clearCookie(ctx fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}
