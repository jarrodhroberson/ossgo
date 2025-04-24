package functions

import (
	"sync"
	"time"

	sls "github.com/jarrodhroberson/ossgo/slices"
	"github.com/joomcode/errorx"
)

func InsteadOfNil[T any](a *T, b *T) *T {
	if b == nil {
		panic(errorx.IllegalArgument.New("second argument to function \"b\" can not be \"nil\""))
	}
	if a == nil {
		return b
	}
	return a
}

// Deprecated: use strings.FirstNonEmpty instead
func FirstNonEmpty(data ...string) string {
	idx, err := sls.FindFirst[string](data, func(t string) bool {
		return t != ""
	})
	if err != nil {
		return ""
	}
	return data[idx]
}

// DeDuplicate creates a deduplicated version of the given function.
//
// - fn: The function to deduplicate.
//
// Returns a new function that behaves as a deduplicated version of `fn`.
func DeDuplicate[T comparable, R any](fn func(T) R) func(T) R {
	var (
		mu         sync.Mutex
		lastArg    T
		lastResult R
		lastCalled time.Time
		hasResult  bool // Track if we have a valid result yet.
	)

	return func(arg T) R {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		// Check if the function was called recently and with the same argument.
		if now.Sub(lastCalled) < time.Minute && arg == lastArg && hasResult { // Changed from duration to a fixed 1 minute to decouple it from the Debounce Duration
			return lastResult
		}

		// Execute the function and store the result.
		lastResult = fn(arg)
		lastArg = arg
		lastCalled = now
		hasResult = true
		return lastResult
	}
}

// Debounce creates a debounced version of the given function.
//
// - fn: The function to debounce.
// - duration: The debounce duration.
//
// Returns a new function that behaves as a debounced version of `fn`.
func Debounce[T any, R any](fn func(T) R, duration time.Duration) func(T) R {
	var (
		mu         sync.Mutex
		timer      *time.Timer
		resultChan = make(chan R, 1) // Channel to pass the result back
	)

	return func(arg T) R {
		mu.Lock()
		defer mu.Unlock()

		// If a timer is already running, stop it.
		if timer != nil {
			timer.Stop()
		}

		// Define a function to execute after the debounce duration.
		execute := func(arg T) {
			res := fn(arg)
			resultChan <- res // Send the result to the channel
		}

		// Start a new timer that will execute the function after the debounce duration.
		timer = time.AfterFunc(duration, func() {
			execute(arg)
		})

		select {
		case res := <-resultChan: // Try to receive a result from the channel immediately
			return res
		default: // No result available immediately; the debounced function hasn't executed yet.
			var zero R
			return zero // Return the zero value for the return type.
		}
	}
}
