package firefly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/httperror"
	"github.com/shopspring/decimal"
)

func (f *Firefly) HandleTxn(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		hasID := regexp.MustCompile(`/[0-9]+$`)
		if hasID.MatchString(req.URL.Path) {
			f.fetchTxn(w, req)
		} else {
			f.listTxns(w, req)
		}
	case "POST":
		f.createTxn(w, req)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type txnsResponse struct {
	Data  []Transactions `json:"data"`
	Meta  meta           `json:"meta"`
	Links links          `json:"links"`
}

type Transactions struct {
	ID         string                `json:"id"`
	Attributes TransactionAttributes `json:"attributes"`
}

type TransactionAttributes struct {
	GroupTitle   string        `json:"group_title"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Type            string          `json:"type"`
	Date            string          `json:"date"` // "2018-09-17T12:46:47+01:00"
	Amount          decimal.Decimal `json:"amount"`
	Description     string          `json:"description"`
	CategoryID      string          `json:"category_id,omitempty"`
	CategoryName    string          `json:"category_name"`
	SourceID        string          `json:"source_id,omitempty"`
	SourceName      string          `json:"source_name,omitempty"`
	DestinationID   string          `json:"destination_id,omitempty"`
	DestinationName string          `json:"destination_name,omitempty"`
}

type createRequest struct {
	Transactions []Transaction `json:"transactions"`
}

func (f *Firefly) createTxn(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, "Could not parse POST data")
		return
	}

	// Build the transaction struct
	amt, err := decimal.NewFromString(req.Form.Get("amount"))
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not parse amount: %s", req.Form.Get("amount")))
		return
	}
	t := Transaction{
		Date:            req.Form.Get("date"),
		Amount:          amt,
		Description:     req.Form.Get("description"),
		CategoryID:      req.Form.Get("category_id"),
		CategoryName:    req.Form.Get("category_name"),
		SourceID:        req.Form.Get("source_id"),
		SourceName:      req.Form.Get("source_name"),
		DestinationID:   req.Form.Get("destination_id"),
		DestinationName: req.Form.Get("destination_name"),
	}

	//
	// Validate the request
	//
	// Verify that a provided category ID is valid. If only a category name is
	// provided, add the ID. Allow an empty category (e.g. for a transfer).
	if t.CategoryID != "" || t.CategoryName != "" {
		cats, _ := f.CachedCategories()
		var ok bool
		for _, c := range cats {
			if strconv.Itoa(c.ID) == t.CategoryID {
				ok = true
				break
			}
			if c.Name == t.CategoryName {
				t.CategoryID = strconv.Itoa(c.ID)
				ok = true
				break
			}
		}
		if !ok {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not find Category with ID = '%s' or Name = '%s'", t.CategoryID, t.CategoryName))
			return
		}
	}

	dateFormat := "2006-01-02"
	// firefly internal dateFormat := "2006-01-02T15:04:05-07:00"
	txnDate, err := time.Parse(dateFormat, t.Date)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse date '%s'", t.Date))
		return
	}

	if t.Description == "" {
		httperror.Send(w, req, http.StatusBadRequest, "description must be provided")
		return
	}

	if t.SourceID == "" && t.SourceName == "" {
		httperror.Send(w, req, http.StatusBadRequest, "source_id or source_name must be provided")
		return
	}

	if t.DestinationID == "" && t.DestinationName == "" {
		httperror.Send(w, req, http.StatusBadRequest, "destination_id or destination_name must be provided")
		return
	}

	if t.Amount.IsZero() {
		httperror.Send(w, req, http.StatusBadRequest, "amount must be provided")
		return
	}

	// Determine the transaction type
	t.SourceID, t.SourceName = f.resolveAccount(t.SourceID, t.SourceName)
	t.DestinationID, t.DestinationName = f.resolveAccount(t.DestinationID, t.DestinationName)
	t.Type = f.calcTxnType(t.SourceID, t.SourceName, t.DestinationID, t.DestinationName)
	if t.Type == "" {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not determine transaction type with provided account information: sourceID: %s, sourceName: %s; destID: %s, destName: %s\n", t.SourceID, t.SourceName, t.DestinationID, t.DestinationName))
		return
	}

	// Send to the firefly API
	doc := createRequest{
		Transactions: []Transaction{t},
	}

	const path = "/api/v1/transactions"

	body, err := json.Marshal(doc)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not marshal CreateRequest: %s", err))
		return
	}
	r, _ := http.NewRequest("POST", f.url+path, bytes.NewBuffer(body))
	r.Header.Add("Authorization", "Bearer "+f.token)
	r.Header.Add("Content-Type", "application/json")
	resp, err := f.client.Do(r)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not create transaction: %s", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not create transaction, got status %d %s", resp.StatusCode, resp.Status))
		return
	}

	// Check for successful response
	var result struct {
		Data Transactions `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Data.ID == "" {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not create transaction: %s", err))
		return
	}

	// Invalidate any matching cache entries. Since the transaction was
	// successfully created, the conversions should not raise errors
	catID, _ := strconv.Atoi(t.CategoryID)
	key := categoryTotalsKey{
		CategoryID: catID,
		Start:      txnDate,
		End:        txnDate,
	}
	f.refreshTransactions(1) // since user is going to txns page next, update now
	go func() {              // we can update other caches after returning
		f.refreshCategoryTxnCache(key)
		_ = f.refreshAccounts()
	}()

	// Successful txn creation should redirect the client to the transactions page
	http.Redirect(w, r, "/app/txns", http.StatusFound)
}

// resolveAccount will determine the ID of an account, provided a name; or the
// name, provided an ID.
func (f *Firefly) resolveAccount(id, name string) (string, string) {
	// Both name and ID missing or provided
	if (id == "" && name == "") || (id != "" && name != "") {
		return id, name
	}

	accts, _ := f.CachedAccounts()

	// Name provided, ID missing
	if id == "" && name != "" {
		for _, a := range accts {
			if a.Attributes.Name != name {
				continue
			}
			return a.ID, name
		}
	}

	// ID provided, name missing
	if id != "" && name == "" {
		for _, a := range accts {
			if a.ID != id {
				continue
			}
			return id, a.Attributes.Name
		}
	}

	return id, name
}

const (
	AcctTypeAsset   = "asset"
	AcctTypeExpense = "expense"
	AcctTypeRevenue = "revenue"
)

// calcTxnType determines whether the transaction is a deposit, withdrawal, or
// transfer. If source is an asset account and dest is not: withdrawal. If dest
// is an asset account but source is not: deposit. If both accounts are of the
// same type: transfer.
//
// TODO(davidschlachter): this may be confused if we have two accounts with the
// same name but different types, e.g. expense and revenue
func (f *Firefly) calcTxnType(srcID, srcName, destID, destName string) string {
	var srcType, destType string
	accts, _ := f.CachedAccounts()
	// Find type of existing accounts
	for _, a := range accts {
		switch a.ID {
		case srcID:
			srcType = a.Attributes.Type
		case destID:
			destType = a.Attributes.Type
		}
		if srcType != "" && destType != "" {
			break
		}
	}
	// TODO(davidschlachter): maybe support cash accounts one day
	if srcType == "" && srcName != "" && destType == AcctTypeAsset {
		srcType = AcctTypeRevenue
	}
	if destType == "" && destName != "" && srcType == AcctTypeAsset {
		destType = AcctTypeExpense
	}
	// Determine transaction type
	if srcType == AcctTypeAsset && destType == AcctTypeExpense {
		return "withdrawal"
	} else if srcType == AcctTypeRevenue && destType == AcctTypeAsset {
		return "deposit"
	} else if srcType == AcctTypeAsset && destType == AcctTypeAsset {
		return "transfer"
	} else {
		return ""
	}
}

func (f *Firefly) listTxns(w http.ResponseWriter, req *http.Request) {
	var page int
	pageStr, ok := req.URL.Query()["page"]
	if ok && len(pageStr) > 0 {
		page, _ = strconv.Atoi(pageStr[0]) // if page cannot be parsed, we'll return page 1
	}

	txns, err := f.CachedTransactions(page)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list transactions: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txns)
}

func (f *Firefly) ListTransactions(page int) ([]Transactions, error) {
	const path = "/api/v1/transactions"

	if page == 0 {
		page = 1
	}

	params := fmt.Sprintf("?page=%d", page)
	req, _ := http.NewRequest("GET", f.url+path+params, nil)
	req.Header.Add("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Transactions: %s", err)
	}
	defer resp.Body.Close()

	var txns txnsResponse
	json.NewDecoder(resp.Body).Decode(&txns)

	var results []Transactions
	results = append(results, txns.Data...)

	return results, nil
}

func (f *Firefly) fetchTxn(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	if _, err := strconv.Atoi(id); err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse transaction ID: %s", id))
		return
	}

	txn, err := f.FetchTransaction(id)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not fetch transaction: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func (f *Firefly) FetchTransaction(id string) (*Transactions, error) {
	const path = "/api/v1/transactions/"

	req, _ := http.NewRequest("GET", f.url+path+id, nil)
	req.Header.Add("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Transaction: %s", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data Transactions `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Data.ID == "" {
		return nil, fmt.Errorf("no transaction found")
	}

	return &result.Data, nil
}
