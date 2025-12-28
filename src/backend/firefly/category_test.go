package firefly_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/shopspring/decimal"
)

func TestCategories(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/categories/", nil)
	f.HandleCategory(w, req)

	var c []firefly.Category
	json.NewDecoder(w.Body).Decode(&c)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}

	if len(c) != 1 {
		t.Fatalf("Got %d Categories, wanted 1", len(c))
	}
	if c[0].Name != "Apartment" {
		t.Fatalf("Got %s as Category name, wanted Apartment", c[0].Name)
	}
	if c[0].ID != 4 {
		t.Fatalf("Got %d as Category ID, wanted 4", c[0].ID)
	}
	if c[0].AutocompleteIgnore != false {
		t.Fatalf("Expected AutocompleteIgnore to be false")
	}
}

func TestListCategoryTotals(t *testing.T) {
	// (Interval not considered in test)
	start := time.Now().Add(time.Hour * -1)
	end := time.Now()
	c, err := f.CachedListCategoryTotals(start, end)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if len(c) != 1 {
		t.Fatalf("Got %d CategoryTotals, wanted 1", len(c))
	}
	if c[0].Name != "Apartment" {
		t.Fatalf("Got %s as Category name, wanted Apartment", c[0].Name)
	}
	if c[0].ID != 4 {
		t.Fatalf("Got %d as Category ID, wanted 4", c[0].ID)
	}
	spent, _ := decimal.NewFromString("-237.80")
	if !c[0].Spent.Equal(spent) {
		t.Fatalf("Got %d as Spent, wanted -237.80", c[0].Spent)
	}
	earned, _ := decimal.NewFromString("54.23")
	if !c[0].Earned.Equal(earned) {
		t.Fatalf("Got %d as Earned, wanted 54.23", c[0].Earned)
	}
	if !c[0].Start.Equal(start) {
		t.Fatalf("Got %s as Start, wanted %s", c[0].Start, start)
	}
	if !c[0].End.Equal(end) {
		t.Fatalf("Got %s as End, wanted %s", c[0].End, end)
	}
}

func TestFetchCategoryTotal(t *testing.T) {
	// (Interval not considered in test)
	start := time.Now().Add(time.Hour * -1)
	end := time.Now()
	c, err := f.CachedFetchCategoryTotals(4, start, end)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if len(c) != 1 {
		t.Fatalf("Got %d CategoryTotals, wanted 1", len(c))
	}
	if c[0].Name != "Apartment" {
		t.Fatalf("Got %s as Category name, wanted Apartment", c[0].Name)
	}
	if c[0].ID != 4 {
		t.Fatalf("Got %d as Category ID, wanted 4", c[0].ID)
	}
	spent, _ := decimal.NewFromString("-323.75")
	if !c[0].Spent.Equal(spent) {
		t.Fatalf("Got %d as Spent, wanted -323.75", c[0].Spent)
	}
	earned, _ := decimal.NewFromString("54.23")
	if !c[0].Earned.Equal(earned) {
		t.Fatalf("Got %d as Earned, wanted 54.23", c[0].Earned)
	}
	if !c[0].Start.Equal(start) {
		t.Fatalf("Got %s as Start, wanted %s", c[0].Start, start)
	}
	if !c[0].End.Equal(end) {
		t.Fatalf("Got %s as End, wanted %s", c[0].End, end)
	}
}
