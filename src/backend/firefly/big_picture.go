package firefly

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shopspring/decimal"

	"github.com/davidschlachter/lychnos/src/backend/httperror"
)

type bigPicture struct {
	Expenses12Months decimal.Decimal `json:"expenses_twelve_months"`
	Income12Months   decimal.Decimal `json:"income_twelve_months"`

	Income3Months   decimal.Decimal `json:"income_three_months"`
	Expenses3Months decimal.Decimal `json:"expenses_three_months"`

	NetWorth decimal.Decimal `json:"net_worth"`
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
	// Note that any transactions without a category will be ignored in the 'Big
	// Picture' summary.
	var bp bigPicture

	// Get our current net worth.
	accounts, err := f.CachedAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not list accounts: %s", err)
	}

	for _, a := range accounts {
		if a.Attributes.Type != AcctTypeAsset || !a.Attributes.IncludeNetWorth {
			continue
		}
		bp.NetWorth = bp.NetWorth.Add(a.Attributes.CurrentBalance)
	}

	categories, err := f.CachedCategories()
	if err != nil {
		return nil, fmt.Errorf("getting cached categories: %w", err)
	}

	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)
	twelveMonthsAgo := now.AddDate(-1, 0, 0)

	for _, c := range categories {
		if _, ok := f.config.BigPictureIgnore[c.ID]; ok {
			continue // Don't include in totals.
		}
		_, alwaysIncome := f.config.BigPictureIncome[c.ID]

		// Last three months
		categoryTotalThreeMonths, err := f.FetchCategoryTotal(c.ID, threeMonthsAgo, now)
		if err != nil {
			return nil, fmt.Errorf("listing '%s' three-month totals: %s", c.Name, err)
		}
		if len(categoryTotalThreeMonths) != 1 {
			return nil, fmt.Errorf("expected '%s' array to contain 1 item, contained %d", c.Name, len(categoryTotalThreeMonths))
		}

		// Last twelve months
		categoryTotalTwelveMonths, err := f.FetchCategoryTotal(c.ID, twelveMonthsAgo, now)
		if err != nil {
			return nil, fmt.Errorf("listing '%s' twelve-month totals: %s", c.Name, err)
		}
		if len(categoryTotalTwelveMonths) != 1 {
			return nil, fmt.Errorf("expected '%s' array to contain 1 item, contained %d", c.Name, len(categoryTotalTwelveMonths))
		}

		categorySum3Months := categoryTotalThreeMonths[0].Earned.Add(categoryTotalThreeMonths[0].Spent)
		categorySum12Months := categoryTotalTwelveMonths[0].Earned.Add(categoryTotalTwelveMonths[0].Spent)

		if alwaysIncome || categoryTotalThreeMonths[0].Earned.Abs().GreaterThan(categoryTotalThreeMonths[0].Spent.Abs()) {
			// An income category
			bp.Income3Months = bp.Income3Months.Add(categorySum3Months)
			bp.Income12Months = bp.Income12Months.Add(categorySum12Months)
		} else {
			// An expense category
			bp.Expenses3Months = bp.Expenses3Months.Add(categorySum3Months)
			bp.Expenses12Months = bp.Expenses12Months.Add(categorySum12Months)
		}
	}

	return &bp, nil
}
