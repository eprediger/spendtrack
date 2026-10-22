package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPurchaseItemValid(t *testing.T) {
	item, errs := NewPurchaseItem("  Candies ", dec("5"), dec("150.00"), dec("750.00"), nil)

	require.Empty(t, errs)
	assert.Equal(t, "Candies", item.Description)
	assert.True(t, item.Quantity.Equal(*dec("5")))
	assert.True(t, item.UnitPrice.Equal(*dec("150.00")))
	assert.True(t, item.LineTotal.Equal(*dec("750.00")))
	assert.Empty(t, item.Discounts)
}

func TestNewPurchaseItemSignedAmounts(t *testing.T) {
	item, errs := NewPurchaseItem("RETURNABLE BOTTLE", dec("3"), dec("-3001.00"), dec("-9003.00"), nil)
	require.Empty(t, errs)
	assert.True(t, item.UnitPrice.Equal(*dec("-3001.00")))
	assert.True(t, item.LineTotal.Equal(*dec("-9003.00")))

	item, errs = NewPurchaseItem("Free item", dec("1"), dec("0"), dec("0"), nil)
	require.Empty(t, errs)
	assert.True(t, item.UnitPrice.IsZero())
	assert.True(t, item.LineTotal.IsZero())
}

func TestNewPurchaseItemDiscounts(t *testing.T) {
	item, errs := NewPurchaseItem("COCA COLA R 2L", dec("1"), dec("12750.00"), dec("12750.00"), []ItemDiscountInput{
		{Description: "  2DO AL 40% LINEA COC ", Amount: dec("1700.00")},
	})
	require.Empty(t, errs)
	require.Len(t, item.Discounts, 1)
	assert.Equal(t, "2DO AL 40% LINEA COC", item.Discounts[0].Description)
	assert.True(t, item.Discounts[0].Amount.Equal(*dec("1700.00")))

	tests := []struct {
		name      string
		discounts []ItemDiscountInput
		want      map[string]string
	}{
		{
			name:      "blank description",
			discounts: []ItemDiscountInput{{Description: "   ", Amount: dec("1.00")}},
			want:      map[string]string{"discounts[0].description": "required"},
		},
		{
			name:      "amount missing",
			discounts: []ItemDiscountInput{{Description: "Promo"}},
			want:      map[string]string{"discounts[0].amount": "required"},
		},
		{
			name:      "amount zero",
			discounts: []ItemDiscountInput{{Description: "Promo", Amount: dec("0")}},
			want:      map[string]string{"discounts[0].amount": "not_positive"},
		},
		{
			name: "amount negative",
			discounts: []ItemDiscountInput{
				{Description: "Promo", Amount: dec("-1.00")},
				{Description: "Promo 2", Amount: dec("1.00")},
			},
			want: map[string]string{"discounts[0].amount": "not_positive"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			item, errs := NewPurchaseItem("COCA COLA R 2L", dec("1"), dec("10.00"), dec("10.00"), tc.discounts)
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

func TestNewPurchaseItemValidation(t *testing.T) {
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
			item, errs := NewPurchaseItem(tc.description, tc.quantity, tc.unitPrice, tc.lineTotal, nil)
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
