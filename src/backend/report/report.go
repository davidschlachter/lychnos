package report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/firefly"
	"github.com/davidschlachter/lychnos/src/backend/interval"
	"github.com/shopspring/decimal"
)

type Reports struct {
	f *firefly.Firefly
	c *categorybudget.CategoryBudgets
	b *budget.Budgets
}

func New(f *firefly.Firefly, c *categorybudget.CategoryBudgets, b *budget.Budgets) (*Reports, error) {
	if f == nil || c == nil || b == nil {
		return nil, fmt.Errorf("must provide valid clients")
	}
	return &Reports{
		f: f,
		c: c,
		b: b,
	}, nil
}

func (r *Reports) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		if strings.Contains(req.URL.Path, "/categorysummary/") {
			hasID := regexp.MustCompile(`/[0-9]+$`)
			if hasID.MatchString(req.URL.Path) {
				r.fetchCategorySummaries(w, req)
			} else {
				r.listCategorySummaries(w, req)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

type CategorySummary struct {
	firefly.Category
	CategoryBudgetID int             `json:"category_budget_id"`
	Amount           decimal.Decimal `json:"amount"`
	Sum              decimal.Decimal `json:"sum"`
	Start            time.Time       `json:"start"`
	End              time.Time       `json:"end"`
}

func (r *Reports) listCategorySummaries(w http.ResponseWriter, req *http.Request) {
	var (
		budget int
		err    error
	)

	// If a budget was not provided, fetch the current one.
	budgetStr, ok := req.URL.Query()["budget"]
	if !ok || len(budgetStr) == 0 {
		// get the current budget
		bgts, err := r.b.List()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not fetch budgets to find latest budget for report: %s\n", err)
			return
		}
		now := time.Now()
		for _, b := range bgts {
			if now.After(b.Start) && now.Before(b.End) {
				budget = b.ID
				break
			}
		}
		if budget == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not identify a current budget for summary")
			return
		}
	} else if len(budgetStr) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Got %d budget IDs, wanted 0 or 1", len(budgetStr))
		return
	} else {
		budget, err = strconv.Atoi(budgetStr[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not parse budget ID: %s\n", budgetStr[0])
			return
		}
	}

	summaries, err := r.ListCategorySummaries(budget)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not generate CategorySummaries: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func (r *Reports) ListCategorySummaries(budgetID int) ([]CategorySummary, error) {
	budget, err := r.b.Fetch(strconv.Itoa(budgetID))
	if err != nil || len(budget) != 1 {
		return nil, fmt.Errorf("could not find budget with ID = %d", budgetID)
	}
	categorybudgets, err := r.c.List()
	if err != nil {
		return nil, fmt.Errorf("could not list categorybudgets: %s", err)
	}
	categories, err := r.f.CachedListCategoryTotals(budget[0].Start, budget[0].End)
	if err != nil {
		return nil, fmt.Errorf("could not list Categories: %s", err)
	}

	var results []CategorySummary

	for _, c := range categorybudgets {
		if c.Budget != budgetID {
			continue
		}
		var cs CategorySummary
		cs.ID = c.Category
		cs.CategoryBudgetID = c.ID
		cs.Amount = c.Amount
		cs.Start = budget[0].Start
		cs.End = budget[0].End
		results = append(results, cs)
	}

	// TODO(davidschlachter): n^2 complexity. Shouldn't have too many
	// Categories, but this could be improved.
	for i := range results {
		for j := range categories {
			if categories[j].ID == results[i].ID {
				results[i].Name = categories[j].Name
				results[i].Sum = categories[j].Earned.Add(categories[j].Spent)
				break
			}
		}
	}

	return results, nil
}

type CategorySummaryDetail struct {
	CategorySummary
	Totals []firefly.CategoryTotal `json:"totals"`
}

func (r *Reports) fetchCategorySummaries(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse categorybudget ID: %s\n", idStr)
		return
	}

	summary, err := r.FetchCategorySummary(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not generate CategorySummary: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (r *Reports) FetchCategorySummary(catBgtID int) ([]CategorySummaryDetail, error) {
	catBgt, err := r.c.Fetch(strconv.Itoa(catBgtID))
	if err != nil || len(catBgt) != 1 {
		return nil, fmt.Errorf("could not fetch categorybudgets: %s", err)
	}
	categories, err := r.f.CachedCategories()
	if err != nil {
		return nil, fmt.Errorf("could not list Categories: %s", err)
	}
	budget, err := r.b.Fetch(strconv.Itoa(catBgt[0].Budget))
	if err != nil || len(budget) != 1 {
		return nil, fmt.Errorf("could not find budget with ID = %d", catBgt[0].Budget)
	}

	// Populate the name and category ID
	var cs CategorySummary
	for _, c := range categories {
		if c.ID != catBgt[0].Category {
			continue
		}
		cs.ID = c.ID
		cs.Name = c.Name
		cs.Start = budget[0].Start
		cs.End = budget[0].End
		break
	}
	if cs.ID == 0 && cs.Name == "" {
		return nil, fmt.Errorf("could not find categorysummary")
	}
	cs.Amount = catBgt[0].Amount
	results := []CategorySummaryDetail{{CategorySummary: cs}}

	// Fetch the summaries for each month, from the start of the budget to the
	// current month
	if budget[0].ReportingInterval != 0 {
		return nil, fmt.Errorf("unknown reporting interval %d, only monthly is supported", budget[0].ReportingInterval)
	}

	intervals := interval.Get(budget[0].Start, budget[0].End, time.Now().UTC().Location())

	for _, i := range intervals {
		ct, err := r.f.CachedFetchCategoryTotals(cs.ID, i.Start, i.End)
		if err != nil {
			return nil, fmt.Errorf("could not generate category summary: %s", err)
		}
		results[0].Totals = append(results[0].Totals, ct[0])
	}

	var sum decimal.Decimal
	for _, t := range results[0].Totals {
		sum = sum.Add(t.Earned.Add(t.Spent))
	}
	results[0].Sum = sum

	return results, nil
}
