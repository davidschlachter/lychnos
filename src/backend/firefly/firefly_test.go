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
	mux.HandleFunc("/api/v1/transactions/2763", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"type":"transactions","id":"2763","attributes":{"created_at":"2022-01-01T10:11:24-05:00","updated_at":"2022-01-01T10:11:24-05:00","user":"1","group_title":null,"transactions":[{"user":"1","transaction_journal_id":"2809","type":"deposit","date":"2022-01-01T00:00:00-05:00","order":0,"currency_id":"9","currency_code":"CAD","currency_name":"Canadian dollar","currency_symbol":"C$","currency_decimal_places":2,"foreign_currency_id":"0","foreign_currency_code":null,"foreign_currency_symbol":null,"foreign_currency_decimal_places":0,"amount":"4.500000000000000000000000","foreign_amount":null,"description":"Interest","source_id":"79","source_name":"Bank","source_iban":null,"source_type":"Revenue account","destination_id":"3","destination_name":"Savings account","destination_iban":"","destination_type":"Asset account","budget_id":"0","budget_name":null,"category_id":"24","category_name":"Interest or Fees","bill_id":null,"bill_name":null,"reconciled":false,"notes":null,"tags":[],"internal_reference":null,"external_id":null,"original_source":"ff3-v5.6.2|api-v1.5.4","recurrence_id":null,"recurrence_total":null,"recurrence_count":null,"bunq_payment_id":null,"external_uri":null,"import_hash_v2":"f776fdea04fa0854fa33a1a2c75660e291fb48114f8d781ae916f6c5c40b3dc3","sepa_cc":null,"sepa_ct_op":null,"sepa_ct_id":null,"sepa_db":null,"sepa_country":null,"sepa_ep":null,"sepa_ci":null,"sepa_batch_id":null,"interest_date":null,"book_date":null,"process_date":null,"due_date":null,"payment_date":null,"invoice_date":null,"longitude":null,"latitude":null,"zoom_level":null}]},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions\/2763","0":{"rel":"self","uri":"\/transactions\/2763"}}}}`))
	})
	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte(`{"data":[{"type":"transactions","id":"2763","attributes":{"created_at":"2022-01-01T10:11:24-05:00","updated_at":"2022-01-01T10:11:24-05:00","user":"1","group_title":null,"transactions":[{"user":"1","transaction_journal_id":"2809","type":"deposit","date":"2022-01-01T00:00:00-05:00","order":0,"currency_id":"9","currency_code":"CAD","currency_name":"Canadian dollar","currency_symbol":"C$","currency_decimal_places":2,"foreign_currency_id":"0","foreign_currency_code":null,"foreign_currency_symbol":null,"foreign_currency_decimal_places":0,"amount":"4.500000000000000000000000","foreign_amount":null,"description":"Interest","source_id":"79","source_name":"Bank","source_iban":null,"source_type":"Revenue account","destination_id":"3","destination_name":"Savings account","destination_iban":"","destination_type":"Asset account","budget_id":"0","budget_name":null,"category_id":"24","category_name":"Interest or Fees","bill_id":null,"bill_name":null,"reconciled":false,"notes":null,"tags":[],"internal_reference":null,"external_id":null,"original_source":"ff3-v5.6.2|api-v1.5.4","recurrence_id":null,"recurrence_total":null,"recurrence_count":null,"bunq_payment_id":null,"external_uri":null,"import_hash_v2":"f776fdea04fa0854fa33a1a2c75660e291fb48114f8d781ae916f6c5c40b3dc3","sepa_cc":null,"sepa_ct_op":null,"sepa_ct_id":null,"sepa_db":null,"sepa_country":null,"sepa_ep":null,"sepa_ci":null,"sepa_batch_id":null,"interest_date":null,"book_date":null,"process_date":null,"due_date":null,"payment_date":null,"invoice_date":null,"longitude":null,"latitude":null,"zoom_level":null}]},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions\/2763","0":{"rel":"self","uri":"\/transactions\/2763"}}}],"meta":{"pagination":{"total":1,"count":1,"per_page":1,"current_page":1,"total_pages":1}},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions?type=default&page=1","first":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions?type=default&page=1","next":"","last":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions?type=default&page=1"}}`))
		case "POST":
			w.Write([]byte(`{"data":{"type":"transactions","id":"2774","attributes":{"created_at":"2022-01-01T23:39:35-05:00","updated_at":"2022-01-01T23:39:35-05:00","user":"1","group_title":null,"transactions":[{"user":"1","transaction_journal_id":"2820","type":"withdrawal","date":"2022-01-01T00:00:00-05:00","order":0,"currency_id":"9","currency_code":"CAD","currency_name":"Canadian dollar","currency_symbol":"C$","currency_decimal_places":2,"foreign_currency_id":"0","foreign_currency_code":null,"foreign_currency_symbol":null,"foreign_currency_decimal_places":0,"amount":"13.370000000000000000000000","foreign_amount":null,"description":"Mirror","source_id":"3","source_name":"Savings accounts","source_iban":"","source_type":"Asset account","destination_id":"529","destination_name":"Structube","destination_iban":null,"destination_type":"Expense account","budget_id":"0","budget_name":null,"category_id":"4","category_name":"Apartment","bill_id":null,"bill_name":null,"reconciled":false,"notes":null,"tags":[],"internal_reference":null,"external_id":null,"original_source":"ff3-v5.6.2|api-v1.5.4","recurrence_id":null,"recurrence_total":null,"recurrence_count":null,"bunq_payment_id":null,"external_uri":null,"import_hash_v2":"599815725d6b01876c21e41b650d981a15b76e0a622633b2af234a210d51616f","sepa_cc":null,"sepa_ct_op":null,"sepa_ct_id":null,"sepa_db":null,"sepa_country":null,"sepa_ep":null,"sepa_ci":null,"sepa_batch_id":null,"interest_date":null,"book_date":null,"process_date":null,"due_date":null,"payment_date":null,"invoice_date":null,"longitude":null,"latitude":null,"zoom_level":null}]},"links":{"self":"http:\/\/192.168.6.4:8753\/api\/v1\/transactions\/2774","0":{"rel":"self","uri":"\/transactions\/2774"}}}}`))
		}
	})

	server = httptest.NewServer(mux)
}
