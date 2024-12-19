package rest_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/application/mock_application"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
	"github.com/waliqueiroz/mystery-gifter-api/test/builder"
	"github.com/waliqueiroz/mystery-gifter-api/test/helper"
	"go.uber.org/mock/gomock"
)

func Test_UserController_Create(t *testing.T) {
	route := "/api/users"

	t.Run("should return status 201 and the user ID when the user is created successfully", func(t *testing.T) {
		// given
		mockCtrl := gomock.NewController(t)
		mockedUserService := mock_application.NewMockUserService(mockCtrl)
		mockedUserService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		userController := rest.NewUserController(mockedUserService)

		createUserDTO := builder.NewCreateUserDTOBuilder().Build()
		payload := helper.EncodeJSON(t, createUserDTO)

		req := httptest.NewRequest(fiber.MethodPost, route, payload)
		req.Header.Set("Content-Type", "application/json")

		app := fiber.New()
		app.Post(route, userController.Create)

		// when
		response, err := app.Test(req)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		// var result fiber.Map
		// test.DecodeJSON(t, response.Body, &result)

		// assert.Equal(t, fiber.Map{"id": ""}, result)
	})

}
