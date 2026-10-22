package postgres_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spendtrack/internal/adapters/postgres"
	"spendtrack/internal/domain"
)

// testPool connects to the test database, running migrations. Tests skip when
// DATABASE_URL is not set so the unit suite stays database-free.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	pool, err := postgres.Connect(context.Background(), url)
	require.NoError(t, err)

	var database string
	require.NoError(t, pool.QueryRow(context.Background(), "SELECT current_database()").Scan(&database))
	require.True(t, strings.HasSuffix(database, "_test"), "integration tests require a dedicated database ending in _test; got %q", database)

	t.Cleanup(pool.Close)
	return pool
}

func testExpense() *domain.Expense {
	return &domain.Expense{
		ID:                    uuid.New(),
		OccurredAt:            time.Date(2026, 9, 20, 14, 30, 0, 0, time.UTC),
		OccurredOffsetMinutes: 0,
		Money:                 domain.Money{Amount: decimal.RequireFromString("12.00"), Currency: "USD"},
		Description:           "Lunch",
	}
}

func TestSaveRollsBackWhenAnItemIsRejected(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewExpenseRepository(pool)

	expense := testExpense()
	expense.Items = []domain.PurchaseItem{{
		Description: "Coffee",
		Quantity:    decimal.RequireFromString("0"),
		UnitPrice:   decimal.RequireFromString("12.00"),
		LineTotal:   decimal.RequireFromString("12.00"),
	}}

	require.Error(t, repo.Save(context.Background(), expense))

	var count int
	require.NoError(t, pool.QueryRow(context.Background(),
		"SELECT count(*) FROM expenses WHERE id = $1", expense.ID).Scan(&count))
	assert.Equal(t, 0, count)
}

func TestSaveRollsBackWhenADiscountIsRejected(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewExpenseRepository(pool)

	expense := testExpense()
	expense.Items = []domain.PurchaseItem{{
		Description: "Coffee",
		Quantity:    decimal.RequireFromString("1"),
		UnitPrice:   decimal.RequireFromString("12.00"),
		LineTotal:   decimal.RequireFromString("12.00"),
		Discounts: []domain.ItemDiscount{{
			Description: "PROMO",
			Amount:      decimal.RequireFromString("0"),
		}},
	}}

	require.Error(t, repo.Save(context.Background(), expense))

	var expenseCount, itemCount int
	require.NoError(t, pool.QueryRow(context.Background(),
		"SELECT count(*) FROM expenses WHERE id = $1", expense.ID).Scan(&expenseCount))
	require.NoError(t, pool.QueryRow(context.Background(),
		"SELECT count(*) FROM expense_items WHERE expense_id = $1", expense.ID).Scan(&itemCount))
	assert.Equal(t, 0, expenseCount)
	assert.Equal(t, 0, itemCount)
}
