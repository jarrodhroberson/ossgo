package logs

import (
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jarrodhroberson/ossgo/containers"
	"github.com/jarrodhroberson/ossgo/functions/must"
)

// zeroLogAdapter
type zeroLogAdapter struct {
	logger *zerolog.Logger
}

func (z zeroLogAdapter) Debugf(format string, v ...interface{}) {
	z.logger.Debug().Msgf(format, v...)
}

func (z zeroLogAdapter) Errorf(format string, v ...interface{}) {
	z.logger.Error().Msgf(format, v)
}

func (z zeroLogAdapter) Infof(format string, v ...interface{}) {
	z.logger.Info().Msgf(format, v)
}

func (z zeroLogAdapter) Warnf(format string, v ...interface{}) {
	z.logger.Warn().Msgf(format, v)
}

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

type lom[T any] struct {
	delegate T
}

func (l lom[T]) MarshalZerologObject(e *zerolog.Event) {
	containers.WalkMap(must.MarshallMap(l.delegate), "", func(k string, v interface{}) {
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			e.Int64(k, v.(int64))
		case float32, float64:
			e.Float64(k, v.(float64))
		case string:
			e.Str(k, v.(string))
		case bool:
			e.Bool(k, v.(bool))
		default:
			log.Warn().Msgf("unknown type %s:", reflect.TypeOf(v))
			e.Str(k, fmt.Sprintf("%s", v))
		}
	})
}
