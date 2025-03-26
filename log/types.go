package log

import (
	"github.com/rs/zerolog"
)

type restyLogger struct {
	log zerolog.Logger
}

func (l restyLogger) Errorf(format string, v ...interface{}) {
	l.log.Error().Msgf(format, v...)
}

func (l restyLogger) Warnf(format string, v ...interface{}) {
	l.log.Warn().Msgf(format, v...)
}

func (l restyLogger) Debugf(format string, v ...interface{}) {
	l.log.Debug().Msgf(format, v...)
}

func (l restyLogger) Tracef(format string, v ...interface{}) {
	l.log.Trace().Msgf(format, v...)
}

func (l restyLogger) Info(args ...interface{}) {
	l.log.Info().Msg(args[0].(string))
}
