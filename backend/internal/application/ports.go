package application

import (
	"context"

	"github.com/google/uuid"

	"spendtrack/internal/domain"
)

// ExpenseRepository is the driven port for expense persistence.
type ExpenseRepository interface {
	Save(ctx context.Context, e *domain.Expense) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
}
