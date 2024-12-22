package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/config"
)

func TestLoad(t *testing.T) {
	t.Run("should load configuration successfully", func(t *testing.T) {
		// given
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_DATABASE", "test_db")
		os.Setenv("DB_USERNAME", "test_user")
		os.Setenv("DB_PASSWORD", "test_pass")
		defer os.Clearenv()

		// when
		cfg, err := config.Load()

		// then
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, "5432", cfg.Database.Port)
		assert.Equal(t, "test_db", cfg.Database.Database)
		assert.Equal(t, "test_user", cfg.Database.Username)
		assert.Equal(t, "test_pass", cfg.Database.Password)
	})

	t.Run("should return an error if environment variables are missing", func(t *testing.T) {
		// given
		os.Clearenv()

		// when
		cfg, err := config.Load()

		// then
		assert.Nil(t, cfg)
		assert.Error(t, err)
	})
}
