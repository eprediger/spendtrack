package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func occurred(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)
	return v
}

func ptr[T any](v T) *T { return &v }

func TestNewExpenseValid(t *testing.T) {
	now := occurred(t, "2026-09-27T18:00:00Z")
	e, verr := NewExpense(occurred(t, "2026-09-20T14:30:00-03:00"), ptr("45.90"), ptr("USD"), "Groceries", now)
	require.Nil(t, verr)
	require.NotNil(t, e)
	assert.Equal(t, "45.90", e.Money.String())
	assert.Equal(t, Currency("USD"), e.Money.Currency)
	assert.Equal(t, "Groceries", e.Description)
	assert.Equal(t, -180, e.OccurredOffsetMinutes)
	assert.Equal(t, "2026-09-20T14:30:00-03:00", e.LocalTime().Format(time.RFC3339))
}

func TestNewExpenseValidation(t *testing.T) {
	now := occurred(t, "2026-09-27T18:00:00Z")
	past := occurred(t, "2026-09-20T14:30:00-03:00")
	future := occurred(t, "2026-09-28T00:00:00Z")

	cases := []struct {
		name        string
		occurredAt  time.Time
		amount      *string
		currency    *string
		description string
		wantFields  map[string]string
	}{
		{"negative amount", past, ptr("-45.90"), ptr("USD"), "Groceries", map[string]string{"amount": "not_positive"}},
		{"zero amount", past, ptr("0"), ptr("USD"), "Groceries", map[string]string{"amount": "not_positive"}},
		{"bad decimal", past, ptr("abc"), ptr("USD"), "Groceries", map[string]string{"amount": "invalid"}},
		{"malformed currency", past, ptr("45.90"), ptr("XX1"), "Groceries", map[string]string{"currency": "unsupported_currency"}},
		{"unknown currency", past, ptr("45.90"), ptr("ZZZ"), "Groceries", map[string]string{"currency": "unsupported_currency"}},
		{"empty description", past, ptr("45.90"), ptr("USD"), "   ", map[string]string{"description": "required"}},
		{"future", future, ptr("45.90"), ptr("USD"), "Groceries", map[string]string{"occurred_at": "future"}},
		{"multiple", future, ptr("0"), ptr("XYZ"), "", map[string]string{"occurred_at": "future", "amount": "not_positive", "currency": "unsupported_currency", "description": "required"}},
		{"amount missing", past, nil, ptr("USD"), "Groceries", map[string]string{"amount": "required"}},
		{"currency missing", past, ptr("45.90"), nil, "Groceries", map[string]string{"currency": "required"}},
		{"amount missing and currency invalid reported together", past, nil, ptr("XX1"), "Groceries", map[string]string{"amount": "required", "currency": "unsupported_currency"}},
		{"currency missing and amount invalid reported together", past, ptr("-1"), nil, "Groceries", map[string]string{"amount": "not_positive", "currency": "required"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, verr := NewExpense(tc.occurredAt, tc.amount, tc.currency, tc.description, now)
			assert.Nil(t, e)
			require.NotNil(t, verr)
			got := map[string]string{}
			for _, f := range verr.Fields {
				got[f.Field] = f.Code
			}
			assert.Equal(t, tc.wantFields, got)
		})
	}
}

func TestNewCurrency(t *testing.T) {
	_, err := newCurrency("USD")
	assert.Nil(t, err)
	_, err = newCurrency("XX1")
	require.NotNil(t, err)
	assert.Equal(t, "unsupported_currency", err.Code)
}
