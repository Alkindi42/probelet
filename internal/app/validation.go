package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ValidationError represents an input validation failure.
type ValidationError struct {
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return e.Message
}

// ParseDurationParam validates a duration string against format and a maximum limit.
// It returns the parsed time.Duration or an error with a user-friendly message.
func ParseDurationParam(d string, max time.Duration) (time.Duration, error) {
	if d == "" {
		return 0, errors.New("duration is required (e.g. 100ms, 5s, 1m)")
	}

	duration, err := time.ParseDuration(d)
	if err != nil {
		return 0, errors.New("invalid duration (examples: 100ms, 5s, 1m)")
	}

	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	if duration > max {
		return 0, fmt.Errorf("duration must be <= %s", max)
	}

	return duration, nil
}

// ParseOptionalDurationParam attempts to parse a duration string.
// It returns 0 and no error if the string is empty, allowing for optional
// configuration. If a value is provided, it must satisfy all constraints
// defined in ParseDurationParam.
func ParseOptionalDurationParam(d string, max time.Duration) (time.Duration, error) {
	if d == "" {
		return 0, nil
	}

	return ParseDurationParam(d, max)
}

// ParseRateParam validates a floating-point rate parameter.
// It returns 0 when the input is empty. If provided, the value
// must be between 0 and 1 inclusive.
func ParseRateParam(r string) (float64, error) {
	if r == "" {
		return 0, nil
	}

	rate, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return 0, errors.New("invalid rate")
	}

	if rate < 0 || rate > 1 {
		return 0, errors.New("rate must be between 0 and 1")
	}

	return rate, nil
}
