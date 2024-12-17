package entrypoint

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

func CustomErrorHandler(ctx fiber.Ctx, err error) error {

	switch e := err.(type) {
	case *fiber.Error:
		return sendError(ctx, e.Code, e.Error())
	case domain.CustomError:
		return sendError(ctx, e.StatusCode(), err.Error())
	default:
		return sendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

}

func sendError(ctx fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"code":    strings.ReplaceAll(strings.ToLower(http.StatusText(statusCode)), " ", "_"),
		"message": message,
	})
}
