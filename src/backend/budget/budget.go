package budget

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Budget struct {
	ID                int       `json:"id"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
	ReportingInterval int       `json:"reporting_interval"`
}

type Budgets struct {
	db *sql.DB
}

func New(db *sql.DB) *Budgets {
	return &Budgets{db: db}
}

func (b *Budgets) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		hasID := regexp.MustCompile(`/[0-9]+$`)
		if hasID.MatchString(req.URL.Path) {
			b.fetch(w, req)
		} else {
			b.list(w, req)
		}
	case "POST":
		b.upsert(w, req)
	case "DELETE":
		b.delete(w, req)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}

func (b *Budgets) fetch(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse budget ID: %s\n", id)
		return
	}

	budgets, err := b.Fetch(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not fetch budget: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func (b *Budgets) Fetch(id string) ([]Budget, error) {
	const q = "SELECT id, start, end, reporting_interval FROM budgets WHERE id = ?;"

	row := b.db.QueryRow(q, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var budgets []Budget

	var bgt Budget
	row.Scan(&bgt.ID, &bgt.Start, &bgt.End, &bgt.ReportingInterval)
	budgets = append(budgets, bgt)

	return budgets, nil
}

func (b *Budgets) list(w http.ResponseWriter, req *http.Request) {
	budgets, err := b.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not list budgets: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func (b *Budgets) List() ([]Budget, error) {
	const q = "SELECT id, start, end, reporting_interval FROM budgets;"
	rows, err := b.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget

	for rows.Next() {
		var bgt Budget
		rows.Scan(&bgt.ID, &bgt.Start, &bgt.End, &bgt.ReportingInterval)
		budgets = append(budgets, bgt)
	}

	return budgets, nil
}

func (b *Budgets) upsert(w http.ResponseWriter, req *http.Request) {
	const (
		q      = "INSERT INTO budgets VALUES(?, ?, ?, ?)"
		format = "2006-01-02 15:04:05"
	)

	var (
		err      error
		interval int
		id       int
	)

	err = req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not parse POST data\n")
		return
	}

	// Get parameters
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
	startString := req.Form.Get("start")
	endString := req.Form.Get("end")
	if len(startString) == 0 || len(endString) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Must provide start and end datetimes\n")
		return
	}
	intervalString := req.Form.Get("interval")
	if len(intervalString) == 0 {
		interval = 0
	} else {
		interval, err = strconv.Atoi(intervalString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Interval must be an integer\n")
			return
		}
	}

	// Validate dates
	start, err := time.Parse(format, startString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse start time: %s", err)
		return
	}
	end, err := time.Parse(format, endString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse end time: %s", err)
		return
	}
	if start.After(end) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "start must be before end\n")
		return
	}

	// TODO(davidschlachter): check that the new budget does not overlap the
	// date ranges of any other one

	// Insert the budget into the database
	_, err = b.db.Exec(q, id, start.Format(format), end.Format(format), interval)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not upsert budget: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (b *Budgets) delete(w http.ResponseWriter, req *http.Request) {
	const q = "DELETE FROM budgets WHERE id = ?;"

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

	_, err = b.db.Exec(q, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not delete budget: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
