package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	var listenPort int
	listenPortFromUser := os.Getenv("PORT")
	if listenPortFromUser == "" {
		listenPort = 8080
	} else {
		listenPort, err = strconv.Atoi(listenPortFromUser)
		if err != nil {
			log.Fatalf("Invalid PORT, expected an integer: %s", err.Error())
		}
	}

	bigPictureIgnore := map[int]struct{}{}
	bigPictureIgnoreString := os.Getenv("BIG_PICTURE_IGNORE")
	if bigPictureIgnoreString != "" {
		bigPictureIgnoreSlice := strings.Split(bigPictureIgnoreString, ",")
		for i := range bigPictureIgnoreSlice {
			categoryID, err := strconv.Atoi(bigPictureIgnoreSlice[i])
			if err != nil {
				log.Fatalf("converting '%s' to integer for category ID: %s", bigPictureIgnoreSlice[i], err)
			}
			bigPictureIgnore[categoryID] = struct{}{}
		}
	}

	bigPictureIncome := map[int]struct{}{}
	bigPictureIncomeString := os.Getenv("BIG_PICTURE_INCOME")
	if bigPictureIncomeString != "" {
		bigPictureIncomeSlice := strings.Split(bigPictureIncomeString, ",")
		for i := range bigPictureIncomeSlice {
			categoryID, err := strconv.Atoi(bigPictureIncomeSlice[i])
			if err != nil {
				log.Fatalf("converting '%s' to integer for category ID: %s", bigPictureIncomeSlice[i], err)
			}
			bigPictureIncome[categoryID] = struct{}{}
		}
	}

	autocompleteIgnoredCategories := map[int]struct{}{}
	autocompleteIgnoredCategoriesString := os.Getenv("AUTOCOMPLETE_CATEGORIES_IGNORE")
	if autocompleteIgnoredCategoriesString != "" {
		autocompleteIgnoredCategoriesSlice := strings.Split(autocompleteIgnoredCategoriesString, ",")
		for i := range autocompleteIgnoredCategoriesSlice {
			categoryID, err := strconv.Atoi(autocompleteIgnoredCategoriesSlice[i])
			if err != nil {
				log.Fatalf("converting '%s' to integer for category ID: %s", autocompleteIgnoredCategoriesSlice[i], err)
			}
			autocompleteIgnoredCategories[categoryID] = struct{}{}
		}
	}

	f, err := firefly.New(
		&http.Client{Timeout: time.Second * 30},
		firefly.Config{
			Token: token,
			URL:   fireflyBase,

			BigPictureIgnore:              bigPictureIgnore,
			BigPictureIncome:              bigPictureIncome,
			AutocompleteIgnoredCategories: autocompleteIgnoredCategories,
		},
	)
	if err != nil {
		log.Fatalf("Could not initialize Firefly-III client: %s", err)
	}

	http.HandleFunc("/api/transactions/", f.HandleTxn)
	http.HandleFunc("/api/accounts/", f.HandleAccount)
	http.HandleFunc("/api/categories/", f.HandleCategory)
	http.HandleFunc("/api/bigpicture/", f.HandleBigPicture)

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
				log.Printf("Failed to check for stale cache: %s", err)
			}
			err = f.InvalidateCacheIfCategoriesHaveChanged(c, b)
			if err != nil {
				log.Printf("Failed to check for stale categories: %s", err)
			}
		}
	}(c, b)

	log.Println("Listening for connections...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}
