package events

import (
	"github.com/jarrodhroberson/ossgo/cuid2"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

func New[T any](mimeType string, message *T, received timestamp.Timestamp) Event[T] {
	return event[T]{
		id: cuid2.New(16).String(),
		messageMimeType: mimeType,
		message:         message,
		receivedAt:      received,
	}
}
