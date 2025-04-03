package firefly

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/davidschlachter/lychnos/src/backend/httperror"
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
	Active          bool            `json:"active"`
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	CurrentBalance  decimal.Decimal `json:"current_balance"`
	IncludeNetWorth bool            `json:"include_net_worth"`
}

func (f *Firefly) HandleAccount(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.RequestURI)
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
	accounts, err := f.CachedAccounts()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list accounts: %s", err))
		return
	}

	// If a type parameter was provided, filter the returned accounts
	acctTypes, ok := req.URL.Query()["type"]
	if ok && len(acctTypes) > 0 {
		var filtered_accounts []Account
		for _, t := range acctTypes {
			for _, a := range accounts {
				if a.Attributes.Type != t {
					continue
				}
				filtered_accounts = append(filtered_accounts, a)
			}
		}
		accounts = filtered_accounts
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
		req, _ := http.NewRequest("GET", f.config.URL+path+params, nil)
		req.Header.Add("Authorization", "Bearer "+f.config.Token)
		resp, err := f.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Accounts: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("got status %d", resp.StatusCode)
		}

		var accs accountsResponse
		json.NewDecoder(resp.Body).Decode(&accs)

		results = append(results, accs.Data...)

		more = accs.Meta.Pagination.CurrentPage < accs.Meta.Pagination.TotalPages
	}

	// Only include 'active' accounts.
	filteredResults := make([]Account, 0, len(results))
	for i := range results {
		if results[i].Attributes.Active {
			filteredResults = append(filteredResults, results[i])
		}
	}

	return filteredResults, err
}
