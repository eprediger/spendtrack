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
		_, err := pool.Exec(ctx, "TRUNCATE expense_item_discounts, expense_items, expenses")
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
	sc.Step(`^the expense is recorded$`, recordedExpense)
	sc.Step(`^the expense is recorded in the Financial Context$`, recordedExpense)
	sc.Step(`^the Financial Manager reviews the expense$`, reviewExpense)
	sc.Step(`^the Financial Manager reviews an expense with an unknown id$`, reviewUnknownExpense)
	sc.Step(`^the Financial Manager reviews an expense with a malformed id$`, reviewMalformedExpense)
	sc.Step(`^the system indicates the expense was not found$`, notFound)
	sc.Step(`^the system indicates the id is invalid$`, invalidID)
	sc.Step(`^an expense with the unsupported currency "([^"]*)" is stored in the database$`, storedWithCurrency)
	sc.Step(`^an expense with an empty description is stored in the database$`, storedWithEmptyDescription)
	sc.Step(`^the system indicates an internal error$`, internalError)
	sc.Step(`^the expense occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the expense shows occurred at "([^"]*)" with amount (\S+) "([^"]*)" and description "([^"]*)"$`, expenseMatches)
	sc.Step(`^the Financial Manager registers an expense at "([^"]*)" for (\S+) "([^"]*)" described as "([^"]*)" with items:$`, registerExpenseWithItems)
	sc.Step(`^the expense has (\d+) items$`, expenseHasItems)
	sc.Step(`^the expense has no items$`, expenseHasNoItems)
	sc.Step(`^an expense with purchase items was registered$`, registeredExpenseWithItems)
	sc.Step(`^the expense shows each purchase item with its description, quantity, unit price and line total$`, expenseShowsPurchaseItems)
	sc.Step(`^an expense with a purchase item with an empty description is stored in the database$`, storedWithEmptyPurchaseItemDescription)
	sc.Step(`^the Financial Manager registers an expense with a purchase item of quantity (\S+), unit price (\S+) and line total (\S+)$`, registerPurchaseItem)
	sc.Step(`^the Financial Manager registers an expense with a purchase item "([^"]*)" of quantity (\S+), unit price (\S+) and line total (\S+)$`, registerNamedPurchaseItem)
	sc.Step(`^the Financial Manager registers an expense with a purchase item without "([^"]*)"$`, registerPurchaseItemWithout)
	sc.Step(`^the Financial Manager registers an expense with a purchase item whose "([^"]*)" is "([^"]*)"$`, registerPurchaseItemWithInvalid)
	sc.Step(`^the system indicates the purchase item "([^"]*)" is (required|invalid)$`, indicatesPurchaseItemError)
	sc.Step(`^the purchase item keeps the line total (\S+) as printed$`, purchaseItemKeepsLineTotal)
	sc.Step(`^the purchase item shows quantity (\S+) and line total (\S+)$`, purchaseItemShowsQuantityAndLineTotal)
	sc.Step(`^the Financial Manager registers an expense without "([^"]*)"$`, registerWithout)
	sc.Step(`^the Financial Manager registers an expense with:$`, registerWithTable)
	sc.Step(`^the Financial Manager registers an expense that occurs in the future for (\S+) "([^"]*)" described as "([^"]*)"$`, registerFuture)
	sc.Step(`^the Financial Manager registers an expense with a charged total of (\S+) "([^"]*)" and item "([^"]*)" line total of (\S+) "([^"]*)" without discounts$`, registerItemizedExpenseNoDiscounts)
	sc.Step(`^the Financial Manager registers an expense with charged total (\S+) "([^"]*)", item "([^"]*)" line total (\S+) "([^"]*)" with a store discount of (\S+) "([^"]*)" described as "([^"]*)", and other items totaling (\S+) "([^"]*)"$`, registerExpenseWithStoreDiscount)
	sc.Step(`^the Financial Manager registers an expense with a charged total of (\S+) "([^"]*)" and item "([^"]*)" line total of (\S+) "([^"]*)" with item discounts totaling (\S+) "([^"]*)"$`, registerItemDiscountsTotaling)
	sc.Step(`^the discount is recorded on the item "([^"]*)"$`, discountRecordedOnItem)
	sc.Step(`^an expense with a store discount was registered$`, registeredExpenseWithStoreDiscount)
	sc.Step(`^the expense shows the discount amount and description on the item it applies to$`, expenseShowsDiscount)
	sc.Step(`^the registration is rejected because the itemized expense does not reconcile$`, rejectedNotReconciled)
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
		"amount":      "4150.00",
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

func registerPurchaseItem(quantity, unitPrice, lineTotal string) error {
	return registerNamedPurchaseItem("Purchase item", quantity, unitPrice, lineTotal)
}

func registerNamedPurchaseItem(description, quantity, unitPrice, lineTotal string) error {
	return postExpenseWithPurchaseItem(map[string]string{
		"description": description,
		"quantity":    quantity,
		"unit_price":  unitPrice,
		"line_total":  lineTotal,
	})
}

func registerPurchaseItemWithout(field string) error {
	item := map[string]string{
		"description": "Purchase item",
		"quantity":    "1",
		"unit_price":  "150.00",
		"line_total":  "150.00",
	}
	delete(item, strings.ReplaceAll(field, " ", "_"))
	return postExpenseWithPurchaseItem(item)
}

func registerPurchaseItemWithInvalid(field, value string) error {
	item := map[string]string{
		"description": "Purchase item",
		"quantity":    "1",
		"unit_price":  "150.00",
		"line_total":  "150.00",
	}
	item[strings.ReplaceAll(field, " ", "_")] = value
	return postExpenseWithPurchaseItem(item)
}

func postExpenseWithPurchaseItem(item map[string]string) error {
	amount := "14000.00"
	if _, err := decimal.NewFromString(item["line_total"]); err == nil {
		amount = item["line_total"]
	}
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": "2026-10-05T16:53:34-03:00",
		"amount":      amount,
		"currency":    "ARS",
		"description": "Gas station store",
		"items":       []map[string]string{item},
	})
}

func indicatesPurchaseItemError(field, code string) error {
	field = strings.ReplaceAll(field, " ", "_")
	errs, _ := st.body["errors"].([]any)
	for _, e := range errs {
		m, _ := e.(map[string]any)
		gotField, _ := m["field"].(string)
		gotCode, _ := m["code"].(string)
		if strings.HasSuffix(gotField, "."+field) && (code == "invalid" || gotCode == code) {
			return nil
		}
	}
	return fmt.Errorf("no purchase item error for %q (%s) in %v", field, code, st.body)
}

func purchaseItemKeepsLineTotal(lineTotal string) error {
	items, ok := st.body["items"].([]any)
	if !ok || len(items) == 0 {
		return fmt.Errorf("expected a purchase item in response, got %v", st.body["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		return fmt.Errorf("purchase item has unexpected shape: %v", items[0])
	}
	if got := item["line_total"]; got != lineTotal {
		return fmt.Errorf("line_total: want %q got %v", lineTotal, got)
	}
	return nil
}

func purchaseItemShowsQuantityAndLineTotal(quantity, lineTotal string) error {
	items, ok := st.body["items"].([]any)
	if !ok || len(items) == 0 {
		return fmt.Errorf("expected a purchase item in response, got %v", st.body["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		return fmt.Errorf("purchase item has unexpected shape: %v", items[0])
	}
	if got := item["quantity"]; got != quantity {
		return fmt.Errorf("quantity: want %q got %v", quantity, got)
	}
	if got := item["line_total"]; got != lineTotal {
		return fmt.Errorf("line_total: want %q got %v", lineTotal, got)
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

func storedWithEmptyDescription() error {
	return pool.QueryRow(context.Background(),
		`INSERT INTO expenses (id, occurred_at, occurred_offset_minutes, amount, currency, description)
		 VALUES (gen_random_uuid(), now(), 0, 10, 'ARS', '') RETURNING id::text`).
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

func registerItemizedExpenseNoDiscounts(amount, currency, itemDescription, lineTotal, _ string) error {
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": "2026-10-05T16:53:34-03:00",
		"amount":      amount,
		"currency":    currency,
		"description": "Store purchase",
		"items": []map[string]string{
			{"description": itemDescription, "quantity": "1", "unit_price": lineTotal, "line_total": lineTotal},
		},
	})
}

func registerExpenseWithStoreDiscount(amount, currency, itemDescription, lineTotal, _ string, discountAmount, _ string, discountDescription, otherTotal, _ string) error {
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": "2026-10-05T16:53:34-03:00",
		"amount":      amount,
		"currency":    currency,
		"description": "Supermarket",
		"items": []map[string]any{
			{
				"description": itemDescription, "quantity": "1",
				"unit_price": lineTotal, "line_total": lineTotal,
				"discounts": []map[string]string{
					{"amount": discountAmount, "description": discountDescription},
				},
			},
			{"description": "Other items", "quantity": "1", "unit_price": otherTotal, "line_total": otherTotal},
		},
	})
}

func discountRecordedOnItem(description string) error {
	items, _ := st.body["items"].([]any)
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if item["description"] != description {
			continue
		}
		discounts, _ := item["discounts"].([]any)
		if len(discounts) != 1 {
			return fmt.Errorf("expected 1 discount on item %q, got %v", description, item["discounts"])
		}
		return nil
	}
	return fmt.Errorf("no item %q in %v", description, st.body["items"])
}

func registerItemDiscountsTotaling(amount, currency, itemDescription, lineTotal, _ string, discountTotal, _ string) error {
	return post("/api/v1/expenses", map[string]any{
		"occurred_at": "2026-10-05T16:53:34-03:00",
		"amount":      amount,
		"currency":    currency,
		"description": "Store purchase",
		"items": []map[string]any{
			{
				"description": itemDescription, "quantity": "1",
				"unit_price": lineTotal, "line_total": lineTotal,
				"discounts": []map[string]string{
					{"amount": discountTotal, "description": "Store discount"},
				},
			},
		},
	})
}

func registeredExpenseWithStoreDiscount() error {
	if err := registerExpenseWithStoreDiscount(
		"59690.00", "ARS", "COCA COLA R 2L", "12750.00", "ARS",
		"1700.00", "ARS", "2DO AL 40% LINEA COC", "48640.00", "ARS"); err != nil {
		return err
	}
	return recordedExpense()
}

func expenseShowsDiscount() error {
	if st.status != http.StatusOK {
		return fmt.Errorf("expected 200, got %d: %v", st.status, st.body)
	}
	items, _ := st.body["items"].([]any)
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if item["description"] != "COCA COLA R 2L" {
			continue
		}
		discounts, _ := item["discounts"].([]any)
		if len(discounts) != 1 {
			return fmt.Errorf("expected 1 discount on the item, got %v", item["discounts"])
		}
		discount, _ := discounts[0].(map[string]any)
		if discount["amount"] != "1700.00" || discount["description"] != "2DO AL 40% LINEA COC" {
			return fmt.Errorf("unexpected discount %v", discount)
		}
		return nil
	}
	return fmt.Errorf("no item %q in %v", "COCA COLA R 2L", st.body["items"])
}

func rejectedNotReconciled() error {
	if err := rejected(); err != nil {
		return err
	}
	return fieldErrors("items", "not_reconciled")
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
