package firefly

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/davidschlachter/lychnos/src/backend/httperror"
)

type bigPicture struct {
	Expenses12Months decimal.Decimal `json:"expenses_twelve_months"`
	Income12Months   decimal.Decimal `json:"income_twelve_months"`
	Taxes12Months    decimal.Decimal `json:"taxes_twelve_months"`
	Income3Months    decimal.Decimal `json:"income_three_months"`
	Expenses3Months  decimal.Decimal `json:"expenses_three_months"`
	NetWorth         decimal.Decimal `json:"net_worth"`
}

type insight struct {
	Difference decimal.Decimal `json:"difference"`
}

func (f *Firefly) HandleBigPicture(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.RequestURI)
	switch req.Method {
	case "GET":
		if err := f.bigPicture(w); err != nil {
			httperror.Send(w, req, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

func (f *Firefly) bigPicture(w http.ResponseWriter) error {
	bp, err := f.CachedBigPicture()
	if err != nil {
		return fmt.Errorf("loading big picture: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bp)

	return nil
}

func (f *Firefly) fetchBigPicture() (*bigPicture, error) {
	// This logic is based on my assumptions for retirement planning. If you're
	// using this yourself, you'll probably want different assumptions (time to
	// fork this repo!).
	var bp bigPicture

	// Get our current net worth.
	accounts, err := f.CachedAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not list accounts: %s", err)
	}

	for _, a := range accounts {
		if a.Attributes.Type != AcctTypeAsset {
			continue
		}
		bp.NetWorth = bp.NetWorth.Add(a.Attributes.CurrentBalance)
	}

	// Get income/expense data, starting with the most recent three months.
	bp.Income3Months, err = f.insight("income", time.Now().AddDate(0, -3, 0), time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not get three month income: %s", err)
	}
	bp.Expenses3Months, err = f.insight("expense", time.Now().AddDate(0, -3, 0), time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not get three month expenses: %s", err)
	}

	// Now for the last 12 months.
	bp.Income12Months, err = f.insight("income", time.Now().AddDate(-1, 0, 0), time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not get twelve month income: %s", err)
	}
	bp.Expenses12Months, err = f.insight("expense", time.Now().AddDate(-1, 0, 0), time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not get twelve month expenses: %s", err)
	}

	// Also include what we've paid in taxes over the last twelve months, which
	// offsets income in the frontend calculations.
	cts, err := f.CachedCategories()
	if err != nil {
		return nil, fmt.Errorf("listing cached categories: %s", err)
	}
	categoryID := -1
	for _, c := range cts {
		if strings.ToLower(c.Name) == "taxes" {
			categoryID = c.ID
			break
		}
	}
	if categoryID == -1 {
		return nil, fmt.Errorf("finding tax category: %s", err)
	}
	taxTotals, err := f.FetchCategoryTotal(categoryID, time.Now().AddDate(-1, 0, 0), time.Now())
	if err != nil {
		return nil, fmt.Errorf("listing tax totals: %s", err)
	}
	if len(taxTotals) != 1 {
		return nil, fmt.Errorf("expected tax array to contain 1 item, contained %d", len(taxTotals))
	}
	bp.Taxes12Months = taxTotals[0].Earned.Add(taxTotals[0].Spent)

	return &bp, nil
}

func (f *Firefly) insight(transactionType string, start time.Time, end time.Time) (decimal.Decimal, error) {
	path := fmt.Sprintf("/api/v1/insight/%s/total?start=%s&end=%s", transactionType, start.Format("2006-01-02"), end.Format("2006-01-02"))

	req, err := http.NewRequest("GET", f.url+path, nil)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Authorization", "Bearer "+f.token)

	resp, err := f.client.Do(req)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to fetch insight: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decimal.Decimal{}, fmt.Errorf("got status %d", resp.StatusCode)
	}

	var i []insight
	json.NewDecoder(resp.Body).Decode(&i)
	if len(i) != 1 {
		return decimal.Decimal{}, fmt.Errorf("expected array to contain 1 item, contained %d", len(i))
	}
	return i[0].Difference, nil
}
