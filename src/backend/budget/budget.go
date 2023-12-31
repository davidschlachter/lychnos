package budget

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/httperror"
)

const dateFormat = "2006-01-02 15:04:05"

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
	log.Printf("%s %s", req.Method, req.RequestURI)
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
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse budget ID: %s\n", id))
		return
	}

	budgets, err := b.Fetch(id)
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not fetch budget: %s", err))
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
		httperror.Send(w, req, http.StatusInternalServerError, fmt.Sprintf("Could not list budgets: %s", err))
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
	var (
		err      error
		interval int
		id       int
	)

	err = req.ParseForm()
	if err != nil {
		httperror.Send(w, req, http.StatusInternalServerError, "Could not parse POST data")
		return
	}

	// Get parameters
	idStr := req.Form.Get("id")
	if len(idStr) == 0 {
		id = 0
	} else {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse ID: %s", idStr))
			return
		}
	}
	startString := req.Form.Get("start")
	endString := req.Form.Get("end")
	if len(startString) == 0 || len(endString) == 0 {
		httperror.Send(w, req, http.StatusBadRequest, "Must provide start and end datetimes")
		return
	}
	intervalString := req.Form.Get("interval")
	if len(intervalString) == 0 {
		interval = 0
	} else {
		interval, err = strconv.Atoi(intervalString)
		if err != nil {
			httperror.Send(w, req, http.StatusBadRequest, "Interval must be an integer")
			return
		}
	}

	// Validate dates
	// TODO(davidschlachter): if the client doesn't provide a time zone, these
	// will be created in UTC, which will lead to unexpected behaviour.
	start, err := time.Parse(dateFormat, startString)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse start time: %s", err))
		return
	}
	end, err := time.Parse(dateFormat, endString)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse end time: %s", err))
		return
	}
	if start.After(end) {
		httperror.Send(w, req, http.StatusBadRequest, "start must be before end")
		return
	}

	// TODO(davidschlachter): check that the new budget does not overlap the
	// date ranges of any other one

	// Insert the budget into the database
	err = b.Upsert(id, start, end, interval)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not upsert budget: %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (b *Budgets) Upsert(id int, start, end time.Time, interval int) error {
	const q = "REPLACE INTO budgets (id, start, end, reporting_interval) VALUES(?, ?, ?, ?);"

	// Insert the budget into the database
	_, err := b.db.Exec(q, id, start.UTC().Format(dateFormat), end.UTC().Format(dateFormat), interval)
	return err
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
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not parse ID: %s", idStr))
		return
	}

	_, err = b.db.Exec(q, id)
	if err != nil {
		httperror.Send(w, req, http.StatusBadRequest, fmt.Sprintf("Could not delete budget: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
