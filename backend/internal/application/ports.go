package application

import (
	"context"

	"github.com/google/uuid"

	"spendtrack/internal/domain"
)

// ExpenseRepository is the driven port for expense persistence.
type ExpenseRepository interface {
	// Save persists an expense.
	Save(ctx context.Context, e *domain.Expense) error

	// FindByID returns the expense with the given ID, or domain.ErrNotFound if none exists.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
}
