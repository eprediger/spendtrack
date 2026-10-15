package domain

import (
	"strings"

	"github.com/shopspring/decimal"
)

// Item is one line in an expense breakdown.
type Item struct {
	Description                    string
	Quantity, UnitPrice, LineTotal decimal.Decimal
}

// NewItem builds a valid Item or reports every violated field.
func NewItem(description string, quantity, unitPrice, lineTotal *decimal.Decimal) (Item, []FieldError) {
	var errs []FieldError

	description = strings.TrimSpace(description)
	if description == "" {
		errs = append(errs, FieldError{Field: "description", Code: "required", Detail: "description is required"})
	}

	var q decimal.Decimal
	if quantity == nil {
		errs = append(errs, FieldError{Field: "quantity", Code: "required", Detail: "quantity is required"})
	} else if !quantity.IsPositive() {
		errs = append(errs, FieldError{Field: "quantity", Code: "not_positive", Detail: "must be greater than 0"})
	} else {
		q = *quantity
	}

	var price decimal.Decimal
	if unitPrice == nil {
		errs = append(errs, FieldError{Field: "unit_price", Code: "required", Detail: "unit_price is required"})
	} else {
		price = *unitPrice
	}

	var total decimal.Decimal
	if lineTotal == nil {
		errs = append(errs, FieldError{Field: "line_total", Code: "required", Detail: "line_total is required"})
	} else {
		total = *lineTotal
	}

	if len(errs) > 0 {
		return Item{}, errs
	}
	return Item{
		Description: description,
		Quantity:    q,
		UnitPrice:   price,
		LineTotal:   total,
	}, nil
}
