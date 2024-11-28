package timeseries

import (
	"github.com/jarrodhroberson/ossgo/timestamp"
)

type Detail struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Event struct {
	Id        string              `json:"id"`
	Details   []Detail            `json:"details"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
}

type TimeSeries struct {
	Id     string  `json:"id"`
	Events []Event `json:"events"`
}

type Namespace struct {
	Id         string       `json:"id"`
	TimeSeries []TimeSeries `json:"time_series"`
}
