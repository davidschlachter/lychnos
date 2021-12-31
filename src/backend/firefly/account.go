package firefly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/shopspring/decimal"
)

type accountsResponse struct {
	Data  []Account `json:"data"`
	Meta  meta      `json:"meta"`
	Links links     `json:"links"`
}

type Account struct {
	ID         string            `json:"id"`
	Attributes AccountAttributes `json:"attributes"`
}

type AccountAttributes struct {
	Active         bool            `json:"active"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	CurrentBalance decimal.Decimal `json:"current_balance"`
}

func (f *Firefly) HandleAccount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		hasID := regexp.MustCompile(`/[0-9]+$`)
		if hasID.MatchString(req.URL.Path) {
			//f.fetchAccount(w, req)
		} else {
			f.listAccounts(w, req)
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

func (f *Firefly) listAccounts(w http.ResponseWriter, req *http.Request) {
	accounts, err := f.ListAccounts("")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not list accounts: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (f *Firefly) ListAccounts(accountType string) ([]Account, error) {
	const path = "/api/v1/accounts"

	var (
		results []Account
		err     error
	)

	if accountType == "" {
		accountType = "all"
	}
	page := 1

	for more := true; more; page++ {
		params := fmt.Sprintf("?type=%s&page=%d", accountType, page)
		req, _ := http.NewRequest("GET", f.url+path+params, nil)
		req.Header.Add("Authorization", "Bearer "+f.token)
		resp, err := f.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Accounts: %s", err)
		}
		defer resp.Body.Close()

		var accs accountsResponse
		json.NewDecoder(resp.Body).Decode(&accs)

		results = append(results, accs.Data...)

		more = accs.Meta.Pagination.CurrentPage < accs.Meta.Pagination.TotalPages
	}

	return results, err
}
