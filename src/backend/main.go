package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	db = connect()
	setupDB(db)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	http.HandleFunc("/api/budget", handleBudget)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
