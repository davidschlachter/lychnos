package report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/categorybudget"
	"github.com/davidschlachter/lychnos/src/backend/firefly"
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
		if strings.HasSuffix(req.URL.Path, "/categorysummary/") {
			r.categorySummaries(w, req)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

type CategorySummary struct {
	firefly.Category
	Amount decimal.Decimal `json:"amount"`
	Sum    decimal.Decimal `json:"sum"`
}

func (r *Reports) categorySummaries(w http.ResponseWriter, req *http.Request) {
	budgetStr, ok := req.URL.Query()["budget"]
	if !ok || len(budgetStr) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "budget parameter must be provided")
		return
	}
	budget, err := strconv.Atoi(budgetStr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse budget ID: %s\n", budgetStr[0])
		return
	}
	summaries, err := r.CategorySummaries(budget)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not generate CategorySummaries: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func (r *Reports) CategorySummaries(budgetID int) ([]CategorySummary, error) {
	budget, err := r.b.Fetch(strconv.Itoa(budgetID))
	if err != nil || len(budget) != 1 {
		return nil, fmt.Errorf("could not find budget with ID = %d", budgetID)
	}
	categorybudgets, err := r.c.List()
	if err != nil {
		return nil, fmt.Errorf("could not list categorybudgets: %s", err)
	}
	categories, err := r.f.CategoryTotals(budget[0].Start, budget[0].End)
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
		cs.Amount = c.Amount
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
