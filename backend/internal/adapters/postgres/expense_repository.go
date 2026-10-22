package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"spendtrack/internal/application"
	"spendtrack/internal/domain"
)

type expenseRepository struct {
	pool *pgxpool.Pool
}

func NewExpenseRepository(pool *pgxpool.Pool) application.ExpenseRepository {
	return &expenseRepository{pool: pool}
}

func numericText(d decimal.Decimal) string {
	return d.StringFixed(-d.Exponent())
}

func (r *expenseRepository) Save(ctx context.Context, e *domain.Expense) error {
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
		for j, discount := range item.Discounts {
			_, err := tx.Exec(ctx,
				`INSERT INTO expense_item_discounts (expense_id, item_position, position, description, amount)
				 VALUES ($1, $2, $3, $4, $5)`,
				e.ID, i, j, discount.Description, numericText(discount.Amount))
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func (r *expenseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	var (
		e               domain.Expense
		amount          decimal.Decimal
		currency        string
		hasExpense      bool
		itemPosition    *int32
		itemDescription *string
		quantity        *decimal.Decimal
		unitPrice       *decimal.Decimal
		lineTotal       *decimal.Decimal
	)
	var items []domain.PurchaseItemInput
	rows, err := r.pool.Query(ctx,
		`SELECT e.id, e.occurred_at, e.occurred_offset_minutes, e.amount, e.currency, e.description,
		        i.position, i.description, i.quantity, i.unit_price, i.line_total
		 FROM expenses e
		 LEFT JOIN expense_items i ON i.expense_id = e.id
		 WHERE e.id = $1
		 ORDER BY i.position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&e.ID, &e.OccurredAt, &e.OccurredOffsetMinutes, &amount, &currency, &e.Description,
			&itemPosition, &itemDescription, &quantity, &unitPrice, &lineTotal,
		); err != nil {
			return nil, err
		}

		hasExpense = true

		if itemPosition == nil {
			continue
		}
		if itemDescription == nil || quantity == nil || unitPrice == nil || lineTotal == nil {
			return nil, fmt.Errorf("stored expense %s has incomplete purchase item", id)
		}

		items = append(items, domain.PurchaseItemInput{
			Description: *itemDescription,
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			LineTotal:   lineTotal,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !hasExpense {
		return nil, domain.ErrNotFound
	}

	discountRows, err := r.pool.Query(ctx,
		`SELECT item_position, description, amount FROM expense_item_discounts
		 WHERE expense_id = $1 ORDER BY item_position, position`, id)
	if err != nil {
		return nil, err
	}
	defer discountRows.Close()
	for discountRows.Next() {
		var (
			itemPosition        int32
			discountDescription string
			discountAmount      decimal.Decimal
		)
		if err := discountRows.Scan(&itemPosition, &discountDescription, &discountAmount); err != nil {
			return nil, err
		}
		if itemPosition < 0 || int(itemPosition) >= len(items) {
			return nil, fmt.Errorf("stored expense %s has a discount on unknown item position %d", id, itemPosition)
		}
		amount := discountAmount
		items[itemPosition].Discounts = append(items[itemPosition].Discounts, domain.ItemDiscountInput{
			Description: discountDescription,
			Amount:      &amount,
		})
	}
	if err := discountRows.Err(); err != nil {
		return nil, err
	}
	// Restore the registration offset stored separately from timestamptz before domain validation.
	occurredAt := e.OccurredAt.In(time.FixedZone("", e.OccurredOffsetMinutes*60))
	validated, validationErr := domain.NewExpense(
		occurredAt, &amount, &currency, e.Description, items, time.Now(),
	)
	if validationErr != nil {
		return nil, fmt.Errorf("stored expense %s is invalid: %v", id, validationErr)
	}
	validated.ID = e.ID // Keep the persisted identity; NewExpense assigns IDs for new records.
	return validated, nil
}
