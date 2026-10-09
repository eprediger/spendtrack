package features_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/shopspring/decimal"
)

type state struct {
	status    int
	body      map[string]any
	expenseID string
}

var st *state

func initializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		st = &state{}
		_, err := pool.Exec(ctx, "TRUNCATE expenses")
		return ctx, err
	})

	// health
	sc.Step(`^the app is running$`, func() error { return nil })
	sc.Step(`^I request the health status$`, requestHealth)
	sc.Step(`^I should receive a healthy status$`, func() error {
		if st.status != http.StatusOK {
			return fmt.Errorf("status %d", st.status)
		}
		return nil
	})

	// expenses
	sc.Step(`^a Financial Context$`, func() error { return nil })
	sc.Step(`^the Financial Manager registers an expense at "([^"]*)" for (\S+) "([^"]*)" described as "([^"]*)"$`, registerExpense)
	sc.Step(`^an expense at "([^"]*)" for (\S+) "([^"]*)" described as "([^"]*)" was registered$`, registeredExpense)
	sc.Step(`^the expense is recorded in the Financial Context$`, recordedExpense)
	sc.Step(`^the Financial Manager reviews the expense$`, reviewExpense)
	sc.Step(`^the Financial Manager reviews an expense with an unknown id$`, reviewUnknownExpense)
	sc.Step(`^the Financial Manager reviews an expense with a malformed id$`, reviewMalformedExpense)
	sc.Step(`^the system indicates the expense was not found$`, notFound)
	sc.Step(`^the system indicates the id is invalid$`, invalidID)
	sc.Step(`^the expense occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the expense shows occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the Financial Manager registers an expense without "([^"]*)"$`, registerWithout)
	sc.Step(`^the Financial Manager registers an expense with:$`, registerWithTable)
	sc.Step(`^the Financial Manager registers an expense that occurs in the future for (\S+) "([^"]*)" described as "([^"]*)"$`, registerFuture)
	sc.Step(`^the registration is rejected$`, rejected)
	sc.Step(`^the system indicates "([^"]*)" is required$`, indicatesRequired)
	sc.Step(`^the system indicates the (\w+) is invalid$`, indicatesInvalid)
	sc.Step(`^the system indicates the occurrence cannot be in the future$`, indicatesFuture)
}

func post(path string, payload map[string]any) error {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	st.status = resp.StatusCode
	raw, _ := io.ReadAll(resp.Body)
	st.body = map[string]any{}
	_ = json.Unmarshal(raw, &st.body)
	return nil
}

func registerExpense(occurredAt, amount, currency, description string) error {
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": occurredAt, "amount": amount,
		"currency": currency, "description": description,
	})
}

func registeredExpense(occurredAt, amount, currency, description string) error {
	if err := registerExpense(occurredAt, amount, currency, description); err != nil {
		return err
	}
	return recordedExpense()
}

func recordedExpense() error {
	if st.status != http.StatusCreated {
		return fmt.Errorf("expected 201, got %d: %v", st.status, st.body)
	}
	id, _ := st.body["id"].(string)
	if id == "" {
		return fmt.Errorf("no id in response: %v", st.body)
	}
	st.expenseID = id
	return nil
}

func reviewExpense() error {
	return getExpense(st.expenseID)
}

// reviewUnknownExpense uses a well-formed id absent from the (per-scenario
// truncated) table.
func reviewUnknownExpense() error {
	return getExpense("00000000-0000-0000-0000-000000000000")
}

func reviewMalformedExpense() error {
	return getExpense("not-a-uuid")
}

func getExpense(id string) error {
	resp, err := http.Get(apiURL + "/api/v1/expenses/" + id)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	st.status = resp.StatusCode
	raw, _ := io.ReadAll(resp.Body)
	st.body = map[string]any{}
	return json.Unmarshal(raw, &st.body)
}

func notFound() error {
	if st.status != http.StatusNotFound {
		return fmt.Errorf("expected 404, got %d: %v", st.status, st.body)
	}
	return nil
}

func invalidID() error {
	if st.status != http.StatusBadRequest {
		return fmt.Errorf("expected 400, got %d: %v", st.status, st.body)
	}
	return nil
}

func requestHealth() error {
	resp, err := http.Get(apiURL + "/api/v1/health")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	st.status = resp.StatusCode
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}

func expenseMatches(occurredAt, amount, currency, description string) error {
	if st.status != http.StatusOK && st.status != http.StatusCreated {
		return fmt.Errorf("status %d", st.status)
	}
	if got := st.body["occurred_at"]; got != occurredAt {
		return fmt.Errorf("occurred_at: want %q got %v", occurredAt, got)
	}
	want, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}
	gotStr, _ := st.body["amount"].(string)
	got, err := decimal.NewFromString(gotStr)
	if err != nil || !got.Equal(want) {
		return fmt.Errorf("amount: want %s got %v", amount, st.body["amount"])
	}
	if got := st.body["currency"]; got != currency {
		return fmt.Errorf("currency: want %q got %v", currency, got)
	}
	if got := st.body["description"]; got != description {
		return fmt.Errorf("description: want %q got %v", description, got)
	}
	return nil
}

func registerWithout(field string) error {
	payload := map[string]any{
		"occurred_at": "2026-09-20T14:30:00-03:00", "amount": "45.90",
		"currency": "USD", "description": "Groceries",
	}
	delete(payload, strings.ReplaceAll(field, " ", "_"))
	return post("/api/v1/expenses", payload)
}

// registerWithTable builds a payload from only the given field/value rows,
// so a row can be omitted entirely to represent a missing field.
func registerWithTable(table *godog.Table) error {
	payload := map[string]any{}
	for _, row := range table.Rows {
		key := strings.ReplaceAll(row.Cells[0].Value, " ", "_")
		payload[key] = row.Cells[1].Value
	}
	return post("/api/v1/expenses", payload)
}

func registerFuture(amount, currency, description string) error {
	future := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": future, "amount": amount,
		"currency": currency, "description": description,
	})
}

func rejected() error {
	if st.status != http.StatusUnprocessableEntity {
		return fmt.Errorf("expected 422, got %d: %v", st.status, st.body)
	}
	return nil
}

func fieldErrors(field, code string) error {
	errs, _ := st.body["errors"].([]any)
	for _, e := range errs {
		m, _ := e.(map[string]any)
		if m["field"] == field && (code == "" || m["code"] == code) {
			return nil
		}
	}
	return fmt.Errorf("no error for %q (%s) in %v", field, code, st.body)
}

func indicatesRequired(field string) error {
	return fieldErrors(strings.ReplaceAll(field, " ", "_"), "required")
}

func indicatesInvalid(field string) error {
	return fieldErrors(field, "")
}

func indicatesFuture() error {
	return fieldErrors("occurred_at", "future")
}
