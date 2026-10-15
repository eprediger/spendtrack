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
		_, err := pool.Exec(ctx, "TRUNCATE expense_items, expenses")
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
	sc.Step(`^an expense with the unsupported currency "([^"]*)" is stored in the database$`, storedWithCurrency)
	sc.Step(`^the system indicates an internal error$`, internalError)
	sc.Step(`^the expense occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the expense shows occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the Financial Manager registers an expense at "([^"]*)" for (\S+) "([^"]*)" described as "([^"]*)" with items:$`, registerExpenseWithItems)
	sc.Step(`^the expense has (\d+) items$`, expenseHasItems)
	sc.Step(`^the expense has no items$`, expenseHasNoItems)
	sc.Step(`^an expense with purchase items was registered$`, registeredExpenseWithItems)
	sc.Step(`^the expense shows each purchase item with its description, quantity, unit price and line total$`, expenseShowsPurchaseItems)
	sc.Step(`^an expense with a purchase item with an empty description is stored in the database$`, storedWithEmptyPurchaseItemDescription)
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

func registerExpenseWithItems(occurredAt, amount, currency, description string, table *godog.Table) error {
	headers := make([]string, len(table.Rows[0].Cells))
	for i, cell := range table.Rows[0].Cells {
		headers[i] = strings.ReplaceAll(cell.Value, " ", "_")
	}
	items := make([]map[string]string, 0, len(table.Rows)-1)
	for _, row := range table.Rows[1:] {
		item := make(map[string]string, len(headers))
		for i, cell := range row.Cells {
			item[headers[i]] = cell.Value
		}
		items = append(items, item)
	}
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": occurredAt, "amount": amount,
		"currency": currency, "description": description, "items": items,
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

func expenseHasItems(expected int) error {
	items, ok := st.body["items"].([]any)
	if !ok || len(items) != expected {
		return fmt.Errorf("expected %d items, got %v", expected, st.body["items"])
	}
	return nil
}

func expenseHasNoItems() error {
	return expenseHasItems(0)
}

func registeredExpenseWithItems() error {
	if err := post("/api/v1/expenses", map[string]any{
		"occurred_at": "2026-10-05T16:53:34-03:00",
		"amount":      "14000.00",
		"currency":    "ARS",
		"description": "Gas station store",
		"items": []map[string]string{
			{"description": "Candies", "quantity": "5", "unit_price": "150.00", "line_total": "750.00"},
			{"description": "Bottled water", "quantity": "1", "unit_price": "3400.00", "line_total": "3400.00"},
		},
	}); err != nil {
		return err
	}
	return recordedExpense()
}

func expenseShowsPurchaseItems() error {
	if st.status != http.StatusOK {
		return fmt.Errorf("expected 200, got %d: %v", st.status, st.body)
	}
	items, ok := st.body["items"].([]any)
	if !ok || len(items) != 2 {
		return fmt.Errorf("expected two purchase items, got %v", st.body["items"])
	}
	expected := []map[string]string{
		{"description": "Candies", "quantity": "5", "unit_price": "150.00", "line_total": "750.00"},
		{"description": "Bottled water", "quantity": "1", "unit_price": "3400.00", "line_total": "3400.00"},
	}
	for i, item := range items {
		got, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("item %d has unexpected shape: %v", i, item)
		}
		for field, want := range expected[i] {
			if got[field] != want {
				return fmt.Errorf("item %d %s: want %q got %v", i, field, want, got[field])
			}
		}
	}
	return nil
}

func storedWithEmptyPurchaseItemDescription() error {
	err := pool.QueryRow(context.Background(),
		`INSERT INTO expenses (id, occurred_at, occurred_offset_minutes, amount, currency, description)
		 VALUES (gen_random_uuid(), now(), 0, 1, 'ARS', 'seeded') RETURNING id::text`).
		Scan(&st.expenseID)
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(),
		`INSERT INTO expense_items (expense_id, position, description, quantity, unit_price, line_total)
		 VALUES ($1, 0, '', 1, 1, 1)`, st.expenseID)
	return err
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

// storedWithCurrency inserts a row directly, bypassing the API validation,
// to simulate data that did not come through registration (a seed or script).
func storedWithCurrency(currency string) error {
	return pool.QueryRow(context.Background(),
		`INSERT INTO expenses (id, occurred_at, occurred_offset_minutes, amount, currency, description)
		 VALUES (gen_random_uuid(), now(), 0, 10, $1, 'seeded') RETURNING id::text`, currency).
		Scan(&st.expenseID)
}

func internalError() error {
	if st.status != http.StatusInternalServerError || st.body["title"] != "Internal error" {
		return fmt.Errorf("expected 500 Internal error, got %d: %v", st.status, st.body)
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
	if got := st.body["amount"]; got != amount {
		return fmt.Errorf("amount: want %q got %v", amount, got)
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
