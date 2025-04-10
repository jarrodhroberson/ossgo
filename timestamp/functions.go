package timestamp

import (
	"math"
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

// FormatCompact formats a Timestamp using the CompactFormat.
func FormatCompact(ts *Timestamp) string {
	return ts.t.Format(CompactFormat)
}

// ParseCompact parses a Timestamp from a string using the CompactFormat.
func ParseCompact(s string) *Timestamp {
	return MustParse(CompactFormat, s)
}

func Enums() Timestamps {
	once.Do(func() {
		instances = enums{
			beginningOfTime: &Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)},
			endOfTime:       &Timestamp{time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)},
			zeroValue:       From(time.UnixMilli(math.MinInt64).UTC()),
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
