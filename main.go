package main

import (
	"simulation/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Parking")
	window.CenterOnScreen()
	view.NewParkingView(window)
	window.ShowAndRun()
}
