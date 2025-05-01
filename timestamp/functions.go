package timestamp

import (
	"fmt"
	"strings"
	"sync"
	"time"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
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

func Format(ts *Timestamp, format string) string {
	return ts.t.Format(format)
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
		log.Error().Stack().Err(err).Msgf("Could not parse as timestamp %s", s)
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

// FormatCompact formats a Timestamp using the COMPACT_FORMAT.
func FormatCompact(ts *Timestamp) string {
	return ts.t.Format(COMPACT_FORMAT)
}

// ParseCompact parses a Timestamp from a string using the COMPACT_FORMAT.
func ParseCompact(s string) *Timestamp {
	return MustParse(COMPACT_FORMAT, s)
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

// ISO8601ToDuration converts an ISO 8601 duration string into a time.Duration.
// It parses the duration components (years, months, weeks, days, hours, minutes, seconds),
// sums them up, and returns the total duration.
//
// Supported format examples:
//   - "P1Y2M3W4DT5H6M7S" (1 year, 2 months, 3 weeks, 4 days, 5 hours, 6 minutes, 7 seconds)
//   - "PT10M" (10 minutes)
//
// Note: For simplicity, this implementation assumes all months are 30 days
// and all years are 365 days. This can result in approximations when handling
// months and years. It also only returns INTEGER seconds. Any decimal places are stripped off.
//
// Parameters:
//
//	s - An ISO 8601 duration string.
//
// Returns:
//
//	A time.Duration representing the total duration, or an error if the input string
//	is not a valid ISO 8601 duration format.
//
// Errors:
//
//	Returns an error if the input string does not match the ISO 8601 format.
func ISO8601ToDuration(s string) (time.Duration, error) {
	if s == "" {
		return -1, errs.MustNotBeEmpty.New("duration string is empty")
	}
	if len(s) < 3 {
		return -1, errs.MinSizeExceededError.New("duration string must be >= 4 characters")
	}
	if !iso8601DurationRegex.MatchString(s) {
		err := errs.InvalidFormat.New("invalid ISO 8601 duration: %s", iso8601DurationRegex.String())
		return -1, errs.ParseError.Wrap(err, "invalid ISO 8601 duration: %s", s)
	}
	log.Info().Msgf("duration string: %s", s)
	matches := iso8601DurationRegex.FindStringSubmatch(s)
	names := iso8601DurationRegex.SubexpNames()
	groups := make(map[string]string, len(names))
	for idx, name := range names {
		groups[name] = matches[idx]
	}

	var dur int64
	if val, ok := groups["years"]; ok {
		if len(val) > 0 {
			years := must.ParseInt(val)
			dur += int64(years) * int64(time.Hour*24*365)
		}
	}
	if val, ok := groups["months"]; ok {
		if len(val) > 0 {
			months := must.ParseInt(val)
			dur += int64(months) * int64(time.Hour*24*30)
		}
	}
	if val, ok := groups["weeks"]; ok {
		if len(val) > 0 {
			weeks := must.ParseInt(val)
			dur += int64(weeks) * int64(time.Hour*24*7)
		}
	}
	if val, ok := groups["days"]; ok {
		if len(val) > 0 {
			days := must.ParseInt(val)
			dur += int64(days) * int64(time.Hour*24)
		}
	}
	if val, ok := groups["time"]; ok {
		log.Info().Msgf("time: %s", val)
		if len(val) != 0 {
			if val, ok := groups["hours"]; ok {
				if len(val) > 0 {
					hours := must.ParseInt(val)
					dur += int64(hours) * int64(time.Hour)
				}
			}
			if val, ok := groups["minutes"]; ok {
				if len(val) > 0 {
					minutes := must.ParseInt(val)
					dur += int64(minutes) * int64(time.Minute)
				}
			}
			if val, ok := groups["seconds"]; ok {
				if len(val) > 0 {
					decIdx := strings.Index(val, ".")
					if decIdx > -1 {
						val = val[:decIdx]
					}
					seconds := must.ParseInt(val)
					dur += int64(seconds) * int64(time.Second)
				}
			}
		}
	}

	return time.Duration(dur), nil
}
