package firefly_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/davidschlachter/lychnos/src/backend/firefly"
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
	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"type":"accounts","id":"464","attributes":{"created_at":"2021-09-21T19:59:20-04:00","updated_at":"2021-09-21T19:59:20-04:00","active":true,"order":null,"name":"1Password","type":"expense","account_role":null,"currency_id":"9","currency_code":"CAD","currency_symbol":"C$","currency_decimal_places":2,"current_balance":"53.97","current_balance_date":"2022-01-01T23:59:59-05:00","notes":null,"monthly_payment_date":null,"credit_card_type":null,"account_number":null,"iban":null,"bic":null,"virtual_balance":"0.00","openinterest":null,"interest_period":null,"current_debt":null,"include_net_worth":true,"longitude":null,"latitude":null,"zoom_level":null},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/accounts\/464","0":{"rel":"self","uri":"\/accounts\/464"}}},{"type":"accounts","id":"387","attributes":{"created_at":"2021-05-26T13:14:09-04:00","updated_at":"2021-05-26T13:14:09-04:00","active":true,"order":null,"name":"Savings accounts","type":"asset","account_role":null,"currency_id":"9","currency_code":"CAD","currency_symbol":"C$","currency_decimal_places":2,"current_balance":"1.00","current_balance_date":"2022-01-01T23:59:59-05:00","notes":null,"monthly_payment_date":null,"credit_card_type":null,"account_number":null,"iban":null,"bic":null,"virtual_balance":"0.00","opening_balance":"0.00","opening_balance_date":null,"liability_type":null,"liability_direction":null,"interest":null,"interest_period":null,"current_debt":null,"include_net_worth":true,"longitude":null,"latitude":null,"zoom_level":null},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/accounts\/387","0":{"rel":"self","uri":"\/accounts\/387"}}}],"meta":{"pagination":{"total":2,"count":2,"per_page":2,"current_page":1,"total_pages":1}},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/accounts?type=all&page=1","first":"http:\/\/192.168.6.4:8753\/api\/v1\/accounts?type=all&page=1","next":"","last":"http:\/\/192.168.6.4:8753\/api\/v1\/accounts?type=all&page=1"}}`))
	})

	server = httptest.NewServer(mux)
}
