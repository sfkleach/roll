// Package main provides the entry point for the Roll virtual dice rolling application.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/sfkleach/roll/internal/dice"
	"github.com/sfkleach/roll/internal/gui"
	"github.com/sfkleach/roll/internal/info"
)

func main() {
	// Define command line flags with abbreviated versions.
	var ascending = flag.Bool("ascending", false, "Sort individual dice rolls in ascending order")
	flag.BoolVar(ascending, "a", false, "Sort individual dice rolls in ascending order (short form)")
	var descending = flag.Bool("descending", false, "Sort individual dice rolls in descending order")
	flag.BoolVar(descending, "d", false, "Sort individual dice rolls in descending order (short form)")
	var showHelp = flag.Bool("help", false, "Show help and cheatsheet")
	var showVersion = flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag.
	if *showVersion {
		fmt.Printf("Roll Dice Application v%s\n", info.GetVersion())
		os.Exit(0)
	}

	// Handle help flag.
	if *showHelp {
		fmt.Println(info.GetCheatsheetContent())
		os.Exit(0)
	}

	// Get remaining arguments (dice expressions).
	args := flag.Args()

	// If command line arguments are provided, run in command line mode.
	if len(args) > 0 {
		runCommandLine(args, *ascending, *descending)
		return
	}

	// Otherwise, run the GUI application.
	runGUI()
}

// runCommandLine processes dice expressions from command line arguments.
func runCommandLine(diceExpressions []string, ascending, descending bool) {
	// Validate sorting flags.
	if ascending && descending {
		fmt.Fprintf(os.Stderr, "Error: Cannot specify both --ascending and --descending flags\n")
		os.Exit(1)
	}

	// Join all arguments into a single dice expression.
	expression := strings.Join(diceExpressions, " ")

	// Parse the dice notation.
	diceSet, err := dice.ParseDiceNotation(expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing dice notation '%s': %v\n", expression, err)
		os.Exit(1)
	}

	// Roll the dice.
	result := diceSet.Roll()

	// Sort individual rolls if requested.
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

		// Print sorted results.
		printCommandLineResults(sortedRolls, result.Total)
	} else {
		// Print results in original order.
		printCommandLineResults(result.DieRolls, result.Total)
	}
}

// printCommandLineResults prints the dice roll results to stdout.
func printCommandLineResults(dieRolls []dice.DieRoll, total int) {
	for _, roll := range dieRolls {
		if roll.FancyValue != "" {
			// For fancy dice, show the fancy value.
			fmt.Printf("%s: %s\n", roll.Type, roll.FancyValue)
		} else {
			// For regular dice, show the numeric result.
			fmt.Printf("%s: %d\n", roll.Type, roll.Result)
		}
	}
	fmt.Printf("Total: %d\n", total)
}

// runGUI starts the graphical user interface.
func runGUI() {
	myApp := app.NewWithID("com.github.sfkleach.roll")

	myWindow := myApp.NewWindow("Roll - Virtual Dice")
	myWindow.Resize(fyne.NewSize(450, 350))
	myWindow.CenterOnScreen()

	// Create and setup the GUI.
	gui.NewApp(myWindow)

	myWindow.ShowAndRun()
}
