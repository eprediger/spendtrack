package httpapi

import (
	"encoding/json"
	"net/http"

	"spendtrack/internal/domain"
)

const problemType = "https://spendtrack.dev/problems/"

type problem struct {
	Type   string              `json:"type"`
	Title  string              `json:"title"`
	Status int                 `json:"status"`
	Errors []domain.FieldError `json:"errors,omitempty"`
}

func writeProblem(w http.ResponseWriter, status int, kind, title string, errs []domain.FieldError) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem{
		Type:   problemType + kind,
		Title:  title,
		Status: status,
		Errors: errs,
	})
}
