package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemValid(t *testing.T) {
	item, errs := NewItem("  Candies ", dec("5"), dec("150.00"), dec("750.00"))

	require.Empty(t, errs)
	assert.Equal(t, "Candies", item.Description)
	assert.True(t, item.Quantity.Equal(*dec("5")))
	assert.True(t, item.UnitPrice.Equal(*dec("150.00")))
	assert.True(t, item.LineTotal.Equal(*dec("750.00")))
}

func TestNewItemSignedAmounts(t *testing.T) {
	item, errs := NewItem("RETURNABLE BOTTLE", dec("3"), dec("-3001.00"), dec("-9003.00"))
	require.Empty(t, errs)
	assert.True(t, item.UnitPrice.Equal(*dec("-3001.00")))
	assert.True(t, item.LineTotal.Equal(*dec("-9003.00")))

	item, errs = NewItem("Free item", dec("1"), dec("0"), dec("0"))
	require.Empty(t, errs)
	assert.True(t, item.UnitPrice.IsZero())
	assert.True(t, item.LineTotal.IsZero())
}

func TestNewItemValidation(t *testing.T) {
	tests := []struct {
		name        string
		description string
		quantity    *decimal.Decimal
		unitPrice   *decimal.Decimal
		lineTotal   *decimal.Decimal
		want        map[string]string
	}{
		{
			name:        "blank description",
			description: "   ",
			quantity:    dec("1"),
			unitPrice:   dec("1"),
			lineTotal:   dec("1"),
			want:        map[string]string{"description": "required"},
		},
		{
			name:        "quantity missing",
			description: "Candies",
			unitPrice:   dec("1"),
			lineTotal:   dec("1"),
			want:        map[string]string{"quantity": "required"},
		},
		{
			name:        "quantity zero",
			description: "Candies",
			quantity:    dec("0"),
			unitPrice:   dec("1"),
			lineTotal:   dec("1"),
			want:        map[string]string{"quantity": "not_positive"},
		},
		{
			name:        "quantity negative",
			description: "Candies",
			quantity:    dec("-1"),
			unitPrice:   dec("1"),
			lineTotal:   dec("1"),
			want:        map[string]string{"quantity": "not_positive"},
		},
		{
			name:        "unit price missing",
			description: "Candies",
			quantity:    dec("1"),
			lineTotal:   dec("1"),
			want:        map[string]string{"unit_price": "required"},
		},
		{
			name:        "line total missing",
			description: "Candies",
			quantity:    dec("1"),
			unitPrice:   dec("1"),
			want:        map[string]string{"line_total": "required"},
		},
		{
			name: "all fields missing",
			want: map[string]string{
				"description": "required",
				"quantity":    "required",
				"unit_price":  "required",
				"line_total":  "required",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			item, errs := NewItem(tc.description, tc.quantity, tc.unitPrice, tc.lineTotal)
			assert.Empty(t, item)
			require.Len(t, errs, len(tc.want))
			got := make(map[string]string, len(errs))
			for _, err := range errs {
				got[err.Field] = err.Code
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
