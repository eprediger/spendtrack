package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"spendtrack/internal/application"
	"spendtrack/internal/domain"
)

type ExpenseHandler struct {
	svc *application.ExpenseService
}

func NewExpenseHandler(svc *application.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

type createRequest struct {
	OccurredAt  *string       `json:"occurred_at"`
	Amount      *string       `json:"amount"`
	Currency    *string       `json:"currency"`
	Description *string       `json:"description"`
	Items       []itemRequest `json:"items"`
}

type itemRequest struct {
	Description *string `json:"description"`
	Quantity    *string `json:"quantity"`
	UnitPrice   *string `json:"unit_price"`
	LineTotal   *string `json:"line_total"`
}

type expenseResponse struct {
	ID          string         `json:"id"`
	OccurredAt  string         `json:"occurred_at"`
	Amount      string         `json:"amount"`
	Currency    string         `json:"currency"`
	Description string         `json:"description"`
	Items       []itemResponse `json:"items"`
}

type itemResponse struct {
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	LineTotal   string `json:"line_total"`
}

// formatAmount keeps the scale: "45.90" must not become "45.9".
func formatAmount(d decimal.Decimal) string { return d.StringFixed(-d.Exponent()) }

// parseAmount turns the request text into a decimal. A missing amount stays nil
// for the domain to report; text that is not a number is reported here.
func parseAmount(s *string) (*decimal.Decimal, *domain.FieldError) {
	if s == nil {
		return nil, nil
	}
	d, err := decimal.NewFromString(*s)
	if err != nil {
		return nil, &domain.FieldError{Field: "amount", Code: "invalid", Detail: "must be a decimal number"}
	}
	return &d, nil
}

func parseItems(reqs []itemRequest) ([]domain.PurchaseItemInput, []domain.FieldError) {
	items := make([]domain.PurchaseItemInput, len(reqs))
	var errs []domain.FieldError
	for i, req := range reqs {
		description := ""
		if req.Description != nil {
			description = *req.Description
		}
		items[i].Description = description
		for _, field := range []struct {
			name  string
			value *string
			dst   **decimal.Decimal
		}{
			{name: "quantity", value: req.Quantity, dst: &items[i].Quantity},
			{name: "unit_price", value: req.UnitPrice, dst: &items[i].UnitPrice},
			{name: "line_total", value: req.LineTotal, dst: &items[i].LineTotal},
		} {
			value, ferr := parseAmount(field.value)
			if ferr != nil {
				ferr.Field = fmt.Sprintf("items[%d].%s", i, field.name)
				errs = append(errs, *ferr)
			}
			*field.dst = value
		}
	}
	return items, errs
}

// withParseErrors replaces required errors for unparseable values and adds any
// adapter errors that do not have a matching domain validation error.
func withParseErrors(fields, parseErrs []domain.FieldError) []domain.FieldError {
	if len(parseErrs) == 0 {
		return fields
	}
	out := append([]domain.FieldError(nil), fields...)
	fieldIndexes := make(map[string]int, len(out)+len(parseErrs))
	for i, field := range out {
		if _, exists := fieldIndexes[field.Field]; !exists {
			fieldIndexes[field.Field] = i
		}
	}
	for _, parseErr := range parseErrs {
		if i, exists := fieldIndexes[parseErr.Field]; exists {
			out[i] = parseErr
			continue
		}
		fieldIndexes[parseErr.Field] = len(out)
		out = append(out, parseErr)
	}
	return out
}

func toResponse(e *domain.Expense) expenseResponse {
	response := expenseResponse{
		ID:          e.ID.String(),
		OccurredAt:  e.LocalTime().Format(time.RFC3339),
		Amount:      formatAmount(e.Money.Amount),
		Currency:    string(e.Money.Currency),
		Description: e.Description,
		Items:       make([]itemResponse, 0, len(e.Items)),
	}
	for _, item := range e.Items {
		response.Items = append(response.Items, itemResponse{
			Description: item.Description,
			Quantity:    formatAmount(item.Quantity),
			UnitPrice:   formatAmount(item.UnitPrice),
			LineTotal:   formatAmount(item.LineTotal),
		})
	}
	return response
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeProblem(w, http.StatusBadRequest, "malformed-json", "Malformed JSON", nil)
		return
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		writeProblem(w, http.StatusBadRequest, "malformed-json", "Malformed JSON", nil)
		return
	}

	// occurred_at needs adapter-side decoding (RFC 3339 text), so its presence
	// and format are checked here: a missing key or malformed text is reported
	// immediately (there is no placeholder that could flow through without
	// risking a false-positive validation pass, and Register persists on
	// success). amount is parsed here too, but its parse error is merged with
	// the domain's. currency/description pass through as-is — Money owns the
	// pair's presence and content checks.
	if req.OccurredAt == nil {
		writeProblem(w, http.StatusUnprocessableEntity, "validation-failed", "Validation failed",
			[]domain.FieldError{{Field: "occurred_at", Code: "required", Detail: "occurred_at is required"}})
		return
	}
	occurredAt, err := time.Parse(time.RFC3339, *req.OccurredAt)
	if err != nil {
		writeProblem(w, http.StatusUnprocessableEntity, "validation-failed", "Validation failed",
			[]domain.FieldError{{Field: "occurred_at", Code: "invalid", Detail: "must be RFC 3339 with explicit offset"}})
		return
	}

	amount, amountErr := parseAmount(req.Amount)
	items, itemParseErrs := parseItems(req.Items)
	e, err := h.svc.Register(r.Context(), application.RegisterInput{
		OccurredAt:  &occurredAt,
		Amount:      amount,
		Currency:    req.Currency,
		Description: req.Description,
		Items:       items,
	})
	var verr *domain.ValidationError
	if errors.As(err, &verr) {
		parseErrs := itemParseErrs
		if amountErr != nil {
			parseErrs = append(parseErrs, *amountErr)
		}
		writeProblem(w, http.StatusUnprocessableEntity, "validation-failed", "Validation failed", withParseErrors(verr.Fields, parseErrs))
		return
	}
	if err != nil {
		log.Printf("register expense: %v", err)
		writeProblem(w, http.StatusInternalServerError, "internal", "Internal error", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toResponse(e))
}

func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeProblem(w, http.StatusBadRequest, "invalid-id", "Invalid expense id", nil)
		return
	}
	e, err := h.svc.Find(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeProblem(w, http.StatusNotFound, "not-found", "Expense not found", nil)
		return
	}
	if err != nil {
		log.Printf("get expense %s: %v", id, err)
		writeProblem(w, http.StatusInternalServerError, "internal", "Internal error", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toResponse(e))
}
