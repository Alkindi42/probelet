package app

import (
	"context"
	"strings"
	"testing"
)

func TestRunFlaky_Validation(t *testing.T) {
	tests := []struct {
		name         string
		req          FlakyRequest
		wantErr      bool
		errSubstring string
	}{
		{
			name:         "missing_both_rates",
			req:          FlakyRequest{},
			wantErr:      true,
			errSubstring: "at least one of 'error_rate' or 'delay_rate' must be greater than 0",
		},
		{
			name: "invalid_error_rate",
			req: FlakyRequest{
				ErrorRate: "abc",
			},
			wantErr:      true,
			errSubstring: "error_rate invalid rate",
		},
		{
			name: "invalid_delay_rate",
			req: FlakyRequest{
				DelayRate: "abc",
			},
			wantErr:      true,
			errSubstring: "delay_rate invalid rate",
		},
		{
			name: "rate_out_of_range",
			req: FlakyRequest{
				ErrorRate: "1.5",
			},
			wantErr:      true,
			errSubstring: "error_rate rate must be between 0 and 1",
		},
		{
			name: "sum_of_rates_above_one",
			req: FlakyRequest{
				ErrorRate: "0.8",
				DelayRate: "0.3",
				MaxDelay:  "1s",
			},
			wantErr:      true,
			errSubstring: "the sum of 'error_rate' and 'delay_rate' must be <=",
		},
		{
			name: "delay_rate_requires_max_delay",
			req: FlakyRequest{
				DelayRate: "0.5",
			},
			wantErr:      true,
			errSubstring: "duration is required",
		},
		{
			name: "invalid_max_delay",
			req: FlakyRequest{
				DelayRate: "0.5",
				MaxDelay:  "abc",
			},
			wantErr:      true,
			errSubstring: "invalid duration",
		},
		{
			name: "max_delay_too_large",
			req: FlakyRequest{
				DelayRate: "0.5",
				MaxDelay:  "2m",
			},
			wantErr:      true,
			errSubstring: "duration must be <=",
		},
		{
			name: "invalid_status_code_non_integer",
			req: FlakyRequest{
				ErrorRate: "1",
				Status:    "abc",
			},
			wantErr:      true,
			errSubstring: "invalid status code",
		},
		{
			name: "invalid_status_code_out_of_range",
			req: FlakyRequest{
				ErrorRate: "1",
				Status:    "99",
			},
			wantErr:      true,
			errSubstring: "status code must be between 100 and 599",
		},
		{
			name: "valid_error_only",
			req: FlakyRequest{
				ErrorRate: "1",
			},
			wantErr: false,
		},
		{
			name: "valid_delay_only",
			req: FlakyRequest{
				DelayRate: "1",
				MaxDelay:  "5ms",
			},
			wantErr: false,
		},
		{
			name: "valid_error_and_delay",
			req: FlakyRequest{
				ErrorRate: "0.2",
				DelayRate: "0.3",
				MaxDelay:  "5ms",
				Status:    "504",
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RunFlaky(context.Background(), tc.req)

			if (err != nil) != tc.wantErr {
				t.Fatalf("RunFlaky() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr && tc.errSubstring != "" {
				if !strings.Contains(err.Error(), tc.errSubstring) {
					t.Fatalf("error message %q does not contain %q", err.Error(), tc.errSubstring)
				}
			}
		})
	}
}
