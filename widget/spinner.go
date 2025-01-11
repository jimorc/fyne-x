package widget

import (
	"image/color"
	"strconv"
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
	text := canvas.NewText(strconv.Itoa(s.value), th.Color(theme.ColorNameForeground, v))
	text.Alignment = fyne.TextAlignTrailing
	objects := []fyne.CanvasObject{
		box,
		border,
		text,
	}
	r := &spinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}

	return r
}

func (s *Spinner) textSize() fyne.Size {
	minText := canvas.NewText(strconv.Itoa(s.min), color.Black)
	maxText := canvas.NewText(strconv.Itoa(s.max), color.Black)
	minTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(minText.Text,
		minText.TextSize, minText.TextStyle, minText.FontSource)
	maxTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(maxText.Text,
		maxText.TextSize, maxText.TextStyle, maxText.FontSource)
	return fyne.NewSize(max(minTextSize.Width, maxTextSize.Width),
		max(minTextSize.Height, maxTextSize.Height))
}

type spinnerRenderer struct {
	box    *canvas.Rectangle
	border *canvas.Rectangle
	text   *canvas.Text

	objects []fyne.CanvasObject
	spinner *Spinner
}

// Destroy destroys any objects that must be destroyed when the renderer is destroyed.
func (r *spinnerRenderer) Destroy() {}

func (r *spinnerRenderer) Layout(size fyne.Size) {
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	padding := th.Size(theme.SizeNameInnerPadding)

	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	r.border.Resize(fyne.NewSize(size.Width-borderSize-.5, size.Height-borderSize-.5))
	r.border.StrokeWidth = borderSize
	r.border.Move(fyne.NewSquareOffsetPos(borderSize / 2))
	r.box.Resize(size.Subtract(fyne.NewSquareSize(borderSize * 2)))
	r.box.Move(fyne.NewSquareOffsetPos(borderSize))

	textSize := r.spinner.textSize()
	rMinSize := r.MinSize()
	xPos := borderSize + padding + textSize.Width
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))
}

// MinSize calculates the minimum size of the Spinner widget.
// This size is based on the maximum width of the number text to be displayed in the widget,
// and the size of the two spinner buttons.
func (r *spinnerRenderer) MinSize() fyne.Size {
	// TODO: add button sizes
	th := r.spinner.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := r.spinner.textSize()
	tHeight := textSize.Height + padding.Height
	tWidth := textSize.Width + padding.Width
	return fyne.NewSize(tWidth, tHeight)
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
	r.text.Text = strconv.Itoa(r.spinner.value)
}

// max returns the larger of the two arguments.
// When fyne-x is updated to use a go version of 1.21 or greater,
// this function can be deleted and replaced by one of the functions
// in the newer version of go.
func max(a, b float32) float32 {
	max := a
	if a < b {
		max = b
	}
	return max
}
