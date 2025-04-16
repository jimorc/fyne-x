package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Spinner Demo")
	intLabel := widget.NewLabel("Int Spinner (1, 10, 1):")
	intLabel.Alignment = fyne.TextAlignTrailing
	intSpinner := xwidget.NewSpinner(1, 10, 1, 0)
	intC := container.NewGridWithColumns(2, intLabel, intSpinner)

	f2Label := widget.NewLabel("Float Spinner (-2, 25, 1.5, 2):")
	f2Label.Alignment = fyne.TextAlignTrailing
	f2Spinner := xwidget.NewSpinner(-2, 25, 1.5, 2)
	f2C := container.NewGridWithColumns(2, f2Label, f2Spinner)
	c := container.NewVBox(intC, f2C)
	w.SetContent(c)
	w.ShowAndRun()
}
