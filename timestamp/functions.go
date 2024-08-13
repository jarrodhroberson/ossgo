package timestamp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var instances Timestamps
var once sync.Once

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

func ToRange(from Timestamp, to Timestamp, interval int64, d time.Duration) []Timestamp {
	r := make([]Timestamp, 0, interval)
	r = append(r, from)
	id := time.Duration(interval * int64(d))
	i := From(from.t.Add(id))
	for !i.Before(to) {
		i = From(i.t.Add(id))
		r = append(r, i)
	}
	r = append(r, to)
	return r
}
