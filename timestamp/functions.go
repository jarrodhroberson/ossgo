package timestamp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var instances Timestamps
var once sync.Once

// daysIn returns the number of days in a given month and year.
func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// AddMonth returns the same day and clock time as t if possible,
// the day of the month of t does not exist m months from t
// the previous day is returned. ie: you request one month from October 31
// you would get November 30 and NOT December 1. This is OPPOSITE of the
// behavior of the standard library time.AddDate(), which would return December 1
func addMonth(t time.Time, m int) time.Time {
	x := t.AddDate(0, m, 0)
	if d := x.Day(); d != t.Day() {
		return x.AddDate(0, 0, -d)
	}
	return x
}

// AddMonth returns the same day and clock time as t if possible,
// the day of the month of t does not exist m months from t
// the previous day is returned. ie: you request one month from October 31
// you would get November 30 and NOT December 1. This is OPPOSITE of the
// behavior of the standard library time.AddDate(), which would return December 1
// m is the number of months to add
func AddMonth(ts *Timestamp, m int) *Timestamp {
	return From(addMonth(ts.t, m))
}

// ParseYouTubeTimestamp parses a non-standard timestamp format that the YouTube V3 API uses.
func ParseYouTubeTimestamp(s string) *Timestamp {
	return MustParse("2006-01-02T15:04:05.999999Z", s)
}

// FormatYouTubeActivityTimestamp formats a Timestamp using the YouTube timestamp format.
func FormatYouTubeActivityTimestamp(ts *Timestamp) string {
	return ts.t.Format("2006-01-02T15:04:05.0Z")
}

// IsZero returns true if the Timestamp is the zero value.
func IsZero(t *Timestamp) bool {
	return Enums().ZeroValue().Compare(t) == 0
}

// MustParse parses a timestamp from a string using the given format.
// It panics if the string cannot be parsed.
func MustParse(format string, s string) *Timestamp {
	// MustParse parses a timestamp from a string using the given format.
	// It panics if the string cannot be parsed.
	t, err := time.Parse(format, s)
	if err != nil {
		log.Error().Err(err).Msgf("Could not parse as timestamp %s", s)
	}
	return From(t)
}

// From creates a Timestamp from a time.Time.
func From(t time.Time) *Timestamp {
	return &Timestamp{
		t: t.UTC(),
	}
}

// FromMillis creates a Timestamp from milliseconds since the Unix epoch.
func FromMillis(ms int64) *Timestamp {
	return From(time.UnixMilli(ms))
}

// To converts a Timestamp to a time.Time.
func To(ts *Timestamp) time.Time {
	return ts.t
}

// Now returns the current time as a Timestamp.
func Now() *Timestamp {
	return &Timestamp{t: time.Now().UTC()}
}

// FormatCompact formats a Timestamp using the compact_format.
func FormatCompact(ts *Timestamp) string {
	return ts.t.Format(compact_format)
}

// ParseCompact parses a Timestamp from a string using the compact_format.
func ParseCompact(s string) *Timestamp {
	return MustParse(compact_format, s)
}

func Enums() Timestamps {
	once.Do(func() {
		instances = enums{
			beginningOfTime: &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},
			endOfTime:       &Timestamp{time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)},
			zeroValue:       From(time.Unix(0, 0).UTC()),
		}
	})
	return instances
}

// ToRange create an iterable slice of Timestamps with the interval d Duration
func ToRange(from *Timestamp, to *Timestamp, d time.Duration) []*Timestamp {
	interval := to.Sub(from) / d
	r := make([]*Timestamp, 0, interval)
	r = append(r, from)
	i := From(from.t.Add(d))
	for i.Before(to) {
		i = From(i.t.Add(d))
		r = append(r, i)
	}
	r = append(r, to)
	return r
}

// FirstDayOfNextMonth returns a Timestamp representing midnight (00:00:00) UTC on the first day
// of the next month. If the next month would be January, the year is incremented accordingly.
func FirstDayOfNextMonth() *Timestamp {
	nextMonth := AddMonth(Now(), 1).Month()
	year := Now().Year()
	if nextMonth == time.January {
		year = year + 1
	}
	return From(time.Date(year, nextMonth, 1, 0, 0, 0, 0, time.UTC))
}

// HumanReadableDuration converts a time.Duration into a human-readable string format.
// The output will include years, days, hours, minutes, and seconds as applicable,
// in descending order of significance. Components that are not relevant (i.e., with a value of 0)
// are omitted, except seconds, which is always included.
//
// For example:
// - A duration of 90061 seconds will yield "1d 1h 1m 1s"
// - A duration of 86400 seconds will yield "1d 0s"
// - A duration of 31556952 seconds will yield "1y 0s"
func HumanReadableDuration(d time.Duration) string {
	years := int(d.Hours()/24) / 365
	days := int(d.Hours()/24) - years*365
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	formattedDuration := ""
	if years > 0 {
		formattedDuration += fmt.Sprintf("%dy ", years)
	}
	if years > 0 || days > 0 {
		formattedDuration += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 || days > 0 { // Include hours if days are present for clarity
		formattedDuration += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 || days > 0 { // Include minutes similarly
		formattedDuration += fmt.Sprintf("%dm ", minutes)
	}
	formattedDuration += fmt.Sprintf("%ds", seconds)

	return formattedDuration
}
