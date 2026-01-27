package handlers

import (
	"strings"
	"testing"
	"time"
)

func TestParseDurationParam(t *testing.T) {
	tests := []struct {
		name         string
		duration     string
		max          time.Duration
		wantDuration time.Duration
		wantErr      bool
		errSubstring string
	}{
		{"valid_duration", "2m", time.Minute * 5, time.Minute * 2, false, ""},
		{"empty_duration_is_required", "", time.Minute * 5, 0, true, "is required"},
		{"invalid_format", "alkindi", time.Minute * 5, 0, true, "invalid duration"},
		{"exceeds_max", "2m", time.Minute * 1, 0, true, "must be <="},
		{"zero_is_not_allowed", "0m", time.Minute * 1, 0, true, "greater than 0"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseDurationParam(tc.duration, tc.max)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseDurationParam() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr && tc.errSubstring != "" {
				if !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("error message %q does not contain %q", err.Error(), tc.errSubstring)
				}
				return
			}

			if got != tc.wantDuration {
				t.Errorf("got %v, want %v", got, tc.wantDuration)
			}
		})
	}
}

func TestParseOptionalDurationParam(t *testing.T) {
	tests := []struct {
		name         string
		duration     string
		max          time.Duration
		wantDuration time.Duration
		wantErr      bool
		errSubstring string
	}{
		{"valid_duration", "2m", time.Minute * 5, time.Minute * 2, false, ""},
		{"optional_duration", "", time.Minute * 5, 0, false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseOptionalDurationParam(tc.duration, tc.max)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseOptionalDurationParam() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr && tc.errSubstring != "" {
				if !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("error message %q does not contain %q", err.Error(), tc.errSubstring)
				}
				return
			}

			if got != tc.wantDuration {
				t.Errorf("got %v, want %v", got, tc.wantDuration)
			}
		})
	}
}
