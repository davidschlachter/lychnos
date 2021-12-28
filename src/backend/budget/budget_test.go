package budget_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davidschlachter/lychnos/src/backend/budget"
)

func TestHandle(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error opening mock database connection: %s\n", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO budgets`).
		WithArgs(0, "2021-01-01 00:00:00", "2021-12-31 23:59:59", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT id, start, end, reporting_interval FROM budgets;`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "start", "end", "reporting_interval"}).
			AddRow(1, "2021-01-01 00:00:00", "2021-12-31 23:59:59", 0))
	mock.ExpectQuery(`SELECT id, start, end, reporting_interval FROM budgets WHERE id`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "start", "end", "reporting_interval"}).
			AddRow(1, "2021-01-01 00:00:00", "2021-12-31 23:59:59", 0))
	mock.ExpectExec(`INSERT INTO budgets`).
		WithArgs(1, "2022-01-01 00:00:00", "2022-12-31 23:59:59", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM budgets WHERE id`).WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	b := budget.New(db)

	// Create
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/budgets/", strings.NewReader("start=2021-01-01%2000%3A00%3A00&end=2021-12-31%2023%3A59%3A59"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusCreated)
	}

	// Find
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/budgets/", nil)
	b.Handle(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}
	var results []budget.Budget
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Fatalf("No budgets found, want 1")
	}

	// Fetch
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/budgets/1", nil)
	b.Handle(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Fatalf("No budgets found, want 1")
	}

	// Update
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/budgets/", strings.NewReader("id=1&start=2022-01-01%2000%3A00%3A00&end=2022-12-31%2023%3A59%3A59"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b.Handle(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusCreated)
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/budgets/1", nil)
	b.Handle(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusNoContent)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
