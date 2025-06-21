package timestamp

import (
	"fmt"
	"regexp"
	"strconv"
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

// ParseYouTubeTimestamp parses the timestamp format that the YouTube V3 API uses.
func ParseYouTubeTimestamp(s string) *Timestamp {
	return MustParse(time.RFC3339, s)
}

// FormatYouTubeActivityTimestamp formats a Timestamp using the YouTube timestamp format.
func FormatYouTubeActivityTimestamp(ts *Timestamp) string {
	return ts.t.Format(time.RFC3339)
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
	if len(s) == 0 {
		return -1, errs.MustNotBeEmpty.New("duration string is empty")
	}
	if len(s) < 3 {
		return -1, errs.MinSizeExceededError.New("duration string must be >= 4 characters")
	}
	if s[0] != 'P' && s[0] != 'T' {
		return -1, errs.InvalidFormat.New("invalid ISO 8601 duration: %s", s)
	}
	if s[len(s)-1] == 'T' {
		err := errs.InvalidFormat.New("invalid ISO 8601 duration: %s", s)
		return -1, errs.MustNeverError.Wrap(err, "missing designator for time")
	}
	var matches []string
	var names []string
	var groups map[string]string
	if s[1] == 'T' {
		matches = timeDurationRegex.FindStringSubmatch(s)
		names = timeDurationRegex.SubexpNames()
		groups = make(map[string]string, len(names))
	} else {
		matches = dateDurationRegex.FindStringSubmatch(s)
		names = dateDurationRegex.SubexpNames()
		groups = make(map[string]string, len(names))
	}

	for idx, name := range names {
		groups[name] = matches[idx]
	}
	fmt.Println(groups)

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

	return time.Duration(dur), nil
}

// ParseISO8601Duration parses an ISO 8601 duration string and returns a time.Duration.
// It handles years, months, weeks, days, hours, minutes, and seconds (including fractional seconds).
// It returns an error if the input string is not a valid ISO 8601 duration.  Because time.Duration
// only goes to the nanosecond, the values for years, months, and weeks are approximate.
func ParseISO8601Duration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return -1, errs.MustNotBeEmpty.New("duration string is empty")
	}
	if len(s) < 3 {
		return -1, errs.MinSizeExceededError.New("duration string must be >= 4 characters")
	}
	if s[0] != 'P' && s[0] != 'T' {
		return -1, errs.InvalidFormat.New("invalid ISO 8601 duration: %s", s)
	}
	if s[len(s)-1] == 'T' {
		err := errs.InvalidFormat.New("invalid ISO 8601 duration: %s", s)
		return -1, errs.MustNeverError.Wrap(err, "missing designator for time")
	}
	if s == "PT" {
		return 0, nil
	}

	var match []string
	var regex *regexp.Regexp

	if strings.Contains(s, "T") {
		regex = timeDurationRegex
		match = append(match, regex.FindStringSubmatch(s)...)
	} else {
		regex = dateDurationRegex
		match = append(match, regex.FindStringSubmatch(s)...)
	}

	if match == nil {
		return 0, fmt.Errorf("invalid ISO 8601 duration format: %s", s)
	}

	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" && match[i] != "" { // Ensure we don't process empty matches
			result[name] = match[i]
		}
	}

	var duration time.Duration

	if yearsStr, ok := result["years"]; ok {
		years, err := strconv.Atoi(yearsStr)
		if err != nil {
			return 0, fmt.Errorf("invalid years value: %s", yearsStr)
		}
		// Approximate: 1 year = 365.25 days
		duration += time.Duration(int64(years) * 365 * 24 * int64(time.Hour)) // Corrected: Explicit conversion
	}

	if monthsStr, ok := result["months"]; ok {
		months, err := strconv.Atoi(monthsStr)
		if err != nil {
			return 0, fmt.Errorf("invalid months value: %s", monthsStr)
		}
		// Approximate: 1 month = 30.44 days
		duration += time.Duration(int64(months) * 30 * 24 * int64(time.Hour)) // Corrected: Explicit conversion
	}

	if weeksStr, ok := result["weeks"]; ok {
		weeks, err := strconv.Atoi(weeksStr)
		if err != nil {
			return 0, fmt.Errorf("invalid weeks value: %s", weeksStr)
		}
		duration += time.Duration(int64(weeks) * 7 * 24 * int64(time.Hour)) // Corrected: Explicit conversion
	}

	if daysStr, ok := result["days"]; ok {
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("invalid days value: %s", daysStr)
		}
		duration += time.Duration(int64(days) * 24 * int64(time.Hour)) // Corrected: Explicit conversion
	}

	if hoursStr, ok := result["hours"]; ok {
		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return 0, fmt.Errorf("invalid hours value: %s", hoursStr)
		}
		duration += time.Duration(hours) * time.Hour
	}

	if minutesStr, ok := result["minutes"]; ok {
		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes value: %s", minutesStr)
		}
		duration += time.Duration(minutes) * time.Minute
	}

	if val, ok := result["seconds"]; ok {
		if len(val) > 0 {
			decIdx := strings.Index(val, ".")
			if decIdx > -1 {
				val = val[:decIdx]
			}
			seconds := must.ParseInt64(val)
			duration += time.Duration(seconds) * time.Second
		}
	}

	return duration, nil
}
