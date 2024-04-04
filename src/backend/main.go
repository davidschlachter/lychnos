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
	"github.com/davidschlachter/lychnos/src/backend/report"
)

func main() {
	db, err := connect()
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	defer db.Close()

	token := os.Getenv("FIREFLY_TOKEN")
	if token == "" {
		log.Fatal("Got empty FIREFLY_TOKEN, expected a value to be set.")
	}
	fireflyBase := os.Getenv("FIREFLY_URL")
	if fireflyBase == "" {
		log.Fatal("Got empty FIREFLY_URL, expected a value to be set.")
	}

	f, err := firefly.New(
		&http.Client{Timeout: time.Second * 30},
		token,
		fireflyBase,
	)
	if err != nil {
		log.Fatalf("Could not initialize Firefly-III client: %s", err)
	}

	http.HandleFunc("/api/transactions/", f.HandleTxn)
	http.HandleFunc("/api/accounts/", f.HandleAccount)
	http.HandleFunc("/api/categories/", f.HandleCategory)

	b := budget.New(db)
	http.HandleFunc("/api/budgets/", b.Handle)

	c := categorybudget.New(db, b)
	http.HandleFunc("/api/categorybudgets/", c.Handle)

	r, err := report.New(f, c, b)
	if err != nil {
		fmt.Printf("Could not initialize reports: %s\n", err)
		os.Exit(1)
	}
	http.HandleFunc("/api/reports/", r.Handle)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.RequestURI)
		fmt.Fprintf(w, "ok\n")
	})

	err = f.RefreshCaches(c, b)
	if err != nil {
		log.Fatalf("Failed to update caches: %s", err)
	}

	go func(c *categorybudget.CategoryBudgets, b *budget.Budgets) {
		for range time.Tick(time.Minute) {
			err = f.InvalidateCacheIfAccountBalancesHaveChanged(c, b)
			if err != nil {
				fmt.Printf("Failed to check for stale cache: %s", err)
			}
		}
	}(c, b)

	log.Println("Listening for connections...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
