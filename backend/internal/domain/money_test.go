package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMoney exercises the pair contract directly: presence and
// content checks for amount and currency, reported together.
func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   *string
		currency *string
		expected map[string]string
	}{
		{"both missing", nil, nil, map[string]string{"amount": "required", "currency": "required"}},
		{"amount missing, currency invalid", nil, ptr("XX1"), map[string]string{"amount": "required", "currency": "unsupported_currency"}},
		{"amount invalid, currency missing", ptr("abc"), nil, map[string]string{"amount": "invalid", "currency": "required"}},
		{"amount not positive", ptr("-1"), ptr("USD"), map[string]string{"amount": "not_positive"}},
		{"amount zero", ptr("0"), ptr("USD"), map[string]string{"amount": "not_positive"}},
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
	m, errs := NewMoney(ptr("45.90"), ptr("ARS"))
	require.Empty(t, errs)
	assert.Equal(t, "45.90", m.String())
	assert.Equal(t, Currency("ARS"), m.Currency)
}

// TestMoneyString locks the scale-preserving rendering used by the
// HTTP adapter: input scale must survive, not collapse to "45.9".
func TestMoneyString(t *testing.T) {
	m, errs := NewMoney(ptr("45.90"), ptr("USD"))
	require.Empty(t, errs)
	assert.Equal(t, "45.90", m.String())

	m, errs = NewMoney(ptr("45"), ptr("USD"))
	require.Empty(t, errs)
	assert.Equal(t, "45", m.String())
}
