package timeseries

import (
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

type Event struct {
	Id        string              `json:"id"`
	Details   map[string]string   `json:"details"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
}

func (e Event) String() string {
	return string(must.MarshalJson(must.MarshallMap(e)))
}

type TimeSeries struct {
	Id     string  `json:"id"`
	Events []Event `json:"events"`
}

func (t TimeSeries) String() string {
	return string(must.MarshalJson(must.MarshallMap(t)))
}

type Namespace struct {
	Id         string       `json:"id"`
	TimeSeries []TimeSeries `json:"time_series"`
}

func (n Namespace) String() string {
	return string(must.MarshalJson(must.MarshallMap(n)))
}
