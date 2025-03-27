package logs

import (
	"resty.dev/v3"

	"github.com/rs/zerolog"
)

func ZerologToResty(log zerolog.Logger) resty.Logger {
	return restyLogger{
		log: log,
	}
}

func DecorateWithLogObjectMarshaller[T any](s T) zerolog.LogObjectMarshaler {
	return lom[T]{
		delegate: s,
	}
}
