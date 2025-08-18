package timestamp

import (
	"encoding/json"
	"fmt"
	"time"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
)

func validate(p *Period) (*Period, error) {
	if p.start == nil {
		return nil, errs.MustNotBeNil.New("start")
	}
	if p.end == nil {
		return nil, errs.MustNotBeNil.New("end")
	}
	if p.start.After(p.end) {
		return nil, errs.InvalidState.New("start timestamp must be before end timestamp %s:%s", p.Start(), p.End())
	}
	if p.start.t.Equal(p.end.t) {
		return nil, errs.InvalidState.New("start and end timestamps must not be equal %s:%s", p.Start(), p.End())
	}
	return p, nil
}

func NewPeriod(start *Timestamp, end *Timestamp) *Period {
	return must.Must(validate(&Period{
		start: start,
		end:   end,
	}))
}

// period this exists to provide json marshaling and unmarshalling proxy
type period struct {
	Start *Timestamp `json:"start"`
	End   *Timestamp `json:"end"`
}

// Period represents a time period with a start and end timestamp.
type Period struct {
	start *Timestamp
	end   *Timestamp
}

func (p *Period) Start() *Timestamp {
	return p.start
}

func (p *Period) End() *Timestamp {
	return p.end
}

func (p *Period) MarshalJSON() ([]byte, error) {
	jsonStruct := period{
		Start: p.start,
		End:   p.end,
	}
	return json.Marshal(jsonStruct)
}

func (p *Period) UnmarshalJSON(data []byte) error {
	var jsonStruct period
	if err := json.Unmarshal(data, &jsonStruct); err != nil {
		return err
	}
	p.start = jsonStruct.Start
	p.end = jsonStruct.End
	return nil
}

// String returns a string representation of the Period in the format "start|end".
func (p *Period) String() string {
	return fmt.Sprintf("%s|%s", p.start, p.end)
}

// Duration returns the duration of the Period.
// It calculates the absolute difference between the end and start timestamps.
func (p *Period) Duration() time.Duration {
	return p.end.Sub(p.start).Abs()
}

// Contains checks if the given Timestamp is within the Period (inclusive of start and end).
func (p *Period) Contains(ts *Timestamp) bool {
	return (p.start == ts || p.end == ts) || (p.start.Before(ts) && p.end.After(ts))
}

// ToPeriod returns a Period that starts at from and ends at from + d.
// The end time is set to the beginning of the day.
func ToPeriod(from *Timestamp, d time.Duration) *Period {
	return &Period{
		start: from,
		end:   from.Add(d).ZeroTime(),
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
func DayToPeriod(ts *Timestamp) *Period {
	return must.Must(validate(&Period{
		start: ts.ZeroTime(),
		end:   ts.Add(time.Hour * 24).ZeroTime(),
	}))
}

// UntilNow returns a Period that begins at the "beginning of time" (a predefined earliest timestamp)
// and ends at the current moment in time. This function is useful for representing a period that spans
// from the earliest point in time to the present.
func UntilNow() *Period {
	return must.Must(validate(&Period{
		start: Enums().BeginningOfTime(),
		end:   Now(),
	}))
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
	return must.Must(validate(&Period{
		start: From(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)),
		end:   From(time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}))
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

// MonthToPeriod returns a Period representing an entire month for the given year and month.
// The start timestamp corresponds to midnight at the beginning of the first day of the specified month,
// and the end timestamp corresponds to midnight at the beginning of the first day of the next month.
//
// The function automatically handles leap years through Go's time.Date function, which correctly
// calculates dates considering leap years. This means February will automatically have 29 days
// in leap years and 28 days in non-leap years without requiring explicit handling.
//
// For example:
// MonthToPeriod(2024, time.February) returns a Period with:
// Start: 2024-02-01 00:00:00 UTC
// End:   2024-03-01 00:00:00 UTC
//
// Parameters:
//   - year: The year as an integer (e.g., 2024)
//   - month: The month as time.Month (e.g., time.February)
//
// Returns a Period with the start and end timestamps for the specified month.
func MonthToPeriod(year int, month time.Month) *Period {
	var end *Timestamp
	if month == time.December {
		end = From(time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC))
	} else {
		nextMonth := month + 1
		nextYear := year
		if nextMonth > time.December {
			nextMonth = time.January
			nextYear++
		}
		end = From(time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, time.UTC))
	}
	return &Period{
		start: From(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)),
		end:   end,
	}
}
