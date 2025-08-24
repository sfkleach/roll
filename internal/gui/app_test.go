package gui

import (
	"testing"
)

func TestParseFlagsFromInput(t *testing.T) {
	tests := []struct {
		input           string
		expectedNotation string
		expectedAsc     bool
		expectedDesc    bool
		expectedError   bool
	}{
		{"3d6", "3d6", false, false, false},
		{"-a 3d6", "3d6", true, false, false},
		{"--ascending 3d6", "3d6", true, false, false},
		{"-d 3d6", "3d6", false, true, false},
		{"--descending 3d6", "3d6", false, true, false},
		{"3d6 -a", "3d6", true, false, false},
		{"3d6 --descending", "3d6", false, true, false},
		{"-a 2d10 d6", "2d10 d6", true, false, false},
		{"--descending 2d20 3d4", "2d20 3d4", false, true, false},
		{"-a -d 3d6", "", false, false, true}, // Error: both flags
		{"--ascending --descending 3d6", "", false, false, true}, // Error: both flags
		{"-a --descending 3d6", "", false, false, true}, // Error: both flags
		{"-d -a 3d6", "", false, false, true}, // Error: both flags
	}

	for _, test := range tests {
		notation, asc, desc, err := parseFlagsFromInput(test.input)
		
		if test.expectedError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.input)
			}
			continue
		}
		
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		
		if notation != test.expectedNotation {
			t.Errorf("Input '%s': expected notation '%s', got '%s'", test.input, test.expectedNotation, notation)
		}
		
		if asc != test.expectedAsc {
			t.Errorf("Input '%s': expected ascending %v, got %v", test.input, test.expectedAsc, asc)
		}
		
		if desc != test.expectedDesc {
			t.Errorf("Input '%s': expected descending %v, got %v", test.input, test.expectedDesc, desc)
		}
	}
}
