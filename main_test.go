package main

import (
	"bytes"
	"io"
	"os"
	"strings"
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

func TestProcessDiceExpression(t *testing.T) {
	// Test the processDiceExpression function used in interactive mode.
	// Capture stdout to verify the output format.

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test a simple dice expression.
	processDiceExpression("1d6", false, false)

	// Restore stdout and read the output.
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output contains expected patterns.
	if !strings.Contains(output, "d6:") {
		t.Errorf("Expected output to contain 'd6:', got: %s", output)
	}
	if !strings.Contains(output, "Total:") {
		t.Errorf("Expected output to contain 'Total:', got: %s", output)
	}
}

func TestProcessDiceExpressionError(t *testing.T) {
	// Test error handling in processDiceExpression.

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test an invalid dice expression.
	processDiceExpression("invalid", false, false)

	// Restore stdout and read the output.
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output contains an error message.
	if !strings.Contains(output, "Error parsing dice notation") {
		t.Errorf("Expected output to contain error message, got: %s", output)
	}
}
