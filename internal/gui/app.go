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

// SortOrder represents the different ways to sort dice results.
type SortOrder int

const (
	SortDiceOrder SortOrder = iota // Original order as entered by user.
	SortAscending                  // Sort by roll value, lowest first.
	SortDescending                 // Sort by roll value, highest first.
)

// App represents the main application window and its components.
type App struct {
	window      fyne.Window
	diceEntry   *widget.Entry
	rollButton  *widget.Button
	resultsCard *widget.Card
	totalCard   *widget.Card
	currentSort SortOrder
	lastResult  *dice.RollResult // Store last result for re-sorting.
}

// NewApp creates a new GUI application instance.
func NewApp(window fyne.Window) *App {
	app := &App{
		window:      window,
		currentSort: SortDiceOrder, // Default to dice order.
	}
	app.setupUI()
	app.setupMenu()
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

// setupMenu creates the application menubar.
func (a *App) setupMenu() {
	// Create View menu with reset function and sort options.
	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Reset Display", func() {
			a.resetDisplay()
		}),
		fyne.NewMenuItemSeparator(),
		// Sort submenu.
		fyne.NewMenuItem("Sort by Dice Order", func() {
			a.setSortOrder(SortDiceOrder)
		}),
		fyne.NewMenuItem("Sort Ascending", func() {
			a.setSortOrder(SortAscending)
		}),
		fyne.NewMenuItem("Sort Descending", func() {
			a.setSortOrder(SortDescending)
		}),
	)

	// Create File menu with standard items.
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exit", func() {
			a.window.Close()
		}),
	)

	// Set the main menu on the window.
	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu)
	a.window.SetMainMenu(mainMenu)
}

// setSortOrder changes the sort order and re-displays the current results.
func (a *App) setSortOrder(order SortOrder) {
	a.currentSort = order
	
	// Re-display the last result with the new sort order if available.
	if a.lastResult != nil {
		a.updateResults(*a.lastResult)
	}
}

// setupHotkeys configures keyboard shortcuts for the application.
func (a *App) setupHotkeys() {
	// Add Ctrl+R hotkey to reset the display.
	a.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyR {
			a.resetDisplay()
		}
	})
}

// resetDisplay clears all results and resets the display to initial state.
func (a *App) resetDisplay() {
	// Reset results card to initial state.
	a.resultsCard.SetContent(container.NewVBox(
		widget.NewLabel("Click 'Roll Dice' to get started!"),
	))

	// Clear the total card.
	a.totalCard.SetContent(widget.NewLabel(""))
	
	// Clear the last result.
	a.lastResult = nil
}

// onRollButtonClicked handles the roll button click event.
func (a *App) onRollButtonClicked() {
	notation := strings.TrimSpace(a.diceEntry.Text)

	if notation == "" {
		a.showError("Please enter dice notation (e.g. 2d6)")
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

// updateResults updates the result display with separate areas for dice rolls and total.
func (a *App) updateResults(result dice.RollResult) {
	// Store the result for potential re-sorting.
	a.lastResult = &result

	// Create a copy of the die rolls for sorting.
	dieRolls := make([]dice.DieRoll, len(result.DieRolls))
	copy(dieRolls, result.DieRolls)

	// Apply sorting based on current sort order.
	switch a.currentSort {
	case SortAscending:
		sort.Slice(dieRolls, func(i, j int) bool {
			return dieRolls[i].Result < dieRolls[j].Result
		})
	case SortDescending:
		sort.Slice(dieRolls, func(i, j int) bool {
			return dieRolls[i].Result > dieRolls[j].Result
		})
	case SortDiceOrder:
		// No sorting needed - use original order.
	}

	// Create the dice results grid (pre-allocate with capacity for die rolls).
	gridContent := make([]fyne.CanvasObject, 0, len(dieRolls)*2)

	// Add each individual die roll as a row in the grid.
	for _, dieRoll := range dieRolls {
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
	
	// Clear the last result since we're showing an error.
	a.lastResult = nil
}
