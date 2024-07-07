package timestamp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var instances Timestamps
var once sync.Once

func NextMonth() time.Month {
	today := time.Now().UTC()
	month := today.Month()
	if month == time.December {
		return time.January
	}
	return month + 1
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

func ToRange(from Timestamp, to Timestamp, d time.Duration) []Timestamp {
	r := make([]Timestamp, 0, d)
	r = append(r, from)
	i := From(from.t.Add(d))
	for i.Before(to) {
		r = append(r, i)
		i = From(i.t.Add(d))
	}
	r = append(r, to)
	return r
}
