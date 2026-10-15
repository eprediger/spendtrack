package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
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

func dec(s string) *decimal.Decimal {
	d := decimal.RequireFromString(s)
	return &d
}

func TestNewExpenseValid(t *testing.T) {
	now := occurred(t, "2026-09-27T18:00:00Z")
	e, verr := NewExpense(occurred(t, "2026-09-20T14:30:00-03:00"), dec("45.90"), ptr("USD"), "Groceries", nil, now)
	require.Nil(t, verr)
	require.NotNil(t, e)
	assert.True(t, e.Money.Amount.Equal(*dec("45.90")))
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
		amount      *decimal.Decimal
		currency    *string
		description string
		items       []PurchaseItemInput
		wantFields  map[string]string
	}{
		{"negative amount", past, dec("-45.90"), ptr("USD"), "Groceries", nil, map[string]string{"amount": "not_positive"}},
		{"zero amount", past, dec("0"), ptr("USD"), "Groceries", nil, map[string]string{"amount": "not_positive"}},
		{"malformed currency", past, dec("45.90"), ptr("XX1"), "Groceries", nil, map[string]string{"currency": "unsupported_currency"}},
		{"unknown currency", past, dec("45.90"), ptr("ZZZ"), "Groceries", nil, map[string]string{"currency": "unsupported_currency"}},
		{"empty description", past, dec("45.90"), ptr("USD"), "   ", nil, map[string]string{"description": "required"}},
		{"future", future, dec("45.90"), ptr("USD"), "Groceries", nil, map[string]string{"occurred_at": "future"}},
		{"multiple", future, dec("0"), ptr("XYZ"), "", nil, map[string]string{"occurred_at": "future", "amount": "not_positive", "currency": "unsupported_currency", "description": "required"}},
		{"amount missing", past, nil, ptr("USD"), "Groceries", nil, map[string]string{"amount": "required"}},
		{"currency missing", past, dec("45.90"), nil, "Groceries", nil, map[string]string{"currency": "required"}},
		{"amount missing and currency invalid reported together", past, nil, ptr("XX1"), "Groceries", nil, map[string]string{"amount": "required", "currency": "unsupported_currency"}},
		{"currency missing and amount invalid reported together", past, dec("-1"), nil, "Groceries", nil, map[string]string{"amount": "not_positive", "currency": "required"}},
		{"expense and item errors", past, dec("45.90"), ptr("XX1"), "Groceries", []PurchaseItemInput{{Description: "Candies", Quantity: dec("0"), UnitPrice: dec("1"), LineTotal: dec("1")}}, map[string]string{"currency": "unsupported_currency", "items[0].quantity": "not_positive"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, verr := NewExpense(tc.occurredAt, tc.amount, tc.currency, tc.description, tc.items, now)
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

func TestNewExpenseItems(t *testing.T) {
	now := occurred(t, "2026-09-27T18:00:00Z")
	past := occurred(t, "2026-09-20T14:30:00-03:00")
	items := []PurchaseItemInput{
		{Description: "Candies", Quantity: dec("5"), UnitPrice: dec("150.00"), LineTotal: dec("750.00")},
		{Description: "Bottled water", Quantity: dec("1"), UnitPrice: dec("3400.00"), LineTotal: dec("3400.00")},
	}

	e, verr := NewExpense(past, dec("14000.00"), ptr("ARS"), "Gas station store", items, now)
	require.Nil(t, verr)
	require.Len(t, e.Items, 2)
	assert.Equal(t, "Candies", e.Items[0].Description)
	assert.True(t, e.Items[0].Quantity.Equal(*dec("5")))
	assert.True(t, e.Items[0].UnitPrice.Equal(*dec("150.00")))
	assert.True(t, e.Items[0].LineTotal.Equal(*dec("750.00")))
	assert.Equal(t, "Bottled water", e.Items[1].Description)

	withoutItems, verr := NewExpense(past, dec("45.90"), ptr("USD"), "Groceries", nil, now)
	require.Nil(t, verr)
	assert.Empty(t, withoutItems.Items)
}

func TestNewExpenseItemErrors(t *testing.T) {
	now := occurred(t, "2026-09-27T18:00:00Z")
	past := occurred(t, "2026-09-20T14:30:00-03:00")

	_, verr := NewExpense(past, dec("45.90"), ptr("USD"), "Groceries", []PurchaseItemInput{
		{Description: "Candies", Quantity: dec("1"), UnitPrice: dec("1"), LineTotal: dec("1")},
		{Description: "Water", Quantity: dec("0"), UnitPrice: dec("1"), LineTotal: dec("1")},
	}, now)
	require.NotNil(t, verr)
	got := map[string]string{}
	for _, field := range verr.Fields {
		got[field.Field] = field.Code
	}
	assert.Equal(t, map[string]string{"items[1].quantity": "not_positive"}, got)

	_, verr = NewExpense(past, dec("45.90"), ptr("USD"), "Groceries", []PurchaseItemInput{
		{Description: "   ", Quantity: dec("1"), UnitPrice: dec("1"), LineTotal: dec("1")},
		{Description: "Water", Quantity: dec("1"), UnitPrice: dec("1")},
	}, now)
	require.NotNil(t, verr)
	got = map[string]string{}
	for _, field := range verr.Fields {
		got[field.Field] = field.Code
	}
	assert.Equal(t, map[string]string{
		"items[0].description": "required",
		"items[1].line_total":  "required",
	}, got)
}

func TestNewCurrency(t *testing.T) {
	_, err := newCurrency("USD")
	assert.Nil(t, err)
	_, err = newCurrency("XX1")
	require.NotNil(t, err)
	assert.Equal(t, "unsupported_currency", err.Code)
}
