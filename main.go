// Package main provides the entry point for the Roll virtual dice rolling application.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/chzyer/readline"

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
	var fancyFiles = flag.String("fancy", "", "Load custom fancy dice from files matching glob pattern")
	var interactive = flag.Bool("interactive", false, "Run in interactive mode")
	flag.BoolVar(interactive, "i", false, "Run in interactive mode (short form)")
	flag.Parse()

	// Handle version flag.
	if *showVersion {
		fmt.Printf("Roll Dice Application v%s\n", info.GetVersion())
		os.Exit(0)
	}

	// Handle help flag.
	if *showHelp {
		fmt.Printf("Usage: %s [OPTIONS] [DICE_NOTATION]\n\n", os.Args[0])
		fmt.Println("Examples:")
		fmt.Println("  roll 3d6")
		fmt.Println("  roll --ascending 2d10 d6")
		fmt.Println("  roll --fancy='*.dice' 2f6")
		fmt.Println("  roll --interactive")
		fmt.Println()
		fmt.Println(info.GetCheatsheetContent())
		os.Exit(0)
	}

	// Load custom fancy dice files if specified.
	if *fancyFiles != "" {
		err := dice.LoadCustomFancyDice(*fancyFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading fancy dice files: %v\n", err)
			os.Exit(1)
		}
	}

	// Get remaining arguments (dice expressions).
	args := flag.Args()

	// Handle interactive mode.
	if *interactive {
		runInteractive(*ascending, *descending)
		return
	}

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

// runInteractive starts an interactive REPL for dice rolling.
func runInteractive(ascending, descending bool) {
	// Validate sorting flags.
	if ascending && descending {
		fmt.Fprintf(os.Stderr, "Error: Cannot specify both --ascending and --descending flags\n")
		os.Exit(1)
	}

	// Configure readline with better settings.
	config := &readline.Config{
		Prompt:            "roll> ",
		HistoryFile:       "", // No history file for now, could be added later
		AutoComplete:      createAutoCompleter(),
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}

	// Create readline instance.
	rl, err := readline.NewEx(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Printf("Roll Dice Interactive Mode v%s\n", info.GetVersion())
	fmt.Println("Enter dice expressions (e.g., 3d6, 2d10 d6) or 'help' for commands.")
	fmt.Println("Type 'quit' or 'exit' to exit, or press Ctrl+C.")
	fmt.Println()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Handle Ctrl+C gracefully.
				fmt.Println("\nGoodbye!")
				break
			} else if err == io.EOF {
				// Handle Ctrl+D gracefully.
				fmt.Println("\nGoodbye!")
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		// Trim whitespace from input.
		line = strings.TrimSpace(line)

		// Skip empty lines.
		if line == "" {
			continue
		}

		// Handle special commands.
		switch strings.ToLower(line) {
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			printInteractiveHelp()
			continue
		case "version":
			fmt.Printf("Roll Dice Application v%s\n", info.GetVersion())
			continue
		case "cheat", "cheatsheet":
			fmt.Println(info.GetCheatsheetContent())
			continue
		}

		// Process dice expression.
		processDiceExpression(line, ascending, descending)
	}
}

// createAutoCompleter creates an autocompleter for the readline interface.
func createAutoCompleter() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("version"),
		readline.PcItem("cheat"),
		readline.PcItem("cheatsheet"),
		readline.PcItem("quit"),
		readline.PcItem("exit"),
		// Common dice expressions
		readline.PcItem("d4"),
		readline.PcItem("d6"),
		readline.PcItem("d8"),
		readline.PcItem("d10"),
		readline.PcItem("d12"),
		readline.PcItem("d20"),
		readline.PcItem("d100"),
		readline.PcItem("2d6"),
		readline.PcItem("3d6"),
		readline.PcItem("4d6"),
		readline.PcItem("1d20"),
		readline.PcItem("2d10"),
		// Fancy dice
		readline.PcItem("f2"),
		readline.PcItem("f4"),
		readline.PcItem("f6"),
		readline.PcItem("f7"),
		readline.PcItem("f12"),
		readline.PcItem("f13"),
		readline.PcItem("f52"),
		// Exclusive dice
		readline.PcItem("3D6"),
		readline.PcItem("4D6"),
		readline.PcItem("5D6"),
		readline.PcItem("2D10"),
		readline.PcItem("3D10"),
	)
}

// printInteractiveHelp prints help information for interactive mode.
func printInteractiveHelp() {
	fmt.Println("Interactive Mode Commands:")
	fmt.Println("  help           - Show this help")
	fmt.Println("  version        - Show version information")
	fmt.Println("  cheat          - Show dice notation cheatsheet")
	fmt.Println("  quit, exit     - Exit interactive mode")
	fmt.Println()
	fmt.Println("Dice Expression Examples:")
	fmt.Println("  3d6            - Roll three six-sided dice")
	fmt.Println("  2d10 d6        - Roll two ten-sided dice and one six-sided die")
	fmt.Println("  1d20,7d4       - Roll one twenty-sided die and seven four-sided dice")
	fmt.Println("  f2             - Roll a two-sided fancy die (heads/tails)")
	fmt.Println("  3D6            - Roll three exclusive six-sided dice (no repeats)")
	fmt.Println()
}

// processDiceExpression parses and executes a dice expression.
func processDiceExpression(expression string, ascending, descending bool) {
	// Parse the dice notation.
	diceSet, err := dice.ParseDiceNotation(expression)
	if err != nil {
		fmt.Printf("Error parsing dice notation '%s': %v\n", expression, err)
		return
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
