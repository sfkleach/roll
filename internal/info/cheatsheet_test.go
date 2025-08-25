package info

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	// Test that GetVersion returns the current version
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() should not return empty string")
	}

	// Should default to "dev" in test environment
	if version != "dev" {
		t.Logf("Version is %s (expected 'dev' in test environment, but this may be overridden)", version)
	}
}

func TestGetCheatsheetContent(t *testing.T) {
	content := GetCheatsheetContent()

	// Test that content contains expected sections
	expectedSections := []string{
		"Roll Dice Application",
		"BASIC DICE NOTATION",
		"FANCY DICE",
		"EXCLUSIVE DICE",
		"SORTING OPTIONS",
		"EXAMPLES",
	}

	for _, section := range expectedSections {
		if !strings.Contains(content, section) {
			t.Errorf("Cheatsheet content missing expected section: %s", section)
		}
	}

	// Test that version is included in content
	if !strings.Contains(content, GetVersion()) {
		t.Error("Cheatsheet content should include current version")
	}
}

func TestGetCheatsheetMarkdown(t *testing.T) {
	content := GetCheatsheetMarkdown()

	// Test that content contains markdown formatting
	if !strings.Contains(content, "# Roll Dice Application") {
		t.Error("Markdown content should contain H1 header")
	}

	if !strings.Contains(content, "## BASIC DICE NOTATION") {
		t.Error("Markdown content should contain H2 headers")
	}

	if !strings.Contains(content, "**d20**") {
		t.Error("Markdown content should contain bold formatting")
	}

	// Test that version is included in markdown content
	if !strings.Contains(content, GetVersion()) {
		t.Error("Markdown cheatsheet content should include current version")
	}
}
