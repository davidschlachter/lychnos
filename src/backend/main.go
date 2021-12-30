package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/cache"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/davidschlachter/lychnos/src/backend/report"
)

func main() {
	db := connect()
	setupDB(db)

	f, err := firefly.New(&http.Client{Timeout: time.Second * 5}, os.Getenv("FIREFLY_TOKEN"), os.Getenv("FIREFLY_URL"))
	if err != nil {
		fmt.Printf("Could not initialize Firefly-III client: %s\n", err)
		os.Exit(1)
	}

	b := budget.New(db)
	http.HandleFunc("/api/budgets/", b.Handle)

	c := categorybudget.New(db)
	http.HandleFunc("/api/categorybudgets/", c.Handle)

	h := cache.New(f)

	r, err := report.New(f, c, b, h)
	if err != nil {
		fmt.Printf("Could not initialize reports: %s\n", err)
		os.Exit(1)
	}
	http.HandleFunc("/api/reports/", r.Handle)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})

	h.RefreshCaches(c, b)

	log.Println("Listening for connections...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
