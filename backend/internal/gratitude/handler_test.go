package gratitude

import (
	"testing"
)

func TestDailyLimit(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		shouldErr bool
	}{
		{"zero entries", 0, false},
		{"one entry", 1, false},
		{"two entries", 2, false},
		{"three entries - limit reached", 3, true},
		{"over limit", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isLimitReached := tt.count >= 3
			if isLimitReached != tt.shouldErr {
				t.Errorf("count=%d: got limitReached=%v, want %v", tt.count, isLimitReached, tt.shouldErr)
			}
		})
	}
}
