package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func create(w http.ResponseWriter, req *http.Request) {
	const (
		q      = "INSERT INTO budgets VALUES(?, ?, ?, ?)"
		format = "2006-01-02 15:04:05"
	)

	var (
		err      error
		interval int
	)

	err = req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not parse POST data\n")
		return
	}

	// Get parameters
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
	_, err = db.Exec(q, nil, start.Format(format), end.Format(format), interval)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not insert budget: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type Budget struct {
	ID                int       `json:"id"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
	ReportingInterval int       `json:"reporting_interval"`
}

func list(w http.ResponseWriter, req *http.Request) {
	const q = "SELECT id, start, end, reporting_interval FROM budgets;"
	rows, err := db.Query(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not list budgets: %s", err)
		return
	}
	defer rows.Close()

	var budgets []Budget

	for rows.Next() {
		var b Budget
		rows.Scan(&b.ID, &b.Start, &b.End, &b.ReportingInterval)
		budgets = append(budgets, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func handleBudget(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		list(w, req)
	case "POST":
		create(w, req)
	case "PATCH":
		w.WriteHeader(http.StatusNotImplemented)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported method %s", req.Method)
	}
}
