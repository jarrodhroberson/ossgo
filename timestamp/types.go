package timestamp

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Timestamps interface {
	BeginningOfTime() Timestamp
	EndOfTime() Timestamp
	ZeroValue() Timestamp
}

type tsenums struct {
}

func (i tsenums) BeginningOfTime() Timestamp {
	return Timestamp{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}
}
func (i tsenums) EndOfTime() Timestamp {
	return Timestamp{time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)}
}

func (i tsenums) ZeroValue() Timestamp {
	return From(time.UnixMilli(math.MinInt64))
}

type Timestamp struct {
	t time.Time
}

func (ts Timestamp) Year() int {
	return ts.t.Year()
}

func (ts Timestamp) Month() time.Month {
	return ts.t.Month()
}

func (ts Timestamp) Day() time.Weekday {
	return ts.t.Weekday()
}

func (ts Timestamp) Before(ots Timestamp) bool {
	return ts.t.Before(ots.t)
}

func (ts Timestamp) After(ots Timestamp) bool {
	return ts.t.After(ots.t)
}

func (ts Timestamp) Add(d time.Duration) Timestamp {
	return From(ts.t.Add(d))
}

func (ts Timestamp) Sub(ots Timestamp) time.Duration {
	return ts.t.Sub(ots.t)
}

func (ts Timestamp) ZeroTime() Timestamp {
	return From(ts.t.Truncate(time.Hour * 24))
}

func (ts Timestamp) Compare(to Timestamp) int {
	return ts.t.Compare(to.t)
}

func (ts Timestamp) String() string {
	bytes, _ := ts.MarshalText()
	return string(bytes)
}

func (ts Timestamp) MarshalText() (text []byte, err error) {
	return []byte(ts.t.UTC().Format(time.RFC3339Nano)), nil
}

func (ts *Timestamp) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return err
	}
	ts.t = t.UTC()
	return nil
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	bytes, err := ts.MarshalText()
	return []byte(strconv.Quote(string(bytes))), err
}

func (ts *Timestamp) UnmarshallJSON(b []byte) error {
	t, err := time.Parse(time.RFC3339Nano, string(b))
	if err != nil {
		return err
	}
	ts.t = t.UTC()
	return nil
}

func (ts Timestamp) MarshalBinary() (data []byte, err error) {
	return ts.MarshalText()
}

func (ts *Timestamp) UnmarshalBinary(b []byte) error {
	return ts.UnmarshalText(b)
}

type Period struct {
	Start Timestamp `json:"start"`
	End   Timestamp `json:"end"`
}

func (p Period) String() string {
	return fmt.Sprintf("%s|%s", p.Start, p.End)
}

func (p Period) Duration() time.Duration {
	return p.End.Sub(p.Start).Abs()
}

func (p Period) Contains(ts Timestamp) bool {
	return (p.Start == ts || p.End == ts) || (p.Start.Before(ts) && p.End.After(ts))
}
