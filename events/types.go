package events

import (
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

type Events[T any] struct {
	timestamp.Period
	Events []Event[T] `json:"events"`
}

type Event[T any] interface {
	Id() string
	MessageType() string
	Message() *T
	ReceivedAt() timestamp.Timestamp
}

type event[T any] struct {
	id              string
	messageMimeType string
	message         *T
	receivedAt      timestamp.Timestamp
}

func (e event[T]) Id() string {
	return e.id
}

func (e event[T]) MessageType() string {
	return e.messageMimeType
}

func (e event[T]) Message() *T {
	return e.message
}

func (e event[T]) ReceivedAt() timestamp.Timestamp {
	return e.receivedAt
}

func (e event[T]) String() string {
	return string(must.MarshalJson(e))
}
