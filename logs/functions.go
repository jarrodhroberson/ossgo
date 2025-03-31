package logs

import (
	"resty.dev/v3"

	"github.com/rs/zerolog"
)

// ZerologToResty converts a zerolog.Logger to a resty.Logger.
// This allows you to use zerolog as the underlying logger for resty.
func ZerologToResty(log zerolog.Logger) resty.Logger {
	return restyLogger{
		log: log,
	}
}

// DecorateWithLogObjectMarshaller decorates a struct with zerolog.LogObjectMarshaler.
// This allows you to log the struct as an object in zerolog.
func DecorateWithLogObjectMarshaller[T any](s T) zerolog.LogObjectMarshaler {
	return lom[T]{
		delegate: s,
	}
}
