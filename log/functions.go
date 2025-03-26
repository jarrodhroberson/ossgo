package log

import (
	"resty.dev/v3"

	"github.com/rs/zerolog"
)



func ZerologToResty(log zerolog.Logger) resty.Logger {
	return restyLogger{
		log: log,
	}
}
