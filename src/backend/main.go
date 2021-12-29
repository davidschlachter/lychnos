package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/firefly"
)

func main() {
	db := connect()
	setupDB(db)

	_, err := firefly.New(&http.Client{Timeout: time.Second * 5}, os.Getenv("FIREFLY_TOKEN"), os.Getenv("FIREFLY_URL"))
	if err != nil {
		fmt.Printf("Could not initialize Firefly-III client: %s\n", err)
		os.Exit(1)
	}

	b := budget.New(db)
	http.HandleFunc("/api/budgets/", b.Handle)

	c := categorybudget.New(db)
	http.HandleFunc("/api/categorybudgets/", c.Handle)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
