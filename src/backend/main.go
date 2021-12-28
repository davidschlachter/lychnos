package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidschlachter/lychnos/src/backend/budget"
)

func main() {
	db := connect()
	setupDB(db)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})

	b := budget.New(db)
	http.HandleFunc("/api/budgets/", b.Handle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
