package events

import (
	"github.com/jarrodhroberson/ossgo/cuid2"
	"github.com/jarrodhroberson/ossgo/timestamp"
)

func New(mimeType string, message string, received timestamp.Timestamp) Event {
	return event{
		Id:              cuid2.New(16),
		MessageMimeType: mimeType,
		Message:         message,
		ReceivedAt:      received,
	}
}
