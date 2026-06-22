package entrypoint_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint"
)

const testSecretKey = "test-secret-key"

func makeTestToken(t *testing.T, secretKey string, expiresIn time.Duration) string {
	t.Helper()
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(expiresIn).Unix(),
		"userID":     "some-user-id",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func Test_NewAuthMiddleware(t *testing.T) {
	successHandler := func(ctx fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	}

	t.Run("should authenticate successfully via cookie", func(t *testing.T) {
		// given
		token := makeTestToken(t, testSecretKey, time.Hour)

		app := fiber.New(fiber.Config{ErrorHandler: entrypoint.CustomErrorHandler})
		app.Use(entrypoint.NewAuthMiddleware(testSecretKey))
		app.Get("/protected", successHandler)

		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		req.Header.Set("Cookie", fmt.Sprintf("access_token=%s", token))

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})

	t.Run("should authenticate successfully via authorization header when cookie is absent", func(t *testing.T) {
		// given
		token := makeTestToken(t, testSecretKey, time.Hour)

		app := fiber.New(fiber.Config{ErrorHandler: entrypoint.CustomErrorHandler})
		app.Use(entrypoint.NewAuthMiddleware(testSecretKey))
		app.Get("/protected", successHandler)

		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})

	t.Run("should reject when cookie contains invalid token without falling back to header", func(t *testing.T) {
		// given
		validToken := makeTestToken(t, testSecretKey, time.Hour)

		app := fiber.New(fiber.Config{ErrorHandler: entrypoint.CustomErrorHandler})
		app.Use(entrypoint.NewAuthMiddleware(testSecretKey))
		app.Get("/protected", successHandler)

		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		req.Header.Set("Cookie", "access_token=invalid.token.value")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})

	t.Run("should return bad_request when neither cookie nor authorization header is present", func(t *testing.T) {
		// given
		app := fiber.New(fiber.Config{ErrorHandler: entrypoint.CustomErrorHandler})
		app.Use(entrypoint.NewAuthMiddleware(testSecretKey))
		app.Get("/protected", successHandler)

		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})
}
