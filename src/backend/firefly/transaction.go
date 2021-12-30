package firefly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func (f *Firefly) HandleTxn(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		f.listTxns(w, req)
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
	SourceID        string          `json:"source_id,omitempty"`
	SourceName      string          `json:"source_name,omitempty"`
	DestinationID   string          `json:"destination_id,omitempty"`
	DestinationName string          `json:"destination_name,omitempty"`
}

type createRequest struct {
	Transactions []Transaction `json:"transactions"`
}

func (f *Firefly) createTxn(w http.ResponseWriter, req *http.Request) {
	// Decode the request
	decoder := json.NewDecoder(req.Body)
	var t Transaction
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not parse POST body: %s", err)
		return
	}

	//
	// Validate the request
	//
	cats, _ := f.Categories()
	var ok bool
	for _, c := range cats {
		if strconv.Itoa(c.ID) == t.CategoryID {
			ok = true
			break
		}
	}
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not find Category with ID = '%s'", t.CategoryID)
		return
	}

	dateFormat := "2006-01-02T15:04:05-07:00"
	txnDate, err := time.Parse(dateFormat, t.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse date '%s'", t.Date)
		return
	}

	if t.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "description must be provided")
		return
	}

	if t.SourceID == "" && t.SourceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "source_id or source_name must be provided")
		return
	}

	if t.DestinationID == "" && t.DestinationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "destination_id or destination_name must be provided")
		return
	}

	if t.Amount.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "amount must be provided")
		return
	}

	t.Type = f.calcTxnType(t.SourceID, t.SourceName, t.DestinationID, t.DestinationName)
	if t.Type == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not determine transaction type")
		return
	}

	// Send to the firefly API
	doc := createRequest{
		Transactions: []Transaction{t},
	}

	const path = "/api/v1/transactions"

	body, err := json.Marshal(doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not marshal CreateRequest: %s", err)
		return
	}
	r, _ := http.NewRequest("POST", f.url+path, bytes.NewBuffer(body))
	r.Header.Add("Authorization", "Bearer "+f.token)
	r.Header.Add("Content-Type", "application/json")
	resp, err := f.client.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not create transaction: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not create transaction: %s", resp.Status)
		return
	}

	respBody, _ := io.ReadAll(resp.Body)

	// Invalidate any matching cache entries. Since the transaction was
	// successfully create, these conversions should not raise errors
	catID, _ := strconv.Atoi(t.CategoryID)
	key := categoryTotalsKey{
		CategoryID: catID,
		Start:      txnDate,
		End:        txnDate,
	}
	f.invalidateCategoryCache(key)
	f.invalidateAccountsCache()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(respBody))
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
	// New accounts are expense accounts
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

	txns, err := f.ListTransactions(page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not list transactions: %s", err)
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
	for _, t := range txns.Data {
		results = append(results, t)
	}

	return results, nil
}
