package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Expense is an economic consumption recognized as a financial fact.
type Expense struct {
	ID                    uuid.UUID
	OccurredAt            time.Time
	OccurredOffsetMinutes int
	Money                 Money
	Description           string
}

// NewExpense builds a valid Expense or reports every violated field.
// occurredAt must carry the local UTC offset observed by the client;
// now is the instant registration happens (server clock). occurredAt and
// description are concrete values because the adapter is responsible for
// their presence; amount and currency are pointers because Money owns the
// presence and content checks of its pair (a Money is never valid with only
// one side present).
func NewExpense(occurredAt time.Time, amount, currency *string, description string, now time.Time) (*Expense, *ValidationError) {
	var fe []FieldError

	if occurredAt.After(now) {
		fe = append(fe, FieldError{Field: "occurred_at", Code: "future", Detail: "occurrence cannot be in the future"})
	}

	money, merrs := NewMoney(amount, currency)
	fe = append(fe, merrs...)

	description = strings.TrimSpace(description)
	if description == "" {
		fe = append(fe, FieldError{Field: "description", Code: "required", Detail: "description is required"})
	}

	if len(fe) > 0 {
		return nil, &ValidationError{Fields: fe}
	}

	_, offsetSeconds := occurredAt.Zone()
	return &Expense{
		ID:                    uuid.Must(uuid.NewV7()),
		OccurredAt:            occurredAt,
		OccurredOffsetMinutes: offsetSeconds / 60,
		Money:                 money,
		Description:           description,
	}, nil
}

// LocalTime returns the occurrence in the offset observed at registration.
func (e *Expense) LocalTime() time.Time {
	return e.OccurredAt.In(time.FixedZone("", e.OccurredOffsetMinutes*60))
}
