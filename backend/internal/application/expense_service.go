package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"spendtrack/internal/domain"
)

// RegisterInput is the raw registration payload. Fields are pointers so an
// adapter can report an omitted field without having decoded a value for it;
// domain.NewExpense turns that into a "required" error alongside every other
// violation in the same request.
type RegisterInput struct {
	OccurredAt  *time.Time
	Amount      *string
	Currency    *string
	Description *string
}

// ExpenseService holds the expense use cases.
type ExpenseService struct {
	repo ExpenseRepository
	now  func() time.Time
}

func NewExpenseService(repo ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo, now: time.Now}
}

// Register validates and persists a new expense.
func (s *ExpenseService) Register(ctx context.Context, in RegisterInput) (*domain.Expense, error) {
	desc := ""
	if in.Description != nil {
		desc = *in.Description
	}
	e, verr := domain.NewExpense(*in.OccurredAt, in.Amount, in.Currency, desc, s.now())
	if verr != nil {
		return nil, verr
	}
	if err := s.repo.Save(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

// Find returns one expense or domain.ErrNotFound.
func (s *ExpenseService) Find(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	return s.repo.FindByID(ctx, id)
}
