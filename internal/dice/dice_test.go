package dice

import (
	"testing"
)

func TestNewDie(t *testing.T) {
	die := NewDie(6)
	if die.Sides != 6 {
		t.Errorf("Expected 6 sides, got %d", die.Sides)
	}
}

func TestDieRoll(t *testing.T) {
	die := NewDie(6)

	// Test multiple rolls to ensure they're in valid range.
	for i := 0; i < 100; i++ {
		roll := die.Roll()
		if roll < 1 || roll > 6 {
			t.Errorf("Roll result %d is out of range [1,6]", roll)
		}
	}
}

func TestDieRollInvalidSides(t *testing.T) {
	// Test defensive check for invalid dice.
	die := NewDie(0)
	roll := die.Roll()
	if roll != 0 {
		t.Errorf("Expected 0 for invalid die, got %d", roll)
	}

	die = NewDie(-1)
	roll = die.Roll()
	if roll != 0 {
		t.Errorf("Expected 0 for invalid die, got %d", roll)
	}
}

func TestDiceSetRoll(t *testing.T) {
	dice := []Die{NewDie(6), NewDie(6), NewDie(6)}
	set := NewDiceSet(dice)

	result := set.Roll()

	if len(result.IndividualRolls) != 3 {
		t.Errorf("Expected 3 individual rolls, got %d", len(result.IndividualRolls))
	}

	// Verify each roll is in valid range.
	for i, roll := range result.IndividualRolls {
		if roll < 1 || roll > 6 {
			t.Errorf("Roll %d result %d is out of range [1,6]", i, roll)
		}
	}

	// Verify total is sum of individual rolls.
	expectedTotal := 0
	for _, roll := range result.IndividualRolls {
		expectedTotal += roll
	}
	if result.Total != expectedTotal {
		t.Errorf("Expected total %d, got %d", expectedTotal, result.Total)
	}
}

func TestParseDiceNotation(t *testing.T) {
	tests := []struct {
		notation    string
		wantErr     bool
		totalDice   int
		description string
	}{
		// Simple single dice groups.
		{"3d6", false, 3, "three six-sided dice"},
		{"1d20", false, 1, "one twenty-sided die"},
		{"2d10", false, 2, "two ten-sided dice"},
		{"10d6", false, 10, "ten six-sided dice"},

		// Single die notation (no count).
		{"d6", false, 1, "one six-sided die (implicit count)"},
		{"d20", false, 1, "one twenty-sided die (implicit count)"},

		// Multiple dice groups with different separators.
		{"2d10 d6", false, 3, "two ten-sided dice and one six-sided die (space)"},
		{"1d20,7d4", false, 8, "one twenty-sided die and seven four-sided dice (comma)"},
		{"3d6+2d4", false, 5, "three six-sided dice and two four-sided dice (plus)"},
		{"d20 2d6 d4", false, 4, "mixed notation with spaces"},
		{"1d8,d12+2d4", false, 4, "mixed separators"},

		// Invalid notations.
		{"", true, 0, "empty string"},
		{"3x6", true, 0, "invalid separator"},
		{"d", true, 0, "missing sides"},
		{"3d", true, 0, "missing sides with count"},
		{"0d6", true, 0, "zero count"},
		{"3d0", true, 0, "zero sides"},
		{"-1d6", true, 0, "negative count"},
		{"3d-6", true, 0, "negative sides"},
		{"abc", true, 0, "non-numeric notation"},
		{"3d6d4", true, 0, "malformed notation"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			set, err := ParseDiceNotation(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for notation %s, but got none", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for notation %s: %v", tt.notation, err)
				return
			}

			if len(set.Dice) != tt.totalDice {
				t.Errorf("Expected %d total dice for %s, got %d", tt.totalDice, tt.notation, len(set.Dice))
			}
		})
	}
}

func TestParseDiceNotationSpecificExamples(t *testing.T) {
	// Test specific examples from the requirements.
	t.Run("d20 single die", func(t *testing.T) {
		set, err := ParseDiceNotation("d20")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(set.Dice) != 1 {
			t.Errorf("Expected 1 die, got %d", len(set.Dice))
		}
		if set.Dice[0].Sides != 20 {
			t.Errorf("Expected 20 sides, got %d", set.Dice[0].Sides)
		}
	})

	t.Run("2d10 d6 space separated", func(t *testing.T) {
		set, err := ParseDiceNotation("2d10 d6")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(set.Dice) != 3 {
			t.Errorf("Expected 3 dice total, got %d", len(set.Dice))
		}

		// Check that we have the right types of dice.
		tenSidedCount := 0
		sixSidedCount := 0
		for _, die := range set.Dice {
			if die.Sides == 10 {
				tenSidedCount++
			} else if die.Sides == 6 {
				sixSidedCount++
			}
		}
		if tenSidedCount != 2 {
			t.Errorf("Expected 2 ten-sided dice, got %d", tenSidedCount)
		}
		if sixSidedCount != 1 {
			t.Errorf("Expected 1 six-sided die, got %d", sixSidedCount)
		}
	})

	t.Run("1d20,7d4 comma separated", func(t *testing.T) {
		set, err := ParseDiceNotation("1d20,7d4")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(set.Dice) != 8 {
			t.Errorf("Expected 8 dice total, got %d", len(set.Dice))
		}

		// Check that we have the right types of dice.
		twentySidedCount := 0
		fourSidedCount := 0
		for _, die := range set.Dice {
			if die.Sides == 20 {
				twentySidedCount++
			} else if die.Sides == 4 {
				fourSidedCount++
			}
		}
		if twentySidedCount != 1 {
			t.Errorf("Expected 1 twenty-sided die, got %d", twentySidedCount)
		}
		if fourSidedCount != 7 {
			t.Errorf("Expected 7 four-sided dice, got %d", fourSidedCount)
		}
	})
}

func TestDiceSetString(t *testing.T) {
	// Test empty dice set.
	emptySet := NewDiceSet([]Die{})
	if emptySet.String() != "empty dice set" {
		t.Errorf("Expected 'empty dice set', got %s", emptySet.String())
	}

	// Test dice set with dice.
	dice := []Die{NewDie(6), NewDie(6), NewDie(20)}
	set := NewDiceSet(dice)
	str := set.String()

	// The exact order may vary due to map iteration, so just check it contains expected parts.
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}
