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
	l := widget.NewLabel("Int Spinner:")
	l.Alignment = fyne.TextAlignTrailing
	s := xwidget.NewSpinner()
	c := container.NewGridWithColumns(2, l, s)
	w.SetContent(c)
	w.ShowAndRun()
}
