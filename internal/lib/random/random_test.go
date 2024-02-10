package random

import (
	"testing"
	"unicode"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"ZeroLength", 0},
		{"SmallLength", 10},
		{"LargeLength", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.length)
			if len(got) != tt.length {
				t.Errorf("NewRandomString() = %v, want %v", len(got), tt.length)
			}
			for _, runeValue := range got {
				if !unicode.IsLetter(runeValue) && !unicode.IsNumber(runeValue) {
					t.Errorf("NewRandomString() contains invalid characters = %v", runeValue)
				}
			}
		})
	}
}