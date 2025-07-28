package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Test with default values
	cfg, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check default values
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "password", cfg.Database.Password)
	assert.Equal(t, "users_db", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "your-secret-key-here", cfg.JWT.Secret)
	assert.Equal(t, 24, cfg.JWT.ExpiryHours)
}

func TestNewWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_HOST", "db.example.com")
	os.Setenv("DATABASE_PORT", "5433")
	os.Setenv("DATABASE_USER", "testuser")
	os.Setenv("DATABASE_PASSWORD", "testpass")
	os.Setenv("DATABASE_NAME", "testdb")
	os.Setenv("DATABASE_SSLMODE", "require")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRY_HOURS", "48")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("DATABASE_PORT")
		os.Unsetenv("DATABASE_USER")
		os.Unsetenv("DATABASE_PASSWORD")
		os.Unsetenv("DATABASE_NAME")
		os.Unsetenv("DATABASE_SSLMODE")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRY_HOURS")
	}()

	cfg, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check environment variable values
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Equal(t, 48, cfg.JWT.ExpiryHours)
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "users_db",
			SSLMode:  "disable",
		},
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=users_db sslmode=disable"
	assert.Equal(t, expected, cfg.GetDSN())
}

func TestGetServerAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	expected := "localhost:8080"
	assert.Equal(t, expected, cfg.GetServerAddress())
}
