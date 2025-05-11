package calculator

import (
	"testing"
)

func TestCalculator(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected float64
		wantErr  bool
	}{
		{"simple addition", "2+2", 4, false},
		{"multiplication before addition", "2+2*2", 6, false},
		{"with parentheses", "(2+2)*2", 8, false},
		{"division", "10/2", 5, false},
		{"decimal", "2.5 + 3.5", 6, false},
		{"invalid expression", "2 + a", 0, true},
		{"division by zero", "2/0", 0, true},
		{"mismatched parentheses", "(2+2", 0, true},
	}

	c := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.Calculate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expr, result, tt.expected)
			}
		})
	}
}
