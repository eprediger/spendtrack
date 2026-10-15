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

func numericText(d decimal.Decimal) string {
	return d.StringFixed(-d.Exponent())
}

func (r *ExpenseRepository) Save(ctx context.Context, e *domain.Expense) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx,
		`INSERT INTO expenses (id, occurred_at, occurred_offset_minutes, amount, currency, description)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.OccurredAt, e.OccurredOffsetMinutes, numericText(e.Money.Amount), string(e.Money.Currency), e.Description)
	if err != nil {
		return err
	}
	for i, item := range e.Items {
		_, err := tx.Exec(ctx,
			`INSERT INTO expense_items (expense_id, position, description, quantity, unit_price, line_total)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			e.ID, i, item.Description, numericText(item.Quantity), numericText(item.UnitPrice), numericText(item.LineTotal))
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
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
	money, errs := domain.NewMoney(&amount, &currency)
	if len(errs) > 0 {
		return nil, fmt.Errorf("stored expense %s has invalid money: %v", id, errs)
	}
	e.Money = money

	rows, err := r.pool.Query(ctx,
		`SELECT description, quantity, unit_price, line_total
		 FROM expense_items WHERE expense_id = $1 ORDER BY position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			description string
			quantity    decimal.Decimal
			unitPrice   decimal.Decimal
			lineTotal   decimal.Decimal
		)
		if err := rows.Scan(&description, &quantity, &unitPrice, &lineTotal); err != nil {
			return nil, err
		}
		item, errs := domain.NewPurchaseItem(description, &quantity, &unitPrice, &lineTotal)
		if len(errs) > 0 {
			for i := range errs {
				errs[i].Field = fmt.Sprintf("items[%d].%s", len(e.Items), errs[i].Field)
			}
			return nil, fmt.Errorf("stored expense %s has invalid purchase item: %v", id, errs)
		}
		e.Items = append(e.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &e, nil
}
