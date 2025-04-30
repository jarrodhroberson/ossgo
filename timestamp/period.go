package timestamp

import (
	"fmt"
	"time"
)

// Period represents a time period with a start and end timestamp.
type Period struct {
	Start *Timestamp `json:"start"`
	End   *Timestamp `json:"end"`
}

// String returns a string representation of the Period in the format "start|end".
func (p *Period) String() string {
	return fmt.Sprintf("%s|%s", p.Start, p.End)
}

// Duration returns the duration of the Period.
// It calculates the absolute difference between the end and start timestamps.
func (p *Period) Duration() time.Duration {
	return p.End.Sub(p.Start).Abs()
}

// Contains checks if the given Timestamp is within the Period (inclusive of start and end).
func (p *Period) Contains(ts *Timestamp) bool {
	return (p.Start == ts || p.End == ts) || (p.Start.Before(ts) && p.End.After(ts))
}

// ToPeriod returns a Period that starts at from and ends at from + d.
// The end time is set to the beginning of the day.
func ToPeriod(from *Timestamp, d time.Duration) *Period {
	return &Period{
		Start: from,
		End:   from.Add(d).ZeroTime(),
	}
}

// Today returns a Period that represents the current day in UTC.
func Today() *Period {
	today := time.Now().UTC()
	return ToPeriod(From(today).ZeroTime(), time.Hour*24)
}

// DayToPeriod returns a Period that represents an entire day based on the given Timestamp.
// The start of the period corresponds to midnight at the beginning of the day, and
// the end corresponds to midnight at the beginning of the following day.
//
// For example, calling DayToPeriod(ts) where ts corresponds to "2023-10-27 15:00:00 UTC"
// will return a Period with:
// Start: "2023-10-27 00:00:00 UTC"
// End:   "2023-10-28 00:00:00 UTC"
func DayToPeriod(ts *Timestamp) (*Period) {
	return &Period{
		Start: ts.ZeroTime(),
		End:   ts.Add(time.Hour*24).ZeroTime(),
	}
}

// UntilNow returns a Period that begins at the "beginning of time" (a predefined earliest timestamp) 
// and ends at the current moment in time. This function is useful for representing a period that spans
// from the earliest point in time to the present.
func UntilNow() *Period {
	return &Period{
		Start: Enums().BeginningOfTime(),
		End:   Now(),
	}
}

// CurrentMonthPeriod returns a Period representing the current month.
// The start timestamp is set to the beginning of the first day of the current month,
// and the end timestamp corresponds to the beginning of the first day of the next month.
func CurrentMonthPeriod() *Period {
	now := time.Now()
	return MonthToPeriod(now.Year(), now.Month())
}

// YearToPeriod returns a Period that represents the entire year specified by the input year.
// The start of the period corresponds to midnight at the beginning of January 1st of the given year,
// and the end corresponds to midnight at the beginning of January 1st of the following year.
//
// For example, calling YearToPeriod(2023) will return a Period with:
// Start: 2023-01-01 00:00:00 UTC
// End:   2024-01-01 00:00:00 UTC
func YearToPeriod(year int) *Period {
	return &Period{
		Start: From(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)),
		End:   From(time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}
}

// YearToMonthlyPeriods returns a slice of Periods, each representing a month of the given year.
// Each Period will start at the beginning of the month and end at the beginning of the next month.
// For example, calling YearToMonthlyPeriods(2023) will produce 12 Periods:
// - Period 1: Start: 2023-01-01 00:00:00 UTC, End: 2023-02-01 00:00:00 UTC
// - Period 2: Start: 2023-02-01 00:00:00 UTC, End: 2023-03-01 00:00:00 UTC
// ...
// - Period 12: Start: 2023-12-01 00:00:00 UTC, End: 2024-01-01 00:00:00 UTC
//
// This function is useful for representing a year divided into monthly intervals.
func YearToMonthlyPeriods(year int) []*Period {
	var periods []*Period
	for month := time.January; month <= time.December; month++ {
		periods = append(periods, MonthToPeriod(year, month))
	}
	return periods
}

// MonthToPeriod returns a Period that represents the entire month of the given Timestamp.
func MonthToPeriod(year int, month time.Month) *Period {
	var end *Timestamp
	if month == time.December {
		end = From(time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC))
	} else {
		end = From(time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC))
	}
	return &Period{
		Start: From(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)),
		End:   end,
	}
}