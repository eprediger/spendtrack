package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	oldValue, wasSet := os.LookupEnv("DATABASE_URL")
	require.NoError(t, os.Unsetenv("DATABASE_URL"))
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv("DATABASE_URL", oldValue)
		} else {
			_ = os.Unsetenv("DATABASE_URL")
		}
	})

	_, err := Load()

	require.ErrorContains(t, err, "DATABASE_URL")
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/spendtrack_test?sslmode=disable")
	t.Setenv("PORT", "not-a-number")

	_, err := Load()

	require.ErrorContains(t, err, "Port")
}

func TestLoadReadsEnvironmentValues(t *testing.T) {
	const databaseURL = "postgres://postgres:postgres@localhost:5432/spendtrack_test?sslmode=disable"
	t.Setenv("DATABASE_URL", databaseURL)
	t.Setenv("PORT", "9000")

	got, err := Load()

	require.NoError(t, err)
	require.Equal(t, Config{Port: 9000, DatabaseURL: databaseURL}, got)
}
