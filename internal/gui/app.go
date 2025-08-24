// Package gui provides the graphical user interface for the dice rolling application.
package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/sfkleach/roll/internal/dice"
)

// App represents the main application window and its components.
type App struct {
	window      fyne.Window
	diceEntry   *widget.Entry
	rollButton  *widget.Button
	resultLabel *widget.Label
	detailLabel *widget.Label
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
	a.diceEntry.SetPlaceHolder("Enter dice notation (e.g., 3d6, d20, 2d10 d6)")
	a.diceEntry.Text = "3d6" // Default value for convenience.

	// Create roll button.
	a.rollButton = widget.NewButton("Roll Dice", a.onRollButtonClicked)
	a.rollButton.Importance = widget.HighImportance

	// Create result labels.
	a.resultLabel = widget.NewLabel("Total: -")
	a.resultLabel.TextStyle = fyne.TextStyle{Bold: true}
	a.detailLabel = widget.NewLabel("Click 'Roll Dice' to get started!")
	a.detailLabel.Wrapping = fyne.TextWrapWord

	// Allow Enter key to trigger roll.
	a.diceEntry.OnSubmitted = func(string) {
		a.onRollButtonClicked()
	}

	// Create layout.
	inputContainer := container.NewBorder(nil, nil, nil, a.rollButton, a.diceEntry)

	content := container.NewVBox(
		widget.NewCard("Dice Notation", "", inputContainer),
		widget.NewSeparator(),
		widget.NewCard("Results", "", container.NewVBox(
			a.resultLabel,
			a.detailLabel,
		)),
	)

	a.window.SetContent(content)
}

// onRollButtonClicked handles the roll button click event.
func (a *App) onRollButtonClicked() {
	notation := strings.TrimSpace(a.diceEntry.Text)

	if notation == "" {
		a.showError("Please enter dice notation (e.g., 3d6, d20, 2d10 d6)")
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

	// Update the display.
	a.updateResults(result)
}

// updateResults updates the result display with the roll results.
func (a *App) updateResults(result dice.RollResult) {
	// Update total.
	a.resultLabel.SetText(fmt.Sprintf("Total: %d", result.Total))

	// Update individual roll details.
	var details strings.Builder
	details.WriteString("Individual rolls: ")

	for i, roll := range result.IndividualRolls {
		if i > 0 {
			details.WriteString(", ")
		}
		details.WriteString(fmt.Sprintf("%d", roll))
	}

	a.detailLabel.SetText(details.String())
}

// showError displays an error message to the user.
func (a *App) showError(message string) {
	a.resultLabel.SetText("Error")
	a.detailLabel.SetText(message)
}
