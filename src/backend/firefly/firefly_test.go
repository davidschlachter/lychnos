package firefly_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/shopspring/decimal"
)

var server *httptest.Server
var f *firefly.Firefly

func TestMain(m *testing.M) {
	var err error
	setup()
	defer server.Close()
	f, err = firefly.New(server.Client(), "token", server.URL)
	if err != nil {
		fmt.Printf("Unexpected error in setup: %s\n", err)
		os.Exit(1)
	}
	status := m.Run()
	os.Exit(status)
}

func setup() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/categories/4", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"type":"categories","id":"4","attributes":{"created_at":"2019-09-07T20:02:33-04:00","updated_at":"2019-09-07T20:02:33-04:00","name":"Apartment","notes":null,"spent":[{"sum":"-323.75","currency_id":9,"currency_name":"Canadian dollar","currency_symbol":"C$","currency_code":"CAD","currency_decimal_places":2}],"earned":[{"sum":"54.23","currency_id":9,"currency_name":"Canadian dollar","currency_symbol":"C$","currency_code":"CAD","currency_decimal_places":2}]},"links":{"0":{"rel":"self","uri":"/categories/4"},"self":"http://192.168.6.4:8753/api/v1/categories/4"}}}`))
	})
	mux.HandleFunc("/api/v1/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"type":"categories","id":"4","attributes":{"created_at":"2019-09-07T20:02:33-04:00","updated_at":"2019-09-07T20:02:33-04:00","name":"Apartment","notes":null,"spent":[{"sum":"-237.80","currency_id":9,"currency_name":"Canadian dollar","currency_symbol":"C$","currency_code":"CAD","currency_decimal_places":2}],"earned":[{"sum":"54.23","currency_id":9,"currency_name":"Canadian dollar","currency_symbol":"C$","currency_code":"CAD","currency_decimal_places":2}]},"links":{"0":{"rel":"self","uri":"/categories/4"},"self":"http://lychnos/api/v1/categories/4"}}]}`))
	})
	mux.HandleFunc("/api/v1/autocomplete/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"id":"4","name":"Apartment"}]`))
	})
	server = httptest.NewServer(mux)
}

func TestCategories(t *testing.T) {
	c, err := f.Categories()
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
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
}

func TestListCategoryTotals(t *testing.T) {
	// (Interval not considered in test)
	start := time.Now().Add(time.Hour * -1)
	end := time.Now()
	c, err := f.ListCategoryTotals(start, end)
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
	c, err := f.FetchCategoryTotal(4, start, end)
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
