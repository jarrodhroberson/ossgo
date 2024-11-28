package events

import (
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

type Events struct {
	timestamp.Period
	Events []Event `json:"events"`
}

type Event interface {
	Id() string
	MessageType() string
	Message() string
	ReceivedAt() timestamp.Timestamp
}

type event struct {
	id              string
	messageMimeType string
	message         string
	receivedAt      timestamp.Timestamp
}

func (e event) Id() string {
	return e.id
}

func (e event) MessageType() string {
	return e.messageMimeType
}

func (e event) Message() string {
	return e.message
}

func (e event) ReceivedAt() timestamp.Timestamp {
	return e.receivedAt
}

func (e event) String() string {
	return string(must.MarshalJson(e))
}
