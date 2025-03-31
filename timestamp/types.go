package timestamp

import (
	"strconv"
	"time"
)

const CompactFormat = "20060102t150405Z"

type Timestamps interface {
	BeginningOfTime() *Timestamp
	EndOfTime() *Timestamp
	ZeroValue() *Timestamp
}

type enums struct {
	beginningOfTime *Timestamp
	endOfTime       *Timestamp
	zeroValue       *Timestamp
}

func (i enums) BeginningOfTime() *Timestamp {
	return i.beginningOfTime
}
func (i enums) EndOfTime() *Timestamp {
	return i.endOfTime
}
func (i enums) ZeroValue() *Timestamp {
	return i.zeroValue
}

type Timestamp struct {
	t time.Time
}

func (ts *Timestamp) Year() int {
	return ts.t.Year()
}

func (ts *Timestamp) Month() time.Month {
	return ts.t.Month()
}

func (ts *Timestamp) Day() time.Weekday {
	return ts.t.Weekday()
}

func (ts *Timestamp) Before(ots *Timestamp) bool {
	return ts.t.Before(ots.t)
}

func (ts *Timestamp) After(ots *Timestamp) bool {
	return ts.t.After(ots.t)
}

func (ts *Timestamp) Add(d time.Duration) *Timestamp {
	return From(ts.t.Add(d))
}

func (ts *Timestamp) Sub(ots *Timestamp) time.Duration {
	return ts.t.Sub(ots.t)
}

func (ts *Timestamp) ZeroTime() *Timestamp {
	return From(ts.t.Truncate(time.Hour * 24))
}

func (ts *Timestamp) Compare(to *Timestamp) int {
	return ts.t.Compare(to.t)
}

func (ts *Timestamp) String() string {
	bytes, _ := ts.MarshalText()
	return string(bytes)
}

func (ts *Timestamp) MarshalText() (text []byte, err error) {
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

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
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

func (ts *Timestamp) MarshalBinary() (data []byte, err error) {
	return ts.MarshalText()
}

func (ts *Timestamp) UnmarshalBinary(b []byte) error {
	return ts.UnmarshalText(b)
}