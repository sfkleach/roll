// Package dice provides functionality for rolling dice and parsing dice notation.
package dice

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
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

// FancyDieValue represents a single value for a fancy die.
type FancyDieValue struct {
	Name  string // Display name (e.g., "heads", "♠", "Mon")
	Value int    // Scoring value
}

// RollResult represents the result of rolling a set of dice.
type RollResult struct {
	DieRolls        []DieRoll // Individual die rolls with their dice info
	IndividualRolls []int     // Just the roll values (for backward compatibility)
	Total           int       // Sum of all rolls
}

// Standard values for fancy dice.
var fancyDiceValues = map[string][]FancyDieValue{
	"f2":  {{"heads", 1}, {"tails", 0}},
	"f4":  {{"♠", 4}, {"♥", 3}, {"♦", 2}, {"♣", 1}},                           // Suit characters
	"f6":  {{"1⚀", 1}, {"2⚁", 2}, {"3⚂", 3}, {"4⚃", 4}, {"5⚄", 5}, {"6⚅", 6}}, // Unicode dice faces (U+2680-U+2685)
	"f7":  {{"Mon", 1}, {"Tue", 2}, {"Wed", 3}, {"Thu", 4}, {"Fri", 5}, {"Sat", 6}, {"Sun", 7}},
	"f12": generateZodiacValues(),
	"f13": {{"A", 4}, {"2", 0}, {"3", 0}, {"4", 0}, {"5", 0}, {"6", 0}, {"7", 0}, {"8", 0}, {"9", 0}, {"10", 0}, {"J", 1}, {"Q", 2}, {"K", 3}},
	"f52": generatePlayingCardValues(),
}

// generateZodiacValues creates zodiac sign values.
func generateZodiacValues() []FancyDieValue {
	zodiacSigns := []string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}
	values := make([]FancyDieValue, len(zodiacSigns))
	for i, sign := range zodiacSigns {
		values[i] = FancyDieValue{Name: sign, Value: i + 1}
	}
	return values
}

// LoadCustomFancyDice loads custom fancy dice from files matching the glob pattern.
func LoadCustomFancyDice(globPattern string) error {
	files, err := filepath.Glob(globPattern)
	if err != nil {
		return fmt.Errorf("invalid glob pattern '%s': %v", globPattern, err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found matching pattern '%s'", globPattern)
	}

	for _, file := range files {
		err := loadSingleFancyDiceFile(file)
		if err != nil {
			return fmt.Errorf("error loading file '%s': %v", file, err)
		}
	}

	return nil
}

// loadSingleFancyDiceFile loads a single fancy dice file.
func loadSingleFancyDiceFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	var values []FancyDieValue
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line.
		value, err := parseFancyDiceLine(line, len(values)+1)
		if err != nil {
			return fmt.Errorf("line %d: %v", lineNum, err)
		}

		values = append(values, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(values) == 0 {
		return fmt.Errorf("file contains no valid fancy dice values")
	}

	// The dice type is determined by the number of values (rank of the dice).
	diceType := fmt.Sprintf("f%d", len(values))

	// Store the custom fancy dice values.
	fancyDiceValues[diceType] = values

	return nil
}

// parseFancyDiceLine parses a single line from a fancy dice file.
// Format: "name, value" or "name" (defaults to position).
func parseFancyDiceLine(line string, defaultValue int) (FancyDieValue, error) {
	parts := strings.Split(line, ",")

	if len(parts) == 1 {
		// Just name, use default value.
		name := strings.TrimSpace(parts[0])
		if name == "" {
			return FancyDieValue{}, fmt.Errorf("empty name")
		}
		return FancyDieValue{Name: name, Value: defaultValue}, nil
	} else if len(parts) == 2 {
		// Name and value.
		name := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		if name == "" {
			return FancyDieValue{}, fmt.Errorf("empty name")
		}

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return FancyDieValue{}, fmt.Errorf("invalid value '%s': must be an integer", valueStr)
		}

		return FancyDieValue{Name: name, Value: value}, nil
	} else {
		return FancyDieValue{}, fmt.Errorf("invalid format: expected 'name' or 'name, value'")
	}
}

// generatePlayingCardValues creates all 52 playing card values.
func generatePlayingCardValues() []FancyDieValue {
	suits := []string{"♣", "♦", "♥", "♠"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	cards := make([]FancyDieValue, 0, 52)
	for _, suit := range suits {
		for _, rank := range ranks {
			// Add numerical position (1-52) alongside the card symbol.
			card := fmt.Sprintf("%s%s", rank, suit)
			cards = append(cards, FancyDieValue{Name: card, Value: len(cards) + 1})
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

	// Group dice by exclusivity for proper handling.
	exclusiveGroups := ds.groupExclusiveDice()

	for _, group := range exclusiveGroups {
		if group.IsExclusive {
			// Roll exclusive group without replacement.
			values := ds.rollExclusiveGroup(group)
			for i, value := range values {
				die := group.Dice[i]

				var dieType string
				var fancyValue string

				if group.IsFancy {
					// Exclusive fancy dice.
					originalType := -(die.Sides + 1000)
					fancyType := fmt.Sprintf("f%d", originalType)
					dieType = fancyType

					if fancyValues, exists := fancyDiceValues[fancyType]; exists && value > 0 && value <= len(fancyValues) {
						fancyValue = fancyValues[value-1].Name
						total += fancyValues[value-1].Value // Add the scoring value to total
					}

					// Create display die with original sides.
					displayDie := Die{Sides: -originalType}
					dieRoll := DieRoll{
						Die:        displayDie,
						Result:     value,
						Type:       dieType,
						FancyValue: fancyValue,
					}
					dieRolls = append(dieRolls, dieRoll)
				} else {
					// Exclusive regular dice.
					originalSides := die.Sides - 1000
					dieType = fmt.Sprintf("d%d", originalSides)

					// Create display die with original sides.
					displayDie := Die{Sides: originalSides}
					dieRoll := DieRoll{
						Die:        displayDie,
						Result:     value,
						Type:       dieType,
						FancyValue: "",
					}
					dieRolls = append(dieRolls, dieRoll)
					total += value
				}

				rolls = append(rolls, value)
			}
		} else {
			// Roll individual dice normally.
			for _, die := range group.Dice {
				roll := die.Roll()

				var dieType string
				var fancyValue string

				if die.Sides < 0 {
					// This is a fancy die.
					fancyType := fmt.Sprintf("f%d", -die.Sides)
					dieType = fancyType

					if values, exists := fancyDiceValues[fancyType]; exists && roll > 0 && roll <= len(values) {
						fancyValue = values[roll-1].Name // Convert 1-based roll to 0-based index
						total += values[roll-1].Value    // Add the scoring value to total
					}
				} else {
					// Regular die.
					dieType = fmt.Sprintf("d%d", die.Sides)
					fancyValue = ""
					total += roll
				}

				dieRoll := DieRoll{
					Die:        die,
					Result:     roll,
					Type:       dieType,
					FancyValue: fancyValue,
				}
				dieRolls = append(dieRolls, dieRoll)
				rolls = append(rolls, roll)
			}
		}
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

// parseSingleDiceGroup parses a single dice group like "3d6", "d20", "2f4", or "3D6" (exclusive).
func parseSingleDiceGroup(group string) ([]Die, error) {
	group = strings.TrimSpace(group)
	if group == "" {
		return nil, fmt.Errorf("empty dice group")
	}

	// Check for exclusive fancy dice notation first: [count]F[type]
	exclusiveFancyRe := regexp.MustCompile(`^(\d*)F(\d+)$`)
	if matches := exclusiveFancyRe.FindStringSubmatch(group); matches != nil {
		return parseExclusiveFancyDice(matches[1], matches[2])
	}

	// Check for exclusive regular dice notation: [count]D[sides]
	exclusiveRegularRe := regexp.MustCompile(`^(\d*)D(\d+)$`)
	if matches := exclusiveRegularRe.FindStringSubmatch(group); matches != nil {
		return parseExclusiveRegularDice(matches[1], matches[2])
	}

	// Check for fancy dice notation: [count]f[type]
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

// parseExclusiveRegularDice parses exclusive regular dice notation (e.g., "3D6").
func parseExclusiveRegularDice(countStr, sidesStr string) ([]Die, error) {
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			return nil, fmt.Errorf("invalid dice count: %s", countStr)
		}
	}

	sides, err := strconv.Atoi(sidesStr)
	if err != nil || sides <= 0 {
		return nil, fmt.Errorf("invalid dice sides: %s", sidesStr)
	}

	// Validate that we don't request more dice than available faces.
	if count > sides {
		return nil, fmt.Errorf("cannot roll %d exclusive dice with only %d sides", count, sides)
	}

	// Create exclusive dice - encode as positive sides + 1000 to mark as exclusive.
	var dice []Die
	for i := 0; i < count; i++ {
		dice = append(dice, Die{Sides: sides + 1000}) // Mark as exclusive
	}

	return dice, nil
}

// parseExclusiveFancyDice parses exclusive fancy dice notation (e.g., "3F4").
func parseExclusiveFancyDice(countStr, typeStr string) ([]Die, error) {
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			return nil, fmt.Errorf("invalid dice count: %s", countStr)
		}
	}

	fancyType := "f" + typeStr
	values, exists := fancyDiceValues[fancyType]
	if !exists {
		return nil, fmt.Errorf("unsupported fancy dice type: %s", fancyType)
	}

	// Validate that we don't request more dice than available values.
	if count > len(values) {
		return nil, fmt.Errorf("cannot roll %d exclusive %s dice with only %d values", count, fancyType, len(values))
	}

	// Create exclusive fancy dice - encode as negative type - 1000 to mark as exclusive.
	fancyTypeNum, _ := strconv.Atoi(typeStr)
	var dice []Die
	for i := 0; i < count; i++ {
		dice = append(dice, Die{Sides: -fancyTypeNum - 1000}) // Mark as exclusive fancy
	}

	return dice, nil
}

// selectWithoutReplacement selects N unique values from the range [1, K] using shuffle algorithm.
// This is the recursive function you described - picks one at random, swaps with first, reduces slice.
func selectWithoutReplacement(k, n int) []int {
	if n <= 0 || k <= 0 || n > k {
		return nil
	}

	// Create array of K numbers [1, 2, 3, ..., K].
	values := make([]int, k)
	for i := 0; i < k; i++ {
		values[i] = i + 1
	}

	// Select N values using shuffle algorithm.
	return selectFromSlice(values, n)
}

// selectFromSlice recursively selects n values from the slice without replacement.
func selectFromSlice(values []int, n int) []int {
	if n <= 0 || len(values) == 0 {
		return nil
	}

	// Base case: if we only need 1 value, pick one at random.
	if n == 1 {
		randomIndex := rand.IntN(len(values))
		return []int{values[randomIndex]}
	}

	// Pick a random index from the current slice.
	randomIndex := rand.IntN(len(values))

	// Swap the selected value with the first position.
	values[0], values[randomIndex] = values[randomIndex], values[0]

	// Take the first value and recursively select n-1 from the rest.
	selected := []int{values[0]}
	remaining := selectFromSlice(values[1:], n-1)

	return append(selected, remaining...)
}

// ExclusiveGroup represents a group of dice that should be rolled exclusively.
type ExclusiveGroup struct {
	Dice        []Die
	IsExclusive bool
	IsFancy     bool
}

// groupExclusiveDice groups dice by their exclusive nature.
func (ds DiceSet) groupExclusiveDice() []ExclusiveGroup {
	var groups []ExclusiveGroup
	currentGroup := ExclusiveGroup{}

	for _, die := range ds.Dice {
		// Check if this die is exclusive.
		isExclusive := false
		isFancy := false

		if die.Sides > 1000 {
			// Exclusive regular dice.
			isExclusive = true
			isFancy = false
		} else if die.Sides < -1000 {
			// Exclusive fancy dice.
			isExclusive = true
			isFancy = true
		}

		// If this die matches the current group type, add it.
		if len(currentGroup.Dice) == 0 ||
			(currentGroup.IsExclusive == isExclusive && currentGroup.IsFancy == isFancy) {
			currentGroup.Dice = append(currentGroup.Dice, die)
			currentGroup.IsExclusive = isExclusive
			currentGroup.IsFancy = isFancy
		} else {
			// Different type, finish current group and start new one.
			if len(currentGroup.Dice) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = ExclusiveGroup{
				Dice:        []Die{die},
				IsExclusive: isExclusive,
				IsFancy:     isFancy,
			}
		}
	}

	// Add the last group if it has dice.
	if len(currentGroup.Dice) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// rollExclusiveGroup rolls a group of exclusive dice without replacement.
func (ds DiceSet) rollExclusiveGroup(group ExclusiveGroup) []int {
	if !group.IsExclusive || len(group.Dice) == 0 {
		return nil
	}

	if group.IsFancy {
		// Exclusive fancy dice.
		firstDie := group.Dice[0]
		originalType := -(firstDie.Sides + 1000)
		fancyType := fmt.Sprintf("f%d", originalType)

		if values, exists := fancyDiceValues[fancyType]; exists {
			// Use shuffle algorithm to select without replacement.
			indices := selectWithoutReplacement(len(values), len(group.Dice))
			results := make([]int, len(indices))
			for i, index := range indices {
				results[i] = index // Return 1-based indices
			}
			return results
		}

		// Fallback for unknown fancy dice.
		results := make([]int, len(group.Dice))
		for i := range results {
			results[i] = originalType
		}
		return results
	} else {
		// Exclusive regular dice.
		firstDie := group.Dice[0]
		originalSides := firstDie.Sides - 1000

		// Use shuffle algorithm to select without replacement.
		return selectWithoutReplacement(originalSides, len(group.Dice))
	}
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
