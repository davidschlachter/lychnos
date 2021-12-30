package firefly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (f *Firefly) HandleTxn(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		f.createTxn(w, req)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type createRequest struct {
	Transactions []createRequestTransaction `json:"transactions"`
}

type createRequestTransaction struct {
	Type            string `json:"type"`
	Date            string `json:"date"` // "2018-09-17T12:46:47+01:00"
	Amount          string `json:"amount"`
	Description     string `json:"description"`
	CategoryID      string `json:"category_id"`
	SourceID        string `json:"source_id"`
	SourceName      string `json:"source_name"`
	DestinationID   string `json:"destination_id"`
	DestinationName string `json:"destination_name"`
}

func (f *Firefly) createTxn(w http.ResponseWriter, req *http.Request) {
	// Decode the request
	decoder := json.NewDecoder(req.Body)
	var t createRequestTransaction
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

	if t.Type != "withdrawal" && t.Type != "deposit" && t.Type != "transfer" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid transaction type '%s'", t.Type)
		return
	}

	dateFormat := "2006-01-02T15:04:05-07:00"
	_, err = time.Parse(dateFormat, t.Date)
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

	if t.Amount == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "amount must be provided")
		return
	}

	// Send to the firefly API
	doc := createRequest{
		Transactions: []createRequestTransaction{t},
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

	f.invalidateCategoryCache(t.CategoryID)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(respBody))
}
