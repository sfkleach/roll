// Package gui provides the graphical user interface for the dice rolling application.
package gui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/sfkleach/roll/internal/dice"
)

// App represents the main application window and its components.
type App struct {
	window      fyne.Window
	diceEntry   *widget.Entry
	rollButton  *widget.Button
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
	a.diceEntry.SetPlaceHolder("e.g. 2d6, -a 3d6, --descending 2d20")
	// No default text - starts empty so placeholder is visible.

	// Create roll button.
	a.rollButton = widget.NewButton("Roll Dice", a.onRollButtonClicked)
	a.rollButton.Importance = widget.HighImportance

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
	inputContainer := container.NewBorder(nil, nil, nil, a.rollButton, a.diceEntry)

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
			// For fancy dice, use Canvas text with large font size
			canvasText := canvas.NewText(dieRoll.FancyValue, nil)
			canvasText.TextSize = 24 // Much larger font size
			canvasText.TextStyle = fyne.TextStyle{
				Monospace: true,
				Bold:      true,
			}
			canvasText.Alignment = fyne.TextAlignTrailing
			canvasText.Resize(fyne.NewSize(60, 40))
			gridContent = append(gridContent, diceType, canvasText)
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
