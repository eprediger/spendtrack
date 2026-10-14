package domain

import "github.com/shopspring/decimal"

// Money pairs an exact decimal amount with its currency.
type Money struct {
	Amount   decimal.Decimal
	Currency Currency
}

// NewMoney validates an amount/currency pair, keeping both required and
// content checks together since a Money is never valid with only one side
// present. amount and currency are nil when the caller's input omitted them
// (or, for amount, could not be parsed: the adapter reports that case).
func NewMoney(amount *decimal.Decimal, currency *string) (Money, []FieldError) {
	var errs []FieldError

	var d decimal.Decimal
	if amount == nil {
		errs = append(errs, FieldError{Field: "amount", Code: "required", Detail: "amount is required"})
	} else if !amount.IsPositive() {
		errs = append(errs, FieldError{Field: "amount", Code: "not_positive", Detail: "must be greater than 0"})
	} else {
		d = *amount
	}

	var c Currency
	if currency == nil {
		errs = append(errs, FieldError{Field: "currency", Code: "required", Detail: "currency is required"})
	} else if v, cerr := newCurrency(*currency); cerr != nil {
		errs = append(errs, *cerr)
	} else {
		c = v
	}

	if len(errs) > 0 {
		return Money{}, errs
	}
	return Money{Amount: d, Currency: c}, nil
}
