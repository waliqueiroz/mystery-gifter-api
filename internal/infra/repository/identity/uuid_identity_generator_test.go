package identity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/identity"
)

func Test_UUIDIdentityGenerator_Generate(t *testing.T) {
	t.Run("should generate a UUID successfully", func(t *testing.T) {
		// given
		mockUUID := uuid.New()
		mockedNewUUID := func() (uuid.UUID, error) {
			return mockUUID, nil
		}
		generator := identity.NewUUIDIdentityGenerator(mockedNewUUID)

		// when
		id, err := generator.Generate()

		// then
		assert.NoError(t, err)
		assert.Equal(t, mockUUID.String(), id)
	})

	t.Run("should return an error when UUID generation fails", func(t *testing.T) {
		// given
		mockedNewUUID := func() (uuid.UUID, error) {
			return uuid.UUID{}, assert.AnError
		}
		generator := identity.NewUUIDIdentityGenerator(mockedNewUUID)

		// when
		id, err := generator.Generate()

		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, "", id)
	})
}
