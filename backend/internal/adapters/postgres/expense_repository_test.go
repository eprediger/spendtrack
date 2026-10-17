package postgres_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"spendtrack/internal/adapters/postgres"
	"spendtrack/internal/application"
)

var _ func(*pgxpool.Pool) application.ExpenseRepository = postgres.NewExpenseRepository

func TestNewExpenseRepositoryReturnsApplicationPort(t *testing.T) {
	repo := postgres.NewExpenseRepository(nil)
	if repo == nil {
		t.Fatal("NewExpenseRepository returned nil")
	}
}
