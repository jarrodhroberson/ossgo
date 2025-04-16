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

func (e *event[T]) Id() string {
	return e.id
}

func (e *event[T]) MessageType() string {
	return e.messageMimeType
}

func (e *event[T]) Message() *T {
	return e.message
}

func (e *event[T]) ReceivedAt() timestamp.Timestamp {
	return e.receivedAt
}

func (e *event[T]) String() string {
	return string(must.MarshalJson(e))
}

func (e *event[T]) UnmarshalJSON(bytes []byte) error {
	var data struct {
		ID          string              `json:"id"`
		MessageType string              `json:"messageType"`
		Message     *T                  `json:"message"`
		ReceivedAt  timestamp.Timestamp `json:"receivedAt"`
	}

	must.UnMarshalJson(bytes, &data)

	e.id = data.ID
	e.messageMimeType = data.MessageType
	e.message = data.Message
	e.receivedAt = data.ReceivedAt

	return nil
}

func (e *event[T]) MarshalJSON() ([]byte, error) {
	return must.MarshalJson(&struct {
		ID          string              `json:"id"`
		MessageType string              `json:"messageType"`
		Message     *T                  `json:"message"`
		ReceivedAt  timestamp.Timestamp `json:"receivedAt"`
	}{
		ID:          e.id,
		MessageType: e.messageMimeType,
		Message:     e.message,
		ReceivedAt:  e.receivedAt,
	}), nil
}

func (e *event[T]) MarshalBinary() ([]byte, error) {
	return e.MarshalJSON()
}

func (e *event[T]) UnmarshalBinary(data []byte) error {
	return e.UnmarshalJSON(data)
}

func (e *event[T]) MarshalText() ([]byte, error) {
	return e.MarshalJSON()
}

func (e *event[T]) UnmarshalText(data []byte) error {
	return e.UnmarshalJSON(data)
}
