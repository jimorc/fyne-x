package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	spinner := xwidget.NewSpinner(-2002, 300, 200)
	c := container.NewCenter(spinner)
	w := a.NewWindow("Spinner Demo")
	w.SetContent(c)
	w.Resize(fyne.NewSize(100, 100))
	w.ShowAndRun()
}
