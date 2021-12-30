// Package interval provides functions for extracting reporting_intervals from
// the stand and end date of a budget.
package interval

import "time"

type ReportingInterval struct {
	Start, End time.Time
}

// Get returns a ReportingInterval slice, which contains the Start and End of
// each reporting interval between the provided start and end time
func Get(start, end time.Time, location *time.Location) []ReportingInterval {
	var intervals []ReportingInterval

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	lastDayThisMonth := time.Date(currentYear, currentMonth, daysInMonth(int(currentMonth), currentYear), 23, 59, 59, 0, location)

	var bgtYear int
	var bgtMonth time.Month

	// d is the first day in each summary reporting period
	for d := start; d.Before(end) && d.Before(lastDayThisMonth); d = time.Date(bgtYear, bgtMonth+1, 1, 0, 0, 0, 0, location) {
		bgtYear, bgtMonth, _ = d.Date()
		// l is the last day of each reporting period
		l := time.Date(bgtYear, bgtMonth, daysInMonth(int(bgtMonth), bgtYear), 23, 59, 59, 0, location)

		r := ReportingInterval{Start: d, End: l}
		intervals = append(intervals, r)
	}

	return intervals
}

// via https://stackoverflow.com/a/35182930
func daysInMonth(month, year int) int {
	switch time.Month(month) {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) { // leap year
			return 29
		}
		return 28
	default:
		return 31
	}
}
