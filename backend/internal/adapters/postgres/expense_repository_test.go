package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spendtrack/internal/adapters/postgres"
	"spendtrack/internal/application"
	"spendtrack/internal/domain"
)

var _ func(*pgxpool.Pool) application.ExpenseRepository = postgres.NewExpenseRepository

func TestNewExpenseRepositoryReturnsApplicationPort(t *testing.T) {
	repo := postgres.NewExpenseRepository(nil)
	if repo == nil {
		t.Fatal("NewExpenseRepository returned nil")
	}
}

func dec(s string) *decimal.Decimal {
	d := decimal.RequireFromString(s)
	return &d
}

func TestExpenseRepositoryDiscountsRoundTrip(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := postgres.Connect(ctx, databaseURL)
	require.NoError(t, err)
	defer pool.Close()
	_, err = pool.Exec(ctx, "TRUNCATE expense_item_discounts, expense_items, expenses")
	require.NoError(t, err)

	repo := postgres.NewExpenseRepository(pool)
	now := time.Now()
	expense, verr := domain.NewExpense(
		now.Add(-time.Hour), dec("347.00"), ptr("ARS"), "Store purchase",
		[]domain.PurchaseItemInput{
			{
				Description: "Soft drink", Quantity: dec("1"), UnitPrice: dec("2000.00"), LineTotal: dec("2000.00"),
				Discounts: []domain.ItemDiscountInput{
					{Description: "SECOND AT 40%", Amount: dec("1700.00")},
					{Description: "COUPON", Amount: dec("3.00")},
				},
			},
			{Description: "Bread", Quantity: dec("1"), UnitPrice: dec("50.00"), LineTotal: dec("50.00")},
		}, now)
	require.Nil(t, verr)
	require.NoError(t, repo.Save(ctx, expense))

	loaded, err := repo.FindByID(ctx, expense.ID)
	require.NoError(t, err)
	require.Len(t, loaded.Items, 2)

	require.Len(t, loaded.Items[0].Discounts, 2)
	assert.Equal(t, "SECOND AT 40%", loaded.Items[0].Discounts[0].Description)
	assert.Equal(t, "1700.00", loaded.Items[0].Discounts[0].Amount.StringFixed(2))
	assert.Equal(t, "COUPON", loaded.Items[0].Discounts[1].Description)
	assert.Equal(t, "3.00", loaded.Items[0].Discounts[1].Amount.StringFixed(2))

	assert.Empty(t, loaded.Items[1].Discounts)
}

func ptr[T any](v T) *T { return &v }
