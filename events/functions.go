package events

import (
	"github.com/jarrodhroberson/ossgo/cuid2"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

// New creates a new Event with a unique ID, the specified MIME type,
// the provided message, and the received timestamp.
//
// Parameters:
//   - mimeType: The MIME type of the message.
//   - message: A pointer to the message data.
//   - received: The timestamp when the message was received.
func New[T any](mimeType string, message *T, received timestamp.Timestamp) Event[T] {
	return event[T]{
		id: cuid2.New(16).String(),
		messageMimeType: mimeType,
		message:         message,
		receivedAt:      received,
	}
}
