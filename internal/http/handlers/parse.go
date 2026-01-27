package handlers

import (
	"errors"
	"fmt"
	"time"
)

// parseDurationParam validates a duration string against format and a maximum limit.
// It returns the parsed time.Duration or an error with a user-friendly message.
func parseDurationParam(d string, max time.Duration) (time.Duration, error) {
	if d == "" {
		return 0, errors.New("duration is required (e.g. 100ms, 5s, 2m")
	}

	duration, err := time.ParseDuration(d)
	if err != nil {
		return 0, errors.New("invalid duration (examples: 100ms, 5s, 2m")
	}

	if duration <= 0 {
		return 0, errors.New("duration must be greater than 0")
	}

	if duration > max {
		return 0, fmt.Errorf("duration must be <= %s", max)
	}

	return duration, nil
}

// parseOptionalDurationParam attempts to parse a duration string.
// It returns 0 and no error if the string is empty, allowing for optional
// configuration. If a value is provided, it must satisfy all constraints
// defined in ParseDurationParam.
func parseOptionalDurationParam(d string, max time.Duration) (time.Duration, error) {
	if d == "" {
		return 0, nil
	}

	return parseDurationParam(d, max)
}
