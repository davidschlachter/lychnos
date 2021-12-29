package categorybudget

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type CategoryBudget struct {
	ID       int             `json:"id"`
	Budget   int             `json:"budget"`
	Category int             `json:"category"`
	Amount   decimal.Decimal `json:"amount"`
}

type CategoryBudgets struct {
	db *sql.DB
}

func New(db *sql.DB) *CategoryBudgets {
	return &CategoryBudgets{db: db}
}

func (c *CategoryBudgets) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		hasID := regexp.MustCompile(`/[0-9]+$`)
		if hasID.MatchString(req.URL.Path) {
			c.fetch(w, req)
		} else {
			c.list(w, req)
		}
	case "POST":
		c.upsert(w, req)
	case "DELETE":
		c.delete(w, req)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

func (c *CategoryBudgets) fetch(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse categorybudget ID: %s\n", id)
		return
	}

	categoryBudgets, err := c.Fetch(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not fetch categorybudget: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryBudgets)
}

func (c *CategoryBudgets) Fetch(id string) ([]CategoryBudget, error) {
	const q = "SELECT id, budget, category, amount FROM category_budgets WHERE id = ?;"

	row := c.db.QueryRow(q, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var categoryBudgets []CategoryBudget

	var catBgt CategoryBudget
	row.Scan(&catBgt.ID, &catBgt.Budget, &catBgt.Category, &catBgt.Amount)
	categoryBudgets = append(categoryBudgets, catBgt)

	return categoryBudgets, nil
}

// TODO(davidschlachter): allow filtering, e.g. by budget ID
func (c *CategoryBudgets) list(w http.ResponseWriter, req *http.Request) {
	categoryBudgets, err := c.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not list categorybudget: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryBudgets)
}

func (c *CategoryBudgets) List() ([]CategoryBudget, error) {
	const q = "SELECT id, budget, category, amount FROM category_budgets;"
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryBudgets []CategoryBudget

	for rows.Next() {
		var catBgt CategoryBudget
		rows.Scan(&catBgt.ID, &catBgt.Budget, &catBgt.Category, &catBgt.Amount)
		categoryBudgets = append(categoryBudgets, catBgt)
	}

	return categoryBudgets, nil
}

func (c *CategoryBudgets) delete(w http.ResponseWriter, req *http.Request) {
	const q = "DELETE FROM category_budgets WHERE id = ?;"

	var (
		id  int
		err error
	)

	idStr := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	id, err = strconv.Atoi(idStr)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid ID: %s\n", idStr)
		return
	}

	_, err = c.db.Exec(q, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not delete categorybudget: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CategoryBudgets) upsert(w http.ResponseWriter, req *http.Request) {
	const (
		q = "INSERT INTO category_budgets VALUES(?, ?, ?, ?)"
	)

	var (
		err                  error
		id, budget, category int
		amount               decimal.Decimal
	)

	err = req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not parse POST data\n")
		return
	}

	idStr := req.Form.Get("id")
	if len(idStr) == 0 {
		id = 0
	} else {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not parse ID: %s\n", idStr)
			return
		}
	}
	budgetStr := req.Form.Get("budget")
	budget, err = strconv.Atoi(budgetStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse budget ID: %s\n", budgetStr)
		return
	}
	categoryStr := req.Form.Get("category")
	category, err = strconv.Atoi(categoryStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse category ID: %s\n", categoryStr)
		return
	}
	amountStr := req.Form.Get("amount")
	amount, err = decimal.NewFromString(amountStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse amount: %s\n", amountStr)
		return
	}

	// TODO(davidschlachter): check that the new categorybudget has a valid
	// budget and category ID

	_, err = c.db.Exec(q, id, budget, category, amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not upsert categorybudget: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
