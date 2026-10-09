package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

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
	OccurredAt  *string `json:"occurred_at"`
	Amount      *string `json:"amount"`
	Currency    *string `json:"currency"`
	Description *string `json:"description"`
}

type expenseResponse struct {
	ID          string `json:"id"`
	OccurredAt  string `json:"occurred_at"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

func toResponse(e *domain.Expense) expenseResponse {
	return expenseResponse{
		ID:          e.ID.String(),
		OccurredAt:  e.LocalTime().Format(time.RFC3339),
		Amount:      e.Money.String(),
		Currency:    string(e.Money.Currency),
		Description: e.Description,
	}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, http.StatusBadRequest, "malformed-json", "Malformed JSON", nil)
		return
	}

	// occurred_at needs adapter-side decoding (RFC 3339 text), so its presence
	// and format are checked here: a missing key or malformed text is reported
	// immediately (there is no placeholder that could flow through without
	// risking a false-positive validation pass, and Register persists on
	// success). amount/currency/description pass through as-is — Money owns
	// its pair's presence and content checks.
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

	e, err := h.svc.Register(r.Context(), application.RegisterInput{
		OccurredAt:  &occurredAt,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
	})
	var verr *domain.ValidationError
	if errors.As(err, &verr) {
		writeProblem(w, http.StatusUnprocessableEntity, "validation-failed", "Validation failed", verr.Fields)
		return
	}
	if err != nil {
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
		writeProblem(w, http.StatusInternalServerError, "internal", "Internal error", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toResponse(e))
}
