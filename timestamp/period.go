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

// MonthToPeriod returns a Period that represents the entire month of the given Timestamp.
func MonthToPeriod(ts *Timestamp) *Period {
	firstDayOfMonth := From(time.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0, time.UTC))
	daysInMonth := daysIn(ts.Month(), ts.Year())
	lastDayOfMonth := From(time.Date(ts.Year(), ts.Month(), daysInMonth, 0, 0, 0, 0, time.UTC).Add(24 * time.Hour))
	return &Period{
		Start: firstDayOfMonth,
		End:   lastDayOfMonth,
	}
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

func UntilNow() *Period {
	return &Period{
		Start: Enums().BeginningOfTime(),
		End:   Now(),
	}
}
