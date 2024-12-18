package entrypoint

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type webError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func CustomErrorHandler(ctx fiber.Ctx, err error) error {

	switch e := err.(type) {
	case *fiber.Error:
		return sendError(ctx, e.Code, e.Error(), nil)
	case domain.CustomError:
		return sendError(ctx, e.StatusCode(), e.Error(), e.Details())
	default:
		return sendError(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

}

func sendError(ctx fiber.Ctx, statusCode int, message string, details any) error {
	return ctx.Status(statusCode).JSON(webError{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(statusCode)), " ", "_"),
		Message: message,
		Details: details,
	})
}
