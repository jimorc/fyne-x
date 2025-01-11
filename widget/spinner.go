package widget

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Spinner widget has an integer value, with two spinner buttons to increment/decrement the value.
type Spinner struct {
	widget.DisableableWidget

	value int
	min   int
	max   int
	step  int

	propertyLock sync.RWMutex
}

// NewSpinner creates a new Spinner widget with the specified minimum, maximum, and step values.
// The initial value is set to the min value.
func NewSpinner(min, max, step int) *Spinner {
	s := &Spinner{
		min:   min,
		max:   max,
		step:  step,
		value: min,
	}
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	s.ExtendBaseWidget(s)

	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = th.Size(theme.SizeNameInputBorder)
	border.StrokeColor = th.Color(theme.ColorNameInputBorder, v)
	border.CornerRadius = th.Size(theme.SizeNameInputRadius)
	objects := []fyne.CanvasObject{
		box,
		border,
	}
	r := &spinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		objects: objects,
	}

	return r
}

type spinnerRenderer struct {
	box    *canvas.Rectangle
	border *canvas.Rectangle

	objects []fyne.CanvasObject
	spinner *Spinner
}

// Destroy destroys any objects that must be destroyed when the renderer is destroyed.
func (r *spinnerRenderer) Destroy() {}

func (r *spinnerRenderer) Layout(size fyne.Size) {
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)

	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	r.border.Resize(fyne.NewSize(size.Width-borderSize-.5, size.Height-borderSize-.5))
	r.border.StrokeWidth = borderSize
	r.border.Move(fyne.NewSquareOffsetPos(borderSize / 2))
	r.box.Resize(size.Subtract(fyne.NewSquareSize(borderSize * 2)))
	r.box.Move(fyne.NewSquareOffsetPos(borderSize))
}

// MinSize calculates the minimum size of the Spinner widget.
// This size is based on the maximum width of the number text to be displayed in the widget,
// and the size of the two spinner buttons.
func (r *spinnerRenderer) MinSize() fyne.Size {
	// TODO: calculate actual required size
	return fyne.NewSize(40, 30)
}

// Objects returns the objects associated with the spinner renderer.
func (r *spinnerRenderer) Objects() []fyne.CanvasObject {
	r.spinner.propertyLock.RLock()
	defer r.spinner.propertyLock.RUnlock()
	return r.objects
}

// Refresh refreshes (redisplays) the spinner widget.
func (r *spinnerRenderer) Refresh() {
	//	r.spinner.propertyLock.RLock()

	th := r.spinner.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.box.FillColor = th.Color(theme.ColorNameInputBackground, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius

	r.border.StrokeColor = th.Color(theme.ColorNamePrimary, v)
}
