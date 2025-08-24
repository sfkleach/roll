// Package dice provides functionality for rolling dice and parsing dice notation.
package dice

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

// Die represents a single die with a specified number of sides.
type Die struct {
	Sides int
}

// DiceSet represents a collection of dice to be rolled together.
type DiceSet struct {
	Dice []Die
}

// DieRoll represents a single die roll with its result.
type DieRoll struct {
	Die        Die    // The die that was rolled
	Result     int    // The result of the roll
	Type       string // Type identifier (e.g., "d6", "f4")
	FancyValue string // For fancy dice, the display value (e.g., "♠", "heads")
}

// RollResult represents the result of rolling a set of dice.
type RollResult struct {
	DieRolls        []DieRoll // Individual die rolls with their dice info
	IndividualRolls []int     // Just the roll values (for backward compatibility)
	Total           int       // Sum of all rolls
}

// Standard values for fancy dice.
var fancyDiceValues = map[string][]string{
	"f2":  {"heads", "tails"},
	"f4":  {"♠", "♥", "♦", "♣"},           // Suit characters
	"f6":  {"⚀", "⚁", "⚂", "⚃", "⚄", "⚅"}, // Unicode dice faces (U+2680-U+2685)
	"f7":  {"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	"f12": {"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}, // Zodiac signs
	"f13": {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"},
	"f52": generatePlayingCards(),
}

// generatePlayingCards creates all 52 playing card symbols.
func generatePlayingCards() []string {
	suits := []string{"♠", "♥", "♦", "♣"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	cards := make([]string, 0, 52)
	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, rank+suit)
		}
	}
	return cards
}

// NewDie creates a new die with the specified number of sides.
func NewDie(sides int) Die {
	return Die{Sides: sides}
}

// Roll rolls a single die and returns the result.
func (d Die) Roll() int {
	if d.Sides <= 0 {
		// Handle fancy dice (negative sides) or invalid dice.
		if d.Sides < 0 {
			// This is a fancy die - return a random index + 1.
			fancyType := fmt.Sprintf("f%d", -d.Sides)
			if values, exists := fancyDiceValues[fancyType]; exists {
				return rand.IntN(len(values)) + 1
			}
		}
		return 0 // Defensive check: avoid rolling invalid dice.
	}
	return rand.IntN(d.Sides) + 1
}

// NewDiceSet creates a new dice set from the provided dice.
func NewDiceSet(dice []Die) DiceSet {
	return DiceSet{Dice: dice}
}

// Roll rolls all dice in the set and returns the results.
func (ds DiceSet) Roll() RollResult {
	dieRolls := make([]DieRoll, 0, len(ds.Dice)) // Pre-allocate with known capacity.
	rolls := make([]int, 0, len(ds.Dice))        // Pre-allocate with known capacity.
	total := 0

	for _, die := range ds.Dice {
		roll := die.Roll()

		var dieType string
		var fancyValue string

		if die.Sides < 0 {
			// This is a fancy die.
			fancyType := fmt.Sprintf("f%d", -die.Sides)
			dieType = fancyType

			if values, exists := fancyDiceValues[fancyType]; exists && roll > 0 && roll <= len(values) {
				fancyValue = values[roll-1] // Convert 1-based roll to 0-based index
			}
		} else {
			// Regular die.
			dieType = fmt.Sprintf("d%d", die.Sides)
			fancyValue = ""
		}

		dieRoll := DieRoll{
			Die:        die,
			Result:     roll,
			Type:       dieType,
			FancyValue: fancyValue,
		}
		dieRolls = append(dieRolls, dieRoll)
		rolls = append(rolls, roll)
		total += roll
	}

	return RollResult{
		DieRolls:        dieRolls,
		IndividualRolls: rolls, // For backward compatibility
		Total:           total,
	}
}

// ParseDiceNotation parses dice notation and returns a DiceSet.
// Supports multiple formats:
// - "3d6" - three six-sided dice
// - "d20" - one twenty-sided die (count defaults to 1)
// - "2d10 d6" - space-separated groups
// - "1d20,7d4" - comma-separated groups
// - "3d6+2d4" - plus-separated groups
// Returns an error if the notation is invalid.
func ParseDiceNotation(notation string) (DiceSet, error) {
	notation = strings.TrimSpace(notation)
	if notation == "" {
		return DiceSet{}, fmt.Errorf("empty dice notation")
	}

	// Split by separators (space, comma, plus).
	parts := splitDiceExpression(notation)

	var allDice []Die

	for _, part := range parts {
		dice, err := parseSingleDiceGroup(part)
		if err != nil {
			return DiceSet{}, err
		}
		allDice = append(allDice, dice...)
	}

	if len(allDice) == 0 {
		return DiceSet{}, fmt.Errorf("no valid dice found in notation: %s", notation)
	}

	return NewDiceSet(allDice), nil
}

// splitDiceExpression splits a dice expression by separators (space, comma, plus).
func splitDiceExpression(notation string) []string {
	// Replace all separators with spaces for consistent splitting.
	notation = strings.ReplaceAll(notation, ",", " ")
	notation = strings.ReplaceAll(notation, "+", " ")

	// Split by whitespace and filter out empty parts.
	parts := strings.Fields(notation)
	return parts
}

// parseSingleDiceGroup parses a single dice group like "3d6", "d20", or "2f4".
func parseSingleDiceGroup(group string) ([]Die, error) {
	group = strings.TrimSpace(group)
	if group == "" {
		return nil, fmt.Errorf("empty dice group")
	}

	// Check for fancy dice notation first: [count]f[type]
	fancyRe := regexp.MustCompile(`^(\d*)f(\d+)$`)
	if matches := fancyRe.FindStringSubmatch(group); matches != nil {
		return parseFancyDice(matches[1], matches[2])
	}

	// Regular dice notation: [count]d[sides]
	regularRe := regexp.MustCompile(`^(\d*)d(\d+)$`)
	matches := regularRe.FindStringSubmatch(group)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid dice notation: %s", group)
	}

	// Parse count (default to 1 if empty).
	countStr := matches[1]
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid number of dice: %s", countStr)
		}
	}

	// Parse sides.
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid number of sides: %s", matches[2])
	}

	// Validate values.
	if count <= 0 {
		return nil, fmt.Errorf("dice count must be positive, got: %d", count)
	}
	if sides <= 0 {
		return nil, fmt.Errorf("dice sides must be positive, got: %d", sides)
	}

	// Create dice.
	var dice []Die
	for i := 0; i < count; i++ {
		dice = append(dice, NewDie(sides))
	}

	return dice, nil
}

// parseFancyDice parses fancy dice notation and creates special "dice" with negative sides to mark them as fancy.
func parseFancyDice(countStr, typeStr string) ([]Die, error) {
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			return nil, fmt.Errorf("invalid dice count: %s", countStr)
		}
	}

	fancyType := "f" + typeStr
	if _, exists := fancyDiceValues[fancyType]; !exists {
		return nil, fmt.Errorf("unsupported fancy dice type: %s", fancyType)
	}

	// Create "dice" with negative sides to mark them as fancy dice.
	// We'll encode the fancy type in the sides value.
	fancyTypeNum, _ := strconv.Atoi(typeStr)
	var dice []Die
	for i := 0; i < count; i++ {
		// Use negative sides to indicate fancy dice.
		dice = append(dice, Die{Sides: -fancyTypeNum})
	}

	return dice, nil
}

// String returns a string representation of the dice set.
func (ds DiceSet) String() string {
	if len(ds.Dice) == 0 {
		return "empty dice set"
	}

	// Count dice by sides for compact representation.
	sidesCounts := make(map[int]int)
	for _, die := range ds.Dice {
		sidesCounts[die.Sides]++
	}

	parts := make([]string, 0, len(sidesCounts)) // Pre-allocate with estimated capacity.
	for sides, count := range sidesCounts {
		parts = append(parts, fmt.Sprintf("%dd%d", count, sides))
	}

	return fmt.Sprintf("DiceSet{%v}", parts)
}
