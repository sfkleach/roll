// Package info provides shared information content for the application.
package info

import (
	"fmt"
	"regexp"
	"strings"
)

// Version information - will be set at build time via ldflags.
// Default value is used for development builds.
var Version = "dev"

// getCheatsheetMarkdownSource returns the single source of truth for cheatsheet content.
func getCheatsheetMarkdownSource() string {
	return fmt.Sprintf(`# Roll Dice Application v%s

## Cheatsheet

### BASIC DICE NOTATION:
- **d20** - Roll a single 20-sided die  
- **3d6** - Roll three 6-sided dice  
- **2d10 d6** - Roll two 10-sided dice and one 6-sided die  
- **1d20,7d4** - Roll one 20-sided die and seven 4-sided dice  

### FANCY DICE (Custom Unicode Characters):
- **f2** - Two-sided coin (heads/tails)  
- **f4** - Four-sided die with suit symbols (♠♥♦♣)  
- **f6** - Six-sided die with dot patterns (⚀⚁⚂⚃⚄⚅)  
- **f7** - Seven-sided die with days of week (Mon-Sun)  
- **f12** - Twelve-sided die with zodiac signs  
- **f13** - Thirteen-sided die with card ranks (A,2-10,J,Q,K)  
- **f52** - Fifty-two-sided die with playing cards  

### CUSTOM FANCY DICE:
- **--fancy=GLOB** - Load custom fancy dice from files matching pattern  
- File format: one line per value as "name, value" or just "name"  
- Example: **--fancy='*.dice'** loads all .dice files  

### EXCLUSIVE DICE (No Repeats in Group):
- **3D6** - Roll three 6-sided dice with no duplicate values  
- **5D20** - Roll five 20-sided dice with no duplicate values  
- **13F52** - Roll thirteen cards with no duplicates  

### SORTING OPTIONS:
- **-a** or **--ascending** - Sort results in ascending order  
- **-d** or **--descending** - Sort results in descending order  

### EXAMPLES:
- roll 3d6 2d10  
- roll --ascending 5D20  
- roll f52 f52 f52  
- roll --fancy='colors.dice' fcolors  
- -a 3d6 (in GUI)  
- --descending 2d20 3d4 (in GUI)  
`, Version)
}

// markdownToPlainText converts markdown to plain text using simple string replacement.
func markdownToPlainText(md string) string {
	text := md

	// Remove headers - convert ### to nothing, ## to nothing, # to nothing
	headerRegex := regexp.MustCompile(`^#{1,3}\s*`)
	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		// Check if this was a header line before processing
		wasHeader := headerRegex.MatchString(line)

		// Remove markdown headers
		line = headerRegex.ReplaceAllString(line, "")

		// Remove bold formatting
		line = strings.ReplaceAll(line, "**", "")

		// Convert markdown list items to bullet points
		if strings.HasPrefix(line, "- ") {
			line = "• " + line[2:]
		}

		// Clean up extra whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines and the redundant "Cheatsheet" line
		if line != "" && line != "Cheatsheet" {
			// Add extra spacing before headers (except the first one)
			if wasHeader && len(processedLines) > 0 {
				processedLines = append(processedLines, "")
			}
			processedLines = append(processedLines, line)
		}
	}

	// Join lines
	result := strings.Join(processedLines, "\n")

	return strings.TrimSpace(result)
}

// GetCheatsheetContent returns the unified cheatsheet content in plain text format.
func GetCheatsheetContent() string {
	markdown := getCheatsheetMarkdownSource()
	return markdownToPlainText(markdown)
}

// GetCheatsheetMarkdown returns the cheatsheet content formatted for GUI display.
func GetCheatsheetMarkdown() string {
	return getCheatsheetMarkdownSource()
}

// GetVersion returns the current version string.
func GetVersion() string {
	return Version
}
