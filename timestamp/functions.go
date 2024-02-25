package timestamp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var instances Timestamps
var once sync.Once

func MonthFromToday() Period {
	today := time.Now().UTC()
	return Period{
		Start: From(today.AddDate(0, 0, -today.Day()+1)),
		End:   From(today.AddDate(0, 1, -today.Day())),
	}
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

func ToRange(from Timestamp, to Timestamp, d time.Duration) []Timestamp {
	r := make([]Timestamp, 0, 12)
	r = append(r, from)
	i := From(from.t.Add(d))
	for i.Before(to) {
		r = append(r, i)
		i = From(i.t.Add(d))
	}
	r = append(r, to)
	return r
}
