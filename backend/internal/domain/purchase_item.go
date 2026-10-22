package domain

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// ItemDiscountInput contains raw attributes for a discount on a purchase item.
type ItemDiscountInput struct {
	Description string
	Amount      *decimal.Decimal
}

// PurchaseItemInput contains raw attributes for an item inside a new Expense.
type PurchaseItemInput struct {
	Description                    string
	Quantity, UnitPrice, LineTotal *decimal.Decimal
	Discounts                      []ItemDiscountInput
}

// ItemDiscount is a reduction already applied to one purchase item.
type ItemDiscount struct {
	Description string
	Amount      decimal.Decimal
}

// PurchaseItem is one line of a purchase recorded in an expense.
type PurchaseItem struct {
	Description                    string
	Quantity, UnitPrice, LineTotal decimal.Decimal
	Discounts                      []ItemDiscount
}

// NewPurchaseItem builds a valid PurchaseItem or reports every violated field.
func NewPurchaseItem(description string, quantity, unitPrice, lineTotal *decimal.Decimal, discounts []ItemDiscountInput) (PurchaseItem, []FieldError) {
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

	var itemDiscounts []ItemDiscount
	for j, input := range discounts {
		discountDescription := strings.TrimSpace(input.Description)
		if discountDescription == "" {
			errs = append(errs, FieldError{Field: fmt.Sprintf("discounts[%d].description", j), Code: "required", Detail: "description is required"})
		}
		var amount decimal.Decimal
		valid := true
		if input.Amount == nil {
			errs = append(errs, FieldError{Field: fmt.Sprintf("discounts[%d].amount", j), Code: "required", Detail: "amount is required"})
			valid = false
		} else if !input.Amount.IsPositive() {
			errs = append(errs, FieldError{Field: fmt.Sprintf("discounts[%d].amount", j), Code: "not_positive", Detail: "must be greater than 0"})
			valid = false
		} else {
			amount = *input.Amount
		}
		if valid && discountDescription != "" {
			itemDiscounts = append(itemDiscounts, ItemDiscount{Description: discountDescription, Amount: amount})
		}
	}

	if len(errs) > 0 {
		return PurchaseItem{}, errs
	}
	return PurchaseItem{
		Description: description,
		Quantity:    q,
		UnitPrice:   price,
		LineTotal:   total,
		Discounts:   itemDiscounts,
	}, nil
}
