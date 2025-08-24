package main

import (
	"testing"

	"github.com/sfkleach/roll/internal/dice"
)

func TestDiceIntegration(t *testing.T) {
	// Test basic integration without GUI components.
	// This tests that our core dice functionality works correctly.

	diceSet, err := dice.ParseDiceNotation("3d6")
	if err != nil {
		t.Fatalf("Failed to parse dice notation: %v", err)
	}

	result := diceSet.Roll()

	if len(result.IndividualRolls) != 3 {
		t.Errorf("Expected 3 dice rolls, got %d", len(result.IndividualRolls))
	}

	if result.Total <= 0 {
		t.Errorf("Expected positive total, got %d", result.Total)
	}

	// Verify that the total matches the sum of individual rolls.
	expectedTotal := 0
	for _, roll := range result.IndividualRolls {
		expectedTotal += roll
		if roll < 1 || roll > 6 {
			t.Errorf("Individual roll %d is out of valid range [1,6]", roll)
		}
	}

	if result.Total != expectedTotal {
		t.Errorf("Total %d doesn't match sum of individual rolls %d", result.Total, expectedTotal)
	}
}
