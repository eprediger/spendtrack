package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMoney exercises the pair contract directly: presence and
// content checks for amount and currency, reported together.
func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   *decimal.Decimal
		currency *string
		expected map[string]string
	}{
		{"both missing", nil, nil, map[string]string{"amount": "required", "currency": "required"}},
		{"amount missing, currency invalid", nil, ptr("XX1"), map[string]string{"amount": "required", "currency": "unsupported_currency"}},
		{"amount not positive, currency missing", dec("-1"), nil, map[string]string{"amount": "not_positive", "currency": "required"}},
		{"amount not positive", dec("-1"), ptr("USD"), map[string]string{"amount": "not_positive"}},
		{"amount zero", dec("0"), ptr("USD"), map[string]string{"amount": "not_positive"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := NewMoney(tc.amount, tc.currency)
			require.Len(t, errs, len(tc.expected))
			for _, fe := range errs {
				assert.Equal(t, tc.expected[fe.Field], fe.Code)
			}
		})
	}
}

func TestNewMoneyValid(t *testing.T) {
	m, errs := NewMoney(dec("45.90"), ptr("ARS"))
	require.Empty(t, errs)
	assert.True(t, m.Amount.Equal(*dec("45.90")))
	assert.Equal(t, Currency("ARS"), m.Currency)
}
