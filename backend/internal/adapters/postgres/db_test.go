package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spendtrack/internal/adapters/postgres"
)

func TestConnectRejectsUnparseableURL(t *testing.T) {
	pool, err := postgres.Connect(context.Background(), "postgres://%zz")
	require.Error(t, err)
	assert.Nil(t, pool)
}

func TestConnectFailsWhenServerIsUnreachable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, "postgres://postgres:postgres@127.0.0.1:1/x?sslmode=disable&connect_timeout=1")
	require.Error(t, err)
	assert.Nil(t, pool)
}
