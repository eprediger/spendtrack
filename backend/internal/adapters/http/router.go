package httpapi

import (
	"encoding/json"
	"net/http"

	"spendtrack/internal/application"
)

// NewRouter wires the API routes.
func NewRouter(svc *application.ExpenseService) http.Handler {
	h := NewExpenseHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/expenses", h.Create)
	mux.HandleFunc("GET /api/v1/expenses/{id}", h.Get)
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	return mux
}
