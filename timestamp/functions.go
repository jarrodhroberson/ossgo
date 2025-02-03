package timestamp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var instances Timestamps
var once sync.Once

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
func AddMonth(ts Timestamp, m int) Timestamp {
	return From(addMonth(ts.t, m))
}

// ParseYouTubeTimestamp parses a non-standard timestamp format that the YouTube V3 API uses.
func ParseYouTubeTimestamp(s string) Timestamp {
	return MustParse("2006-01-02T15:04:05.999999Z", s)
}

func IsZero(t Timestamp) bool {
	return Enums().ZeroValue().Compare(t) == 0
}

func MonthToPeriod(ts Timestamp) Period {
	firstDayOfMonth := From(time.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0, time.UTC))
	duration := time.Duration(24 * 7 * daysIn(ts.Month(), ts.Year()))
	return ToPeriod(firstDayOfMonth, duration)
}

func ToPeriod(from Timestamp, d time.Duration) Period {
	return Period{
		Start: from,
		End:   from.Add(d),
	}
}

func Today() Period {
	today := time.Now().UTC()
	return ToPeriod(From(today).ZeroTime(), time.Hour*24)
}

func MustParse(format string, s string) Timestamp {
	t, err := time.Parse(format, s)
	if err != nil {
		log.Error().Err(err).Msgf("Could not parse as timestamp %s", s)
	}
	return From(t)
}

func From(t time.Time) Timestamp {
	return Timestamp{
		t: t.UTC(),
	}
}

func FromMillis(ms int64) Timestamp {
	return From(time.UnixMilli(ms))
}

func To(ts Timestamp) time.Time {
	return ts.t
}

func Now() Timestamp {
	return Timestamp{t: time.Now().UTC()}
}

func Enums() Timestamps {
	once.Do(func() {
		instances = tsenums{}
	})
	return instances
}

// ToRange create an iterable slice of Timestamps with the interval d Duration
func ToRange(from Timestamp, to Timestamp, d time.Duration) []Timestamp {
	interval := to.Sub(from) / d
	r := make([]Timestamp, 0, interval)
	r = append(r, from)
	i := From(from.t.Add(d))
	for i.Before(to) {
		i = From(i.t.Add(d))
		r = append(r, i)
	}
	r = append(r, to)
	return r
}
