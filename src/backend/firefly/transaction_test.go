package firefly_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/shopspring/decimal"
)

func TestListTransactions(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/", nil)
	f.HandleTxn(w, req)

	var x []firefly.Transactions
	json.NewDecoder(w.Body).Decode(&x)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}
	if len(x) != 1 {
		t.Fatalf("Got %d transactions, wanted 1", len(x))
	}
	if x[0].ID != "2763" {
		t.Fatalf("Got transaction ID %s, wanted 2763", x[0].ID)
	}

	if len(x[0].Attributes.Transactions) != 1 {
		t.Fatalf("Got %d attributes.Transactions, wanted 1", len(x[0].Attributes.Transactions))
	}

	if x[0].Attributes.Transactions[0].Type != "deposit" {
		t.Fatalf("Got transaction Type %s, wanted deposit", x[0].Attributes.Transactions[0].Type)
	}
	if x[0].Attributes.Transactions[0].Date != "2022-01-01T00:00:00-05:00" {
		t.Fatalf("Got transaction Date %s, wanted 2022-01-01T00:00:00-05:00", x[0].Attributes.Transactions[0].Date)
	}
	expectedAmount, _ := decimal.NewFromString("4.50")
	if !x[0].Attributes.Transactions[0].Amount.Equals(expectedAmount) {
		t.Fatalf("Got transaction Amount %s, wanted 53.97", x[0].Attributes.Transactions[0].Amount)
	}
	if x[0].Attributes.Transactions[0].Description != "Interest" {
		t.Fatalf("Got transaction Description %s, wanted Interest", x[0].Attributes.Transactions[0].Description)
	}
	if x[0].Attributes.Transactions[0].CategoryID != "24" {
		t.Fatalf("Got transaction CategoryID %s, wanted 24", x[0].Attributes.Transactions[0].CategoryID)
	}
	if x[0].Attributes.Transactions[0].CategoryName != "Interest or Fees" {
		t.Fatalf("Got transaction CategoryName %s, wanted Interest or Fees", x[0].Attributes.Transactions[0].CategoryName)
	}
	if x[0].Attributes.Transactions[0].SourceID != "79" {
		t.Fatalf("Got transaction SourceID %s, wanted 79", x[0].Attributes.Transactions[0].SourceID)
	}
	if x[0].Attributes.Transactions[0].SourceName != "Bank" {
		t.Fatalf("Got transaction SourceName %s, wanted Bank", x[0].Attributes.Transactions[0].SourceName)
	}
	if x[0].Attributes.Transactions[0].DestinationID != "3" {
		t.Fatalf("Got transaction DestinationID %s, wanted 3", x[0].Attributes.Transactions[0].DestinationID)
	}
	if x[0].Attributes.Transactions[0].DestinationName != "Savings account" {
		t.Fatalf("Got transaction DestinationName %s, wanted Savings account", x[0].Attributes.Transactions[0].DestinationName)
	}
}

func TestFetchTransaction(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/2763", nil)
	f.HandleTxn(w, req)

	var x firefly.Transactions
	json.NewDecoder(w.Body).Decode(&x)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Status code = %d, want %d\n", w.Result().StatusCode, http.StatusOK)
	}
	if x.ID != "2763" {
		t.Fatalf("Got transaction ID %s, wanted 2763", x.ID)
	}

	if len(x.Attributes.Transactions) != 1 {
		t.Fatalf("Got %d attributes.Transactions, wanted 1", len(x.Attributes.Transactions))
	}

	if x.Attributes.Transactions[0].Type != "deposit" {
		t.Fatalf("Got transaction Type %s, wanted deposit", x.Attributes.Transactions[0].Type)
	}
	if x.Attributes.Transactions[0].Date != "2022-01-01T00:00:00-05:00" {
		t.Fatalf("Got transaction Date %s, wanted 2022-01-01T00:00:00-05:00", x.Attributes.Transactions[0].Date)
	}
	expectedAmount, _ := decimal.NewFromString("4.50")
	if !x.Attributes.Transactions[0].Amount.Equals(expectedAmount) {
		t.Fatalf("Got transaction Amount %s, wanted 53.97", x.Attributes.Transactions[0].Amount)
	}
	if x.Attributes.Transactions[0].Description != "Interest" {
		t.Fatalf("Got transaction Description %s, wanted Interest", x.Attributes.Transactions[0].Description)
	}
	if x.Attributes.Transactions[0].CategoryID != "24" {
		t.Fatalf("Got transaction CategoryID %s, wanted 24", x.Attributes.Transactions[0].CategoryID)
	}
	if x.Attributes.Transactions[0].CategoryName != "Interest or Fees" {
		t.Fatalf("Got transaction CategoryName %s, wanted Interest or Fees", x.Attributes.Transactions[0].CategoryName)
	}
	if x.Attributes.Transactions[0].SourceID != "79" {
		t.Fatalf("Got transaction SourceID %s, wanted 79", x.Attributes.Transactions[0].SourceID)
	}
	if x.Attributes.Transactions[0].SourceName != "Bank" {
		t.Fatalf("Got transaction SourceName %s, wanted Bank", x.Attributes.Transactions[0].SourceName)
	}
	if x.Attributes.Transactions[0].DestinationID != "3" {
		t.Fatalf("Got transaction DestinationID %s, wanted 3", x.Attributes.Transactions[0].DestinationID)
	}
	if x.Attributes.Transactions[0].DestinationName != "Savings account" {
		t.Fatalf("Got transaction DestinationName %s, wanted Savings account", x.Attributes.Transactions[0].DestinationName)
	}
}

func TestCreateTransaction(t *testing.T) {
	data := url.Values{}
	data.Set("date", "2022-01-01")
	data.Set("amount", "13.37")
	data.Set("description", "Mirror")
	data.Set("category_id", "4")
	data.Set("category_name", "Apartment")
	data.Set("source_name", "Savings accounts")
	data.Set("destination_name", "Structube")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	f.HandleTxn(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Success is a 302 redirect to the transactions page
	if w.Result().StatusCode != http.StatusFound {
		body, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("Status code = %d, want %d\n. Response body: %s", w.Result().StatusCode, http.StatusOK, body)
	}
}
