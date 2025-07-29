package logs

import (
	"resty.dev/v3"

	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go/v82"
)

// StripeLeveledLogger creates a Stripe leveled logger adapter that wraps a zerolog.Logger.
// It implements the stripe.LeveledLoggerInterface interface to enable Stripe SDK logging
// through zerolog.
//
// Parameters:
//   - l: A pointer to a zerolog.Logger instance that will be used for logging
//
// Returns:
//   - stripe.LeveledLoggerInterface: An adapter that implements Stripe's logging interface
func StripeLeveledLogger(l *zerolog.Logger) stripe.LeveledLoggerInterface {
	return zeroLogAdapter{
		logger: l,
	}
}

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
