package httpapi

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spendtrack/internal/domain"
)

func text(s string) *string { return &s }

func TestParseAmount(t *testing.T) {
	// A missing amount is the domain's to report as "required".
	amount, ferr := parseAmount(nil)
	assert.Nil(t, amount)
	assert.Nil(t, ferr)

	amount, ferr = parseAmount(text("abc"))
	assert.Nil(t, amount)
	require.NotNil(t, ferr)
	assert.Equal(t, domain.FieldError{Field: "amount", Code: "invalid", Detail: "must be a decimal number"}, *ferr)

	amount, ferr = parseAmount(text("45.90"))
	require.Nil(t, ferr)
	require.NotNil(t, amount)
	assert.True(t, amount.Equal(decimal.RequireFromString("45.90")))
}

// TestFormatAmount locks the scale-preserving rendering: input scale must
// survive, not collapse to "45.9".
func TestFormatAmount(t *testing.T) {
	assert.Equal(t, "45.90", formatAmount(decimal.RequireFromString("45.90")))
	assert.Equal(t, "45", formatAmount(decimal.RequireFromString("45")))
}

func TestWithAmountError(t *testing.T) {
	fields := []domain.FieldError{
		{Field: "amount", Code: "required", Detail: "amount is required"},
		{Field: "currency", Code: "unsupported_currency", Detail: "currency is not supported"},
	}
	invalid := domain.FieldError{Field: "amount", Code: "invalid", Detail: "must be a decimal number"}

	got := withAmountError(fields, &invalid)
	assert.Equal(t, []domain.FieldError{invalid, fields[1]}, got)
	assert.Equal(t, "required", fields[0].Code, "input must not be modified")

	assert.Equal(t, fields, withAmountError(fields, nil))
}
