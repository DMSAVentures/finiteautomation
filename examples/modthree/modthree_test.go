package modthree

import "testing"

// Test ModThreeGeneric
func TestModThreeGeneric(t *testing.T) {
	mod3 := NewModThreeGeneric()

	tests := []struct {
		binary   string
		expected ModState
	}{
		{"", 0},
		{"0", 0},
		{"1", 1},
		{"10", 2},
		{"11", 0},
		{"100", 1},
		{"101", 2},
		{"110", 0},
		{"111", 1},
		{"1101", 1},
		{"1110", 2},
		{"1111", 0},
	}

	for _, tt := range tests {
		t.Run(tt.binary, func(t *testing.T) {
			_ = mod3.IsDivisibleByThree(tt.binary)
			if mod3.fsm.CurrentState() != tt.expected {
				t.Errorf(" %s, want %d", tt.binary, tt.expected)
			}
		})
	}
}
