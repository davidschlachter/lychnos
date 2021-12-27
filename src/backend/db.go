package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Validate database connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Could not connect to database: %s", err)
		os.Exit(1)
	}

	return db
}

// setupDB will create the tables that lychnos uses to store its budget data, if
// they do not already exist.
func setupDB(db *sql.DB) {
	q := []string{`
CREATE TABLE IF NOT EXISTS budgets (
	id INT NOT NULL AUTO_INCREMENT,
	start DATETIME NOT NULL,
	end DATETIME NOT NULL,
	reporting_interval INT NOT NULL,
	PRIMARY KEY ( id )
);
`, `
CREATE TABLE IF NOT EXISTS category_budgets (
	id INT NOT NULL AUTO_INCREMENT,
	budget INT,
	category INT,
	amount DECIMAL(12,4),
	PRIMARY KEY ( id ),
	FOREIGN KEY ( budget ) REFERENCES budgets( id )
);
`}

	for _, s := range q {
		_, err := db.Exec(s)
		if err != nil {
			fmt.Printf("Could not initialize database tables: %s", err)
			os.Exit(1)
		}
	}
}
