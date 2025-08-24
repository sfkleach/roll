// Package main provides the entry point for the Roll virtual dice rolling application.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/sfkleach/roll/internal/gui"
)

func main() {
	myApp := app.NewWithID("com.github.sfkleach.roll")

	myWindow := myApp.NewWindow("Roll - Virtual Dice")
	myWindow.Resize(fyne.NewSize(450, 350))
	myWindow.CenterOnScreen()

	// Create and setup the GUI.
	gui.NewApp(myWindow)

	myWindow.ShowAndRun()
}
