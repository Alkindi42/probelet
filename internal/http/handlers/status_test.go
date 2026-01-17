package handlers

import "testing"

func TestParseStatusCode(t *testing.T) {
	tests := []struct {
		in   string
		want int
		ok   bool
	}{
		{"202", 202, true},
		{"99", 0, false},
		{"600", 0, false},
		{"abc", 0, false},
		{"", 0, false},
	}

	for _, tc := range tests {
		got, ok := parseStatusCode(tc.in)

		if ok != tc.ok {
			t.Fatalf("in=%q: expected ok=%v, got %v", tc.in, tc.ok, ok)
		}
		if got != tc.want {
			t.Fatalf("in=%q: expected %d, got %d", tc.in, tc.want, got)
		}
	}
}
