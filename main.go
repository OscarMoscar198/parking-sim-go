package main

import (
	"simulation/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Parkinglot Simulator")
	window.CenterOnScreen()
	view.NewView(window)
	window.ShowAndRun()
}
