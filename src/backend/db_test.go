package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSetupDB(t *testing.T) {
	var (
		mock sqlmock.Sqlmock
		err  error
	)

	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error opening mock database connection: %s\n", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS budgets.*`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS category_budgets.*`).WillReturnResult(sqlmock.NewResult(1, 1))

	setupDB(db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
