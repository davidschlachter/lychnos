package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/davidschlachter/lychnos/src/backend/budget"
)

var db *sql.DB

func main() {
	db = connect()
	setupDB(db)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})

	b := budget.New(db)
	http.HandleFunc("/api/budget", b.Handle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
