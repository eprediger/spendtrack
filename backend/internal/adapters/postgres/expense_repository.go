package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"spendtrack/internal/domain"
)

// ExpenseRepository is the Postgres adapter for application.ExpenseRepository.
type ExpenseRepository struct {
	pool *pgxpool.Pool
}

func NewExpenseRepository(pool *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{pool: pool}
}

func (r *ExpenseRepository) Save(ctx context.Context, e *domain.Expense) error {
	// Send text, not the Decimal: its driver Value() trims trailing zeros ("45.90" would be stored as "45.9").
	amount := e.Money.Amount.StringFixed(-e.Money.Amount.Exponent())
	_, err := r.pool.Exec(ctx,
		`INSERT INTO expenses (id, occurred_at, occurred_offset_minutes, amount, currency, description)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.OccurredAt, e.OccurredOffsetMinutes, amount, string(e.Money.Currency), e.Description)
	return err
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	var (
		e        domain.Expense
		amount   decimal.Decimal
		currency string
	)
	err := r.pool.QueryRow(ctx,
		`SELECT id, occurred_at, occurred_offset_minutes, amount, currency, description
		 FROM expenses WHERE id = $1`, id).
		Scan(&e.ID, &e.OccurredAt, &e.OccurredOffsetMinutes, &amount, &currency, &e.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	// Same check as new input: a row written outside the API (a seed, a script) must not be served.
	amountText := amount.StringFixed(-amount.Exponent())
	money, errs := domain.NewMoney(&amountText, &currency)
	if len(errs) > 0 {
		return nil, fmt.Errorf("stored expense %s has invalid money: %v", id, errs)
	}
	e.Money = money
	return &e, nil
}
