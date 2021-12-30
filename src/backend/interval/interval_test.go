package interval_test

import (
	"testing"
	"time"

	"github.com/davidschlachter/lychnos/src/backend/interval"
)

const timeFormat = "2006-01-02 15:04:05 -0700"

var expectedIntervals = [][2]string{
	{"2020-01-01 00:00:00 -0500", "2020-01-31 23:59:59 -0500"},
	{"2020-02-01 00:00:00 -0500", "2020-02-29 23:59:59 -0500"},
	{"2020-03-01 00:00:00 -0500", "2020-03-31 23:59:59 -0400"},
	{"2020-04-01 00:00:00 -0400", "2020-04-30 23:59:59 -0400"},
	{"2020-05-01 00:00:00 -0400", "2020-05-31 23:59:59 -0400"},
	{"2020-06-01 00:00:00 -0400", "2020-06-30 23:59:59 -0400"},
	{"2020-07-01 00:00:00 -0400", "2020-07-31 23:59:59 -0400"},
	{"2020-08-01 00:00:00 -0400", "2020-08-31 23:59:59 -0400"},
	{"2020-09-01 00:00:00 -0400", "2020-09-30 23:59:59 -0400"},
	{"2020-10-01 00:00:00 -0400", "2020-10-31 23:59:59 -0400"},
	{"2020-11-01 00:00:00 -0400", "2020-11-30 23:59:59 -0500"},
	{"2020-12-01 00:00:00 -0500", "2020-12-31 23:59:59 -0500"},
}

func TestHandle(t *testing.T) {
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		t.Fatalf("Failed to load location 'America/Toronto': %s", err)
	}
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, location)
	end := time.Date(2020, 12, 30, 23, 59, 59, 0, location)
	intervals := interval.Get(start, end, location)

	if len(intervals) != len(expectedIntervals) {
		t.Fatalf("len(intervals) = %d, wanted %d\n", len(intervals), len(expectedIntervals))
	}

	for i := range intervals {
		expStart, _ := time.Parse(timeFormat, expectedIntervals[i][0])
		expEnd, _ := time.Parse(timeFormat, expectedIntervals[i][1])
		if !intervals[i].Start.Equal(expStart) {
			t.Fatalf("Start = %s, wanted %s\n", intervals[i].Start, expStart)
		}
		if !intervals[i].End.Equal(expEnd) {
			t.Fatalf("End = %s, wanted %s\n", intervals[i].End, expEnd)
		}
	}
}
