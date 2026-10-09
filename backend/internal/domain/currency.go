package domain

// Currency is an ISO 4217 alphabetic currency code.
type Currency string

// supportedCurrencies is the curated set of codes the product accepts.
// It grows deliberately as needs appear; every code here is a valid
// ISO 4217 alphabetic code. Keeping a curated set instead of the full
// ISO list avoids maintaining (or depending on) a large data table that
// drifts.
var supportedCurrencies = []string{"ARS", "USD"}

// newCurrency validates code against the supported set.
func newCurrency(code string) (Currency, *FieldError) {
	for _, supported := range supportedCurrencies {
		if code == supported {
			return Currency(code), nil
		}
	}
	return "", &FieldError{Field: "currency", Code: "unsupported_currency", Detail: "currency is not supported"}
}
