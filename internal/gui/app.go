// Package gui provides the graphical user interface for the dice rolling application.
package gui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/sfkleach/roll/internal/dice"
	"github.com/sfkleach/roll/internal/info"
)

// hasReplacementCharacters checks if a string contains actual replacement characters
// that indicate the font doesn't support the Unicode characters.
func hasReplacementCharacters(text string) bool {
	for _, r := range text {
		if r == '\uFFFD' || // Unicode replacement character �
			r == '\u25A1' || // White square □
			r == '\u2610' || // Ballot box ☐
			r == '\u25AF' || // White vertical rectangle ▯
			r == '\u25AD' || // White rectangle ▭
			r == '?' || // Question mark fallback
			r == '\u003F' || // Another question mark representation
			(r >= '\u2680' && r <= '\u2685') { // Dice face range - often show as replacement
			return true
		}
	}
	return false
}

// App represents the main application window and its components.
type App struct {
	window      fyne.Window
	diceEntry   *widget.Entry
	rollButton  *widget.Button
	infoButton  *widget.Button
	resultsCard *widget.Card
	totalCard   *widget.Card
}

// NewApp creates a new GUI application instance.
func NewApp(window fyne.Window) *App {
	app := &App{
		window: window,
	}
	app.setupUI()
	return app
}

// setupUI initializes the user interface components.
func (a *App) setupUI() {
	// Create input field for dice notation.
	a.diceEntry = widget.NewEntry()
	a.diceEntry.SetPlaceHolder("e.g. 2d6")
	// No default text - starts empty so placeholder is visible.

	// Create roll button.
	a.rollButton = widget.NewButton("Roll Dice", a.onRollButtonClicked)
	a.rollButton.Importance = widget.HighImportance

	// Create info button with theme icon.
	a.infoButton = widget.NewButtonWithIcon("", theme.InfoIcon(), a.onInfoButtonClicked)

	// Create results card (will be populated when rolling).
	a.resultsCard = widget.NewCard("", "", container.NewVBox(
		widget.NewLabel("Click 'Roll Dice' to get started!"),
	))

	// Create total card (will be populated when rolling).
	a.totalCard = widget.NewCard("", "", container.NewVBox(
		widget.NewLabel(""),
	))

	// Allow Enter key to trigger roll.
	a.diceEntry.OnSubmitted = func(string) {
		a.onRollButtonClicked()
	}

	// Create layout.
	buttonsContainer := container.NewHBox(a.infoButton, a.rollButton)
	inputContainer := container.NewBorder(nil, nil, nil, buttonsContainer, a.diceEntry)

	content := container.NewVBox(
		inputContainer,
		widget.NewSeparator(),
		a.resultsCard,
		a.totalCard,
	)

	a.window.SetContent(content)
}

// parseFlagsFromInput extracts sorting flags from the input text and returns cleaned dice notation and sorting preferences.
func parseFlagsFromInput(input string) (diceNotation string, ascending bool, descending bool, err error) {
	parts := strings.Fields(input)
	var cleanParts []string

	for _, part := range parts {
		switch part {
		case "-a", "--ascending":
			if descending {
				return "", false, false, fmt.Errorf("cannot specify both ascending and descending flags")
			}
			ascending = true
		case "-d", "--descending":
			if ascending {
				return "", false, false, fmt.Errorf("cannot specify both ascending and descending flags")
			}
			descending = true
		default:
			cleanParts = append(cleanParts, part)
		}
	}

	diceNotation = strings.Join(cleanParts, " ")
	return diceNotation, ascending, descending, nil
}

// onRollButtonClicked handles the roll button click event.
func (a *App) onRollButtonClicked() {
	input := strings.TrimSpace(a.diceEntry.Text)

	if input == "" {
		a.showError("Please enter dice notation (e.g. 2d6, -a 3d6, --descending 2d20)")
		return
	}

	// Parse flags from input.
	notation, ascending, descending, err := parseFlagsFromInput(input)
	if err != nil {
		a.showError(fmt.Sprintf("Flag error: %v", err))
		return
	}

	if notation == "" {
		a.showError("Please enter dice notation after any flags")
		return
	}

	// Parse the dice notation.
	diceSet, err := dice.ParseDiceNotation(notation)
	if err != nil {
		a.showError(fmt.Sprintf("Invalid dice notation: %v", err))
		return
	}

	// Roll the dice.
	result := diceSet.Roll()

	// Sort if requested.
	if ascending || descending {
		sortedRolls := make([]dice.DieRoll, len(result.DieRolls))
		copy(sortedRolls, result.DieRolls)

		if ascending {
			sort.Slice(sortedRolls, func(i, j int) bool {
				return sortedRolls[i].Result < sortedRolls[j].Result
			})
		} else if descending {
			sort.Slice(sortedRolls, func(i, j int) bool {
				return sortedRolls[i].Result > sortedRolls[j].Result
			})
		}

		// Create a new result with sorted rolls.
		sortedResult := dice.RollResult{
			DieRolls:        sortedRolls,
			IndividualRolls: result.IndividualRolls, // Keep original for compatibility.
			Total:           result.Total,
		}
		a.updateResults(sortedResult)
	} else {
		// Update the display with original order.
		a.updateResults(result)
	}
}

// updateResults updates the result display with separate areas for dice rolls and total.
func (a *App) updateResults(result dice.RollResult) {
	// Create the dice results grid (pre-allocate with capacity for die rolls).
	gridContent := make([]fyne.CanvasObject, 0, len(result.DieRolls)*2)

	// Add each individual die roll as a row in the grid.
	for _, dieRoll := range result.DieRolls {
		// Left column: dice type (e.g., "d6", "d20", "f4", "f12").
		diceType := widget.NewLabel(dieRoll.Type)
		diceType.Alignment = fyne.TextAlignLeading

		// For fancy dice, use monospace font for consistency.
		if dieRoll.FancyValue != "" {
			diceType.TextStyle = fyne.TextStyle{Monospace: true}
		}

		// Right column: roll result (fancy value or numeric).
		if dieRoll.FancyValue != "" {
			// For fancy dice, check if Unicode characters render as replacement characters
			displayText := dieRoll.FancyValue
			if hasReplacementCharacters(dieRoll.FancyValue) {
				// Fall back to showing the score if Unicode shows replacement characters
				displayText = fmt.Sprintf("%d", dieRoll.Result)
			}

			rollValue := widget.NewLabel(displayText)
			rollValue.Alignment = fyne.TextAlignTrailing
			// No special TextStyle to allow system font with natural colors
			gridContent = append(gridContent, diceType, rollValue)
		} else {
			// Regular numeric value
			rollValue := widget.NewLabel(fmt.Sprintf("%d", dieRoll.Result))
			rollValue.Alignment = fyne.TextAlignTrailing
			gridContent = append(gridContent, diceType, rollValue)
		}
	}

	// Create a 2-column grid for dice results.
	diceGrid := container.NewGridWithColumns(2, gridContent...)

	// Update the results card content.
	a.resultsCard.SetContent(diceGrid)

	// Create total display.
	totalLabel := widget.NewLabel(fmt.Sprintf("Total: %d", result.Total))
	totalLabel.Alignment = fyne.TextAlignCenter
	totalLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Update the total card content.
	a.totalCard.SetContent(totalLabel)
}

// showError displays an error message to the user.
func (a *App) showError(message string) {
	errorLabel := widget.NewLabel(message)
	errorLabel.Wrapping = fyne.TextWrapWord
	a.resultsCard.SetContent(errorLabel)

	// Clear the total area.
	a.totalCard.SetContent(widget.NewLabel(""))
}

// onInfoButtonClicked shows information about dice notation and sorting options in a separate window.
func (a *App) onInfoButtonClicked() {
	// Create a new window for the cheatsheet.
	cheatWindow := fyne.CurrentApp().NewWindow("Dice Rolling Cheatsheet")
	cheatWindow.Resize(fyne.NewSize(600, 500))

	// Get the unified cheatsheet content in markdown format.
	cheatContent := info.GetCheatsheetMarkdown()

	// Create a rich text widget with the cheatsheet content.
	richText := widget.NewRichTextFromMarkdown(cheatContent)
	richText.Wrapping = fyne.TextWrapWord

	// Create close button.
	closeBtn := widget.NewButton("Close", func() {
		cheatWindow.Close()
	})

	// Create scroll container for the content.
	scroll := container.NewScroll(richText)

	// Layout the window.
	content := container.NewBorder(
		nil,      // top
		closeBtn, // bottom
		nil,      // left
		nil,      // right
		scroll,   // center
	)

	cheatWindow.SetContent(content)
	cheatWindow.Show()
}
