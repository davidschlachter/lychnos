package budget_test

import (
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
		WithArgs(nil, "2021-01-01 00:00:00", "2021-12-31 23:59:59", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT id, start, end, reporting_interval FROM budgets;`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "start", "end", "reporting_interval"}).
			AddRow(1, "2021-01-01 00:00:00", "2021-12-31 23:59:59", 0))

	b := budget.New(db)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/budget", strings.NewReader("start=2021-01-01%2000%3A00%3A00&end=2021-12-31%2023%3A59%3A59"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b.Handle(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusCreated)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/budget", nil)
	b.Handle(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
