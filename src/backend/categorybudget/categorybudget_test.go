package categorybudget_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
)

func TestHandle(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error opening mock database connection: %s\n", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO category_budgets`).
		WithArgs(0, 1, 1, "1000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets;`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "budget", "category", "amount"}).
			AddRow(1, 1, 1, "1000"))
	mock.ExpectQuery(`SELECT id, budget, category, amount FROM category_budgets WHERE id`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "budget", "category", "amount"}).
			AddRow(1, 1, 1, "1000"))
	mock.ExpectExec(`INSERT INTO category_budgets`).
		WithArgs(1, 2, 1, "1000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM category_budgets WHERE id`).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	c := categorybudget.New(db)

	// Create
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/categorybudgets/", strings.NewReader("budget=1&category=1&amount=1000"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusCreated)
	}

	// Find
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/categorybudgets/", nil)
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
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
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Fatalf("No categorybudgets found, want 1")
	}

	// Update
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/categorybudgets/", strings.NewReader("id=1&budget=2&category=1&amount=1000"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusCreated)
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/categorybudgets/1", nil)
	c.Handle(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusNoContent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
