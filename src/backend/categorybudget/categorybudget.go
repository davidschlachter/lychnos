package categorybudget

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/budget"
	"github.com/davidschlachter/lychnos/src/backend/httperror"
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
	b  *budget.Budgets
}

func New(db *sql.DB, b *budget.Budgets) *CategoryBudgets {
	return &CategoryBudgets{db: db, b: b}
}

func (c *CategoryBudgets) Handle(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.RequestURI)
	switch req.Method {
	case "GET":
		hasID := regexp.MustCompile(`/[0-9]+$`)
		if hasID.MatchString(req.URL.Path) {
			c.fetch(w, req)
		} else {
			c.list(w, req)
		}
	case "POST":
		// Support multiple creates if the request type is application/json,
		// otherwise, assume an application/x-www-form-urlencoded request to
		// POST a single categorybudget.
		if req.Header.Get("Content-type") == "application/json" {
			c.upsertMultiple(w, req)
		} else {
			c.upsertSingle(w, req)
		}

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
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse categorybudget ID: %s", id))
		return
	}

	categoryBudgets, err := c.Fetch(id)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not fetch categorybudget: %s", err))
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
	// If a budget was not provided, fetch the current one.
	var (
		budget      int
		maxBudgetID int
		err         error
	)
	budgetStr, ok := req.URL.Query()["budget"]
	if !ok || len(budgetStr) == 0 {
		bgts, err := c.b.List()
		if err != nil {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not fetch budgets to find latest budget for list of category budgets: %s\n", err))
			return
		}
		now := time.Now()
		for _, b := range bgts {
			if b.ID > maxBudgetID {
				maxBudgetID = b.ID
			}
			if now.After(b.Start) && now.Before(b.End) {
				budget = b.ID
				break
			}
		}
		// If no budget exists, create one for the current year.
		if budget == 0 {
			budget = maxBudgetID + 1
			now := time.Now()
			err = c.b.Upsert(
				budget,
				time.Date(now.Year(), time.January, 01, 0, 0, 0, 0, time.Local),
				time.Date(now.Year(), time.December, 31, 0, 0, 0, 0, time.Local),
				0,
			)
			if err != nil {
				httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("No budget existed, failed to create a new budget: %s", err))
				return
			}
		}
	} else if len(budgetStr) > 1 {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Got %d budget IDs, wanted 0 or 1", len(budgetStr)))
		return
	} else {
		budget, err = strconv.Atoi(budgetStr[0])
		if err != nil {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse budget ID: %s\n", budgetStr[0]))
			return
		}
	}

	categoryBudgets, err := c.List()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list categorybudget: %s", err))
		return
	}

	var result []CategoryBudget
	for _, cb := range categoryBudgets {
		if cb.Budget == budget {
			result = append(result, cb)
		}
	}

	log.Printf("Returning category budgets for budget ID = %d\n", budget)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Invalid ID: %s", idStr))
		return
	}

	_, err = c.db.Exec(q, id)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not delete categorybudget: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CategoryBudgets) upsertSingle(w http.ResponseWriter, req *http.Request) {
	const (
		q = "INSERT INTO category_budgets (id, budget, category, amount) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE budget=VALUES(budget), category=VALUES(category), amount=VALUES(amount);"
	)

	var (
		err                  error
		id, budget, category int
		amount               decimal.Decimal
	)

	err = req.ParseForm()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, "Could not parse POST data")
		return
	}

	idStr := req.Form.Get("id")
	if len(idStr) == 0 {
		id = 0
	} else {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse ID: %s\n", idStr))
			return
		}
	}
	budgetStr := req.Form.Get("budget")
	budget, err = strconv.Atoi(budgetStr)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse budget ID: %s\n", budgetStr))
		return
	}
	categoryStr := req.Form.Get("category")
	category, err = strconv.Atoi(categoryStr)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse category ID: %s\n", categoryStr))
		return
	}
	amountStr := req.Form.Get("amount")
	amount, err = decimal.NewFromString(amountStr)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse amount: %s\n", amountStr))
		return
	}

	// TODO(davidschlachter): check that the new categorybudget has a valid
	// budget and category ID

	// TODO(davidschlachter): check that we have at most one categorybudget for
	// each category in a budget

	_, err = c.db.Exec(q, id, budget, category, amount)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not upsert categorybudget: %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// upsertMultiple will remove all CategoryBudgets for the current Budget,
// replacing them with the provided CategoryBudgets.
func (c *CategoryBudgets) upsertMultiple(w http.ResponseWriter, req *http.Request) {
	const (
		q_create = "INSERT INTO category_budgets (budget, category, amount) VALUES(?, ?, ?);"
		q_delete = "DELETE FROM category_budgets WHERE id = ?;"
	)

	var (
		cbs    []CategoryBudget
		budget int
		err    error
	)

	// Since we are doing multiple database operations, use a transaction
	tx, err := c.db.Begin()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Failed to begin database transaction: %s", err))
		return
	}
	defer tx.Rollback()

	json.NewDecoder(req.Body).Decode(&cbs)
	if len(cbs) == 0 {
		httperror.Send(w, req, http.StatusBadRequest, "Could not find any category budgets in request")
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Printf("body: %s\n", bodyString)
	log.Printf("cbs: %+v\n", cbs)

	// All cb's in a request must refer to the same budget. If the budget is not
	// provided, use the current one.
	if cbs[0].Budget == 0 {
		bgts, err := c.b.List()
		if err != nil {
			httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list budgets: %s", err))
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
			// TODO(davidschlachter): create a new budget if we cannot find an existing one for the current period
			httperror.Send(w, req, http.StatusInternalServerError, "Could not identify the current budget")
			return
		}
	} else {
		budget = cbs[0].Budget
		for _, cb := range cbs {
			if cb.Budget != budget {
				httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Got budget ID %d, expected %d. All category budgets in request must be for a single budget.", cb.Budget, budget))
				return
			}
		}
	}

	// upsertMultiple replaces all categoryBudgets for the budget. Delete before inserting.
	previous, err := c.List()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list previous category budgets: %s", err))
		return
	}
	for _, p := range previous {
		if p.Budget == budget {
			_, err = tx.Exec(q_delete, p.ID)
			if err != nil {
				httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not delete previous category budget: %s", err))
				return
			}
		}
	}

	// Ensure that no two category budgets share the same category
	cats := make(map[int]struct{})
	for _, cb := range cbs {
		_, ok := cats[cb.Category]
		if !ok {
			cats[cb.Category] = struct{}{}
		} else {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Got at least two category budgets for category ID %d, expected at most one.", cb.Category))
			return
		}
	}

	// Insert the new category budgets for the budget
	for _, cb := range cbs {
		if cb.Amount.IsZero() {
			continue // skip empty category budgets
		}
		_, err = tx.Exec(q_create, budget, cb.Category, cb.Amount)
		if err != nil {
			log.Printf("failed to upsert CategoryBudget: %s", err)
			httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not upsert categorybudget: %s", err))
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not commit changes to the database: %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
