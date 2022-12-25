package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type connectionType int

const (
	connectionTypeSQLite connectionType = iota
	connectionTypeMySQL
)

var connectionTypeStrings = [...]string{
	connectionTypeSQLite: "sqlite3",
	connectionTypeMySQL:  "mysql",
}

func (c connectionType) String() string {
	if int(c) < 0 || int(c) >= len(connectionTypeStrings) {
		return ""
	}
	return connectionTypeStrings[int(c)]
}

type connectionString struct {
	DSN  string
	Type connectionType
}

// connect connects to the database, based on the connection string provided by
// the user. If no connection string was provided, attempt to create a SQLite
// database.
func connect() (*sql.DB, error) {
	var db *sql.DB

	dsn, err := getDSN()
	if err != nil {
		return nil, err
	}

	if dsn.DSN == "" {
		db, err = createSQLite()
	} else {
		db, err = sql.Open(dsn.Type.String(), dsn.DSN)
	}
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Validate database connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %s", err)
	}

	// Set up tables
	err = setupDB(&dsn, db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %s", err)
	}

	return db, nil
}

func createSQLite() (*sql.DB, error) {
	// create the DB
	db, err := sql.Open("sqlite3", "./database.sqlite3")
	return db, err
}

// getDSN reads the database connection string from the environment, returning
// it along with its type.
func getDSN() (connectionString, error) {
	var (
		dsn                 connectionString
		mysqlDSN, sqliteDSN string
	)
	mysqlDSN = os.Getenv("MYSQL_DSN")
	sqliteDSN = os.Getenv("SQLITE_DSN")

	if mysqlDSN != "" && sqliteDSN != "" {
		return dsn, fmt.Errorf("MYSQL_DSN and SQLITE_DSN cannot both be set")
	}

	if mysqlDSN != "" {
		dsn.DSN = mysqlDSN
		dsn.Type = connectionTypeMySQL
	} else {
		dsn.DSN = sqliteDSN // okay if empty, will create SQLite database
		dsn.Type = connectionTypeSQLite
	}

	return dsn, nil
}

// setupDB will create the tables that lychnos uses to store its budget data, if
// they do not already exist.
func setupDB(dsn *connectionString, db *sql.DB) error {
	var q []string

	if dsn != nil && dsn.Type == connectionTypeMySQL {
		// MySQL
		q = []string{`
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
	} else {
		// SQLite
		q = []string{`
CREATE TABLE IF NOT EXISTS budgets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	start DATETIME NOT NULL,
	end DATETIME NOT NULL,
	reporting_interval INT NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS category_budgets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	budget INT,
	category INT,
	amount DECIMAL(12,4),
	FOREIGN KEY ( budget ) REFERENCES budgets( id )
);
		`}
	}

	for _, s := range q {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}
