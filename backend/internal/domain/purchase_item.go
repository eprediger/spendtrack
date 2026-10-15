package domain

import (
	"strings"

	"github.com/shopspring/decimal"
)

// PurchaseItemInput contains raw attributes for an item inside a new Expense.
type PurchaseItemInput struct {
	Description                    string
	Quantity, UnitPrice, LineTotal *decimal.Decimal
}

// PurchaseItem is one line of a purchase recorded in an expense.
type PurchaseItem struct {
	Description                    string
	Quantity, UnitPrice, LineTotal decimal.Decimal
}

// NewPurchaseItem builds a valid PurchaseItem or reports every violated field.
func NewPurchaseItem(description string, quantity, unitPrice, lineTotal *decimal.Decimal) (PurchaseItem, []FieldError) {
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
		return PurchaseItem{}, errs
	}
	return PurchaseItem{
		Description: description,
		Quantity:    q,
		UnitPrice:   price,
		LineTotal:   total,
	}, nil
}
