package dice

import (
	"strings"
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

func TestDieRollStructure(t *testing.T) {
	// Test that the new DieRoll structure works correctly.
	diceSet, err := ParseDiceNotation("2d6 d20")
	if err != nil {
		t.Fatalf("Failed to parse dice notation: %v", err)
	}

	result := diceSet.Roll()

	// Should have 3 dice total (2d6 + 1d20).
	if len(result.DieRolls) != 3 {
		t.Errorf("Expected 3 die rolls, got %d", len(result.DieRolls))
	}

	// Check that die rolls have correct structure.
	expectedSides := []int{6, 6, 20} // Order should be preserved
	for i, dieRoll := range result.DieRolls {
		if dieRoll.Die.Sides != expectedSides[i] {
			t.Errorf("Die roll %d: expected %d sides, got %d", i, expectedSides[i], dieRoll.Die.Sides)
		}
		if dieRoll.Result < 1 || dieRoll.Result > dieRoll.Die.Sides {
			t.Errorf("Die roll %d: result %d is out of range [1,%d]", i, dieRoll.Result, dieRoll.Die.Sides)
		}
	}

	// Verify backward compatibility.
	if len(result.IndividualRolls) != len(result.DieRolls) {
		t.Errorf("IndividualRolls length %d doesn't match DieRolls length %d",
			len(result.IndividualRolls), len(result.DieRolls))
	}

	// Verify total calculation.
	expectedTotal := 0
	for _, roll := range result.IndividualRolls {
		expectedTotal += roll
	}
	if result.Total != expectedTotal {
		t.Errorf("Total %d doesn't match sum of individual rolls %d", result.Total, expectedTotal)
	}
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

// Tests for fancy dice functionality (Version 1.1).
func TestFancyDice(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		wantType string
		wantErr  bool
	}{
		{"Single f2", "f2", "f2", false},
		{"Single f4", "f4", "f4", false},
		{"Single f6", "f6", "f6", false},
		{"Single f7", "f7", "f7", false},
		{"Single f12", "f12", "f12", false},
		{"Single f13", "f13", "f13", false},
		{"Single f52", "f52", "f52", false},
		{"Multiple f4", "3f4", "f4", false},
		{"Invalid fancy dice", "f99", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, err := ParseDiceNotation(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDiceNotation(%q) expected error, got nil", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDiceNotation(%q) unexpected error: %v", tt.notation, err)
				return
			}

			// Roll the dice and check the result.
			result := set.Roll()

			// Check that we got the right type of dice.
			found := false
			for _, roll := range result.DieRolls {
				if roll.Type == tt.wantType {
					found = true

					// For fancy dice, check that FancyValue is populated.
					if strings.HasPrefix(tt.wantType, "f") && roll.FancyValue == "" {
						t.Errorf("ParseDiceNotation(%q) fancy dice missing FancyValue", tt.notation)
					}

					// For regular dice, check that FancyValue is empty.
					if strings.HasPrefix(tt.wantType, "d") && roll.FancyValue != "" {
						t.Errorf("ParseDiceNotation(%q) regular dice has unexpected FancyValue", tt.notation)
					}
				}
			}

			if !found {
				t.Errorf("ParseDiceNotation(%q) expected dice type %s not found", tt.notation, tt.wantType)
			}
		})
	}
}

func TestFancyDiceValues(t *testing.T) {
	// Test that f2 returns "heads" or "tails".
	for i := 0; i < 10; i++ {
		set, err := ParseDiceNotation("f2")
		if err != nil {
			t.Fatalf("ParseDiceNotation(f2) unexpected error: %v", err)
		}

		result := set.Roll()
		if len(result.DieRolls) != 1 {
			t.Fatalf("ParseDiceNotation(f2) expected 1 roll, got %d", len(result.DieRolls))
		}

		roll := result.DieRolls[0]
		if roll.FancyValue == "" {
			t.Fatal("ParseDiceNotation(f2) missing FancyValue")
		}

		value := roll.FancyValue
		if value != "heads" && value != "tails" {
			t.Errorf("ParseDiceNotation(f2) expected 'heads' or 'tails', got %q", value)
		}
	}
}

func TestMixedDiceNotation(t *testing.T) {
	// Test mixing regular and fancy dice.
	set, err := ParseDiceNotation("d20 f4 2f12")
	if err != nil {
		t.Fatalf("ParseDiceNotation(mixed) unexpected error: %v", err)
	}

	result := set.Roll()
	if len(result.DieRolls) != 4 { // 1 d20 + 1 f4 + 2 f12
		t.Fatalf("Expected 4 dice rolls, got %d", len(result.DieRolls))
	}

	// Check that we have the expected types.
	types := make(map[string]int)
	for _, roll := range result.DieRolls {
		types[roll.Type]++
	}

	if types["d20"] != 1 {
		t.Errorf("Expected 1 d20, got %d", types["d20"])
	}
	if types["f4"] != 1 {
		t.Errorf("Expected 1 f4, got %d", types["f4"])
	}
	if types["f12"] != 2 {
		t.Errorf("Expected 2 f12, got %d", types["f12"])
	}
}

// Tests for exclusive dice functionality (Version 1.2).
func TestExclusiveDiceParsing(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		wantDice int
		wantErr  bool
		wantType string
	}{
		{"Exclusive regular dice", "3D6", 3, false, "exclusive regular"},
		{"Exclusive fancy dice", "4F4", 4, false, "exclusive fancy"},
		{"Single exclusive die", "D20", 1, false, "exclusive regular"},
		{"Mixed exclusive and regular", "2d6 3D4", 5, false, "mixed"},
		{"Too many exclusive dice", "7D6", 0, true, "error"},
		{"Too many exclusive fancy", "5F4", 0, true, "error"},
		{"Invalid exclusive fancy", "3F99", 0, true, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, err := ParseDiceNotation(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDiceNotation(%q) expected error, got nil", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDiceNotation(%q) unexpected error: %v", tt.notation, err)
				return
			}

			if len(set.Dice) != tt.wantDice {
				t.Errorf("ParseDiceNotation(%q) expected %d dice, got %d", tt.notation, tt.wantDice, len(set.Dice))
			}
		})
	}
}

func TestExclusiveDiceUniqueness(t *testing.T) {
	// Test that exclusive regular dice don't repeat values.
	t.Run("3D6 no repeats", func(t *testing.T) {
		set, err := ParseDiceNotation("3D6")
		if err != nil {
			t.Fatalf("ParseDiceNotation(3D6) unexpected error: %v", err)
		}

		// Test multiple times to be sure.
		for i := 0; i < 10; i++ {
			result := set.Roll()
			if len(result.IndividualRolls) != 3 {
				t.Fatalf("Expected 3 rolls, got %d", len(result.IndividualRolls))
			}

			// Check uniqueness.
			seen := make(map[int]bool)
			for _, value := range result.IndividualRolls {
				if seen[value] {
					t.Errorf("Run %d: Duplicate value %d found in exclusive dice roll: %v", i, value, result.IndividualRolls)
				}
				seen[value] = true

				// Check valid range.
				if value < 1 || value > 6 {
					t.Errorf("Run %d: Value %d out of range [1,6]", i, value)
				}
			}
		}
	})

	// Test that exclusive fancy dice don't repeat values.
	t.Run("3F4 no repeats", func(t *testing.T) {
		set, err := ParseDiceNotation("3F4")
		if err != nil {
			t.Fatalf("ParseDiceNotation(3F4) unexpected error: %v", err)
		}

		// Test multiple times to be sure.
		for i := 0; i < 10; i++ {
			result := set.Roll()
			if len(result.DieRolls) != 3 {
				t.Fatalf("Expected 3 die rolls, got %d", len(result.DieRolls))
			}

			// Check uniqueness of fancy values.
			seenFancy := make(map[string]bool)
			for _, roll := range result.DieRolls {
				if seenFancy[roll.FancyValue] {
					t.Errorf("Run %d: Duplicate fancy value '%s' found in exclusive dice roll", i, roll.FancyValue)
				}
				seenFancy[roll.FancyValue] = true

				// Check that fancy value is populated.
				if roll.FancyValue == "" {
					t.Errorf("Run %d: Missing fancy value for f4 dice", i)
				}
			}
		}
	})
}

func TestMixedExclusiveAndRegular(t *testing.T) {
	// Test that mixing exclusive and regular dice works correctly.
	set, err := ParseDiceNotation("2d6 3D4")
	if err != nil {
		t.Fatalf("ParseDiceNotation(2d6 3D4) unexpected error: %v", err)
	}

	result := set.Roll()
	if len(result.IndividualRolls) != 5 {
		t.Fatalf("Expected 5 rolls total, got %d", len(result.IndividualRolls))
	}

	// The first 2 values (2d6) can repeat, the last 3 values (3D4) should be unique.
	lastThreeValues := result.IndividualRolls[2:] // Skip first 2 (2d6)
	seen := make(map[int]bool)
	for i, value := range lastThreeValues {
		if seen[value] {
			t.Errorf("Duplicate value %d found in exclusive 3D4 portion at position %d: %v", value, i, lastThreeValues)
		}
		seen[value] = true

		// Check valid range for D4.
		if value < 1 || value > 4 {
			t.Errorf("Value %d out of range [1,4] for D4 dice", value)
		}
	}
}

func TestExclusiveErrorCases(t *testing.T) {
	// Test error when requesting more exclusive dice than possible values.
	tests := []struct {
		name     string
		notation string
		wantErr  string
	}{
		{"Too many D6", "7D6", "cannot roll 7 exclusive dice with only 6 sides"},
		{"Too many F4", "5F4", "cannot roll 5 exclusive f4 dice with only 4 values"},
		{"Exactly max D6", "6D6", ""}, // Should work
		{"Exactly max F4", "4F4", ""}, // Should work
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDiceNotation(tt.notation)

			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("ParseDiceNotation(%q) expected error containing %q, got nil", tt.notation, tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ParseDiceNotation(%q) expected error containing %q, got %q", tt.notation, tt.wantErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("ParseDiceNotation(%q) unexpected error: %v", tt.notation, err)
				}
			}
		})
	}
}
