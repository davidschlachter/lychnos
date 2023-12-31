package categorybudget_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
)

func TestHandle(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error opening mock database connection: %s\n", err)
	}
	defer db.Close()

	// Insert
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets;`).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(`INSERT INTO category_budgets`).
		WithArgs(1, 1, "1000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// Find
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets;`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "budget", "category", "amount"}).
			AddRow(1, 1, 1, "1000"))
	// Fetch
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets WHERE id = \?;`).
		WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "budget", "category", "amount"}).
		AddRow(1, 1, 1, "1000"))
	// Delete
	mock.ExpectExec(`DELETE FROM category_budgets WHERE id`).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Replace
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets;`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "budget", "category", "amount"}).
			AddRow(1, 1, 1, "1000"))
	mock.ExpectExec(`DELETE FROM category_budgets WHERE id`).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO category_budgets`).
		WithArgs(1, 1, "25").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c := categorybudget.New(db, budget.New(db))

	// Create
	w := httptest.NewRecorder()
	amount, _ := decimal.NewFromString("1000.00")
	body, _ := json.Marshal([]categorybudget.CategoryBudget{{Budget: 1, Category: 1, Amount: amount}})
	req := httptest.NewRequest(http.MethodPost, "/api/categorybudgets/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusCreated, body)
	}

	// Find
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/categorybudgets/?budget=1", nil)
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusOK, body)
	}
	var results []categorybudget.CategoryBudget
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Fatalf("No categorybudgets found, want 1")
	}

	// Fetch
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/categorybudgets/1", nil)
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusOK, body)
	}
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Fatalf("No categorybudgets found, want 1")
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/categorybudgets/1", nil)
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusCreated, body)
	}

	// Create multiple
	w = httptest.NewRecorder()
	amount, _ = decimal.NewFromString("25.00")
	body, _ = json.Marshal([]categorybudget.CategoryBudget{{Budget: 1, Category: 1, Amount: amount}})
	req = httptest.NewRequest(http.MethodPost, "/api/categorybudgets/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(w.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusCreated, body)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
