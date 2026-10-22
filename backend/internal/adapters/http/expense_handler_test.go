package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spendtrack/internal/application"
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

func TestParseItems(t *testing.T) {
	items, errs := parseItems(nil)
	assert.Empty(t, items)
	assert.Empty(t, errs)

	items, errs = parseItems([]itemRequest{
		{Description: text("Candies"), Quantity: text("5"), UnitPrice: text("150.00"), LineTotal: text("750.00")},
		{Description: text("Water"), Quantity: text("abc"), UnitPrice: text("3400.00"), LineTotal: text("3400.00")},
	})
	require.Len(t, items, 2)
	assert.Nil(t, items[1].Quantity)
	require.Equal(t, []domain.FieldError{
		{Field: "items[1].quantity", Code: "invalid", Detail: "must be a decimal number"},
	}, errs)

	items, errs = parseItems([]itemRequest{{Description: text("Water"), Quantity: text("0.750"), LineTotal: text("2400.00")}})
	require.Empty(t, errs)
	require.Len(t, items, 1)
	assert.Nil(t, items[0].UnitPrice)
	assert.Equal(t, int32(-3), items[0].Quantity.Exponent())
}

func TestParseItemsDiscounts(t *testing.T) {
	items, errs := parseItems([]itemRequest{
		{
			Description: text("COCA COLA R 2L"), Quantity: text("1"),
			UnitPrice: text("12750.00"), LineTotal: text("12750.00"),
			Discounts: []discountRequest{
				{Description: text("2DO AL 40% LINEA COC"), Amount: text("1700.00")},
				{Description: text("PROMO"), Amount: text("abc")},
			},
		},
	})
	require.Len(t, items, 1)
	require.Len(t, items[0].Discounts, 2)
	assert.Equal(t, "2DO AL 40% LINEA COC", items[0].Discounts[0].Description)
	require.NotNil(t, items[0].Discounts[0].Amount)
	assert.True(t, items[0].Discounts[0].Amount.Equal(decimal.RequireFromString("1700.00")))
	assert.Nil(t, items[0].Discounts[1].Amount)
	require.Equal(t, []domain.FieldError{
		{Field: "items[0].discounts[1].amount", Code: "invalid", Detail: "must be a decimal number"},
	}, errs)
}

// TestFormatAmount locks the scale-preserving rendering: input scale must
// survive, not collapse to "45.9".
func TestFormatAmount(t *testing.T) {
	assert.Equal(t, "45.90", formatAmount(decimal.RequireFromString("45.90")))
	assert.Equal(t, "45", formatAmount(decimal.RequireFromString("45")))
}

func TestWithParseErrors(t *testing.T) {
	fields := []domain.FieldError{
		{Field: "amount", Code: "required", Detail: "amount is required"},
		{Field: "items[1].quantity", Code: "required", Detail: "quantity is required"},
		{Field: "currency", Code: "unsupported_currency", Detail: "currency is not supported"},
	}
	parseErrors := []domain.FieldError{
		{Field: "amount", Code: "invalid", Detail: "must be a decimal number"},
		{Field: "items[1].quantity", Code: "invalid", Detail: "must be a decimal number"},
	}

	got := withParseErrors(fields, parseErrors)
	assert.Equal(t, []domain.FieldError{parseErrors[0], parseErrors[1], fields[2]}, got)
	assert.Equal(t, "required", fields[0].Code, "input must not be modified")
	assert.Equal(t, "required", fields[1].Code, "input must not be modified")

	assert.Equal(t, fields, withParseErrors(fields, nil))
}

func TestWithParseErrorsAppendsParseErrorsWithoutDomainMatch(t *testing.T) {
	fields := []domain.FieldError{
		{Field: "currency", Code: "unsupported_currency", Detail: "currency is not supported"},
	}
	parseErrors := []domain.FieldError{
		{Field: "items[0].quantity", Code: "invalid", Detail: "must be a decimal number"},
	}

	got := withParseErrors(fields, parseErrors)
	assert.Equal(t, []domain.FieldError{fields[0], parseErrors[0]}, got)
}

func TestWithParseErrorsHandlesManyInvalidItemsQuickly(t *testing.T) {
	const count = 20_000
	fields := make([]domain.FieldError, count)
	parseErrors := make([]domain.FieldError, count)
	for i := range fields {
		field := fmt.Sprintf("items[%d].quantity", i)
		fields[i] = domain.FieldError{Field: field, Code: "required", Detail: "quantity is required"}
		parseErrors[i] = domain.FieldError{Field: field, Code: "invalid", Detail: "must be a decimal number"}
	}

	start := time.Now()
	got := withParseErrors(fields, parseErrors)
	elapsed := time.Since(start)

	require.Len(t, got, count)
	assert.Equal(t, parseErrors[0], got[0])
	assert.Equal(t, parseErrors[count-1], got[count-1])
	assert.Less(t, elapsed, 100*time.Millisecond, "error merging should scale to large item lists")
}

func TestToResponseItems(t *testing.T) {
	expense := &domain.Expense{}
	response := toResponse(expense)
	require.NotNil(t, response.Items)
	assert.Empty(t, response.Items)

	expense.Items = []domain.PurchaseItem{{
		Description: "Apples",
		Quantity:    decimal.RequireFromString("0.750"),
		UnitPrice:   decimal.RequireFromString("3200.00"),
		LineTotal:   decimal.RequireFromString("2400.00"),
	}}
	response = toResponse(expense)
	require.Len(t, response.Items, 1)
	assert.Equal(t, itemResponse{
		Description: "Apples",
		Quantity:    "0.750",
		UnitPrice:   "3200.00",
		LineTotal:   "2400.00",
		Discounts:   []discountResponse{},
	}, response.Items[0])
	assert.NotNil(t, response.Items[0].Discounts)

	expense.Items[0].Discounts = []domain.ItemDiscount{{
		Description: "2DO AL 40% LINEA COC",
		Amount:      decimal.RequireFromString("1700.00"),
	}}
	response = toResponse(expense)
	require.Equal(t, []discountResponse{
		{Amount: "1700.00", Description: "2DO AL 40% LINEA COC"},
	}, response.Items[0].Discounts)
}

func TestParseItemsNullItem(t *testing.T) {
	var req createRequest
	require.NoError(t, json.Unmarshal([]byte(`{"items":[null]}`), &req))

	items, errs := parseItems(req.Items)
	require.Empty(t, errs)
	require.Len(t, items, 1)
	assert.Empty(t, items[0].Description)
	assert.Nil(t, items[0].Quantity)
	assert.Nil(t, items[0].UnitPrice)
	assert.Nil(t, items[0].LineTotal)
}

func TestCreateItemsNotAnArray(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", strings.NewReader(`{"items":"x"}`))
	recorder := httptest.NewRecorder()

	NewExpenseHandler(nil).Create(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/problem+json", recorder.Header().Get("Content-Type"))
	var got problem
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
	assert.Equal(t, problemType+"malformed-json", got.Type)
	assert.Equal(t, "Malformed JSON", got.Title)
	assert.Equal(t, http.StatusBadRequest, got.Status)
}

func TestCreateRejectsOccurredAtWithoutOffset(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", strings.NewReader(
		`{"occurred_at":"2026-09-20T14:30:00","amount":"12.00","currency":"USD","description":"Lunch"}`,
	))
	recorder := httptest.NewRecorder()

	NewExpenseHandler(nil).Create(recorder, req)

	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Equal(t, "application/problem+json", recorder.Header().Get("Content-Type"))
	var got problem
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
	assert.Equal(t, problemType+"validation-failed", got.Type)
	assert.Equal(t, []domain.FieldError{
		{Field: "occurred_at", Code: "invalid", Detail: "must be RFC 3339 with explicit offset"},
	}, got.Errors)
}

func TestCreateRejectsTrailingJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", strings.NewReader(`{} {}`))
	recorder := httptest.NewRecorder()

	NewExpenseHandler(nil).Create(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

type failingExpenseRepository struct {
	err error
}

func (r failingExpenseRepository) Save(context.Context, *domain.Expense) error {
	return r.err
}

func (r failingExpenseRepository) FindByID(context.Context, uuid.UUID) (*domain.Expense, error) {
	return nil, r.err
}

func TestCreateLogsRegistrationFailureWithoutExposingIt(t *testing.T) {
	persistenceErr := errors.New("database unavailable")
	handler := NewExpenseHandler(application.NewExpenseService(failingExpenseRepository{err: persistenceErr}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", strings.NewReader(
		`{"occurred_at":"2020-01-01T00:00:00Z","amount":"12.00","currency":"USD","description":"Lunch"}`,
	))
	recorder := httptest.NewRecorder()
	var logs bytes.Buffer
	oldLogOutput := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(oldLogOutput)

	handler.Create(recorder, req)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, logs.String(), persistenceErr.Error())
	assert.NotContains(t, recorder.Body.String(), persistenceErr.Error())
}
