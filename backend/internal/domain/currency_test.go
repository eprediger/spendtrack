package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSupportedCurrencies asserts every code in the curated set validates.
func TestSupportedCurrencies(t *testing.T) {
	for _, code := range supportedCurrencies {
		c, err := newCurrency(code)
		require.Nil(t, err, "code %s rejected", code)
		assert.Equal(t, Currency(code), c)
	}
}

// TestUnsupportedCurrency asserts a valid ISO code outside the curated set is
// rejected as unsupported (not as malformed).
func TestUnsupportedCurrency(t *testing.T) {
	_, err := newCurrency("EUR")
	require.NotNil(t, err)
	assert.Equal(t, "unsupported_currency", err.Code)
}

// TestMalformedCurrency asserts a non-ISO-shaped code is rejected.
func TestMalformedCurrency(t *testing.T) {
	_, err := newCurrency("XX1")
	require.NotNil(t, err)
	assert.Equal(t, "unsupported_currency", err.Code)
}
