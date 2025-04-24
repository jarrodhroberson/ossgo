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


// DebounceDeduplicate creates a debounced and deduplicated version of the given function.
//
//   - fn: The function to debounce and deduplicate. Must accept a single argument
//     and return a single value.
//   - duration: The debounce/deduplication duration.  The function will only be called
//     once per this duration.
//
// Returns a new function that behaves as a debounced and deduplicated version of `fn`.
func DebounceDeduplicate[T comparable, R any](fn func(T) R, duration time.Duration) func(T) R {
	var (
		mu         sync.Mutex
		timer      *time.Timer
		lastArg    T
		lastResult R
		lastCalled time.Time
		hasResult  bool // Track if we have a valid result yet.
	)

	return func(arg T) R {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		// Check if the function was called recently.  If so, and if the argument is the same,
		// return the last result.
		if now.Sub(lastCalled) < duration && arg == lastArg && hasResult {
			return lastResult
		}

		// If a timer is already running, stop it.  This effectively resets the debounce.
		if timer != nil {
			timer.Stop()
		}

		// Capture the current argument.
		lastArg = arg

		// Define a function to execute after the debounce duration.
		execute := func() {
			mu.Lock()
			defer mu.Unlock()

			// Execute the function and store the result.
			lastResult = fn(lastArg)
			lastCalled = time.Now()
			timer = nil // Clear the timer after execution.
			hasResult = true
		}

		// Start a new timer that will execute the function after the debounce duration.
		timer = time.AfterFunc(duration, execute)

		// Return the last cached result immediately.
		return lastResult
	}
}