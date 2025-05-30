package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Disableable = (*Spinner)(nil)
var _ fyne.Focusable = (*Spinner)(nil)

var _ desktop.Mouseable = (*Spinner)(nil)
var _ fyne.Scrollable = (*Spinner)(nil)
var _ Spinnable = (*Spinner)(nil)

// Spinner widget has a minimum, maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type Spinner struct {
	widget.DisableableWidget
	base *SpinnerBase

	hovered bool
	focused bool

	OnChanged func(float64) `json:"-"`
}

// NewSpinner creates a new Spinner widget.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 	decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
//
//	onChanged is the callback function that is called whenever the spinner value changes.
func NewSpinner(min, max, step float64, decPlaces uint, onChanged func(float64)) *Spinner {
	s := &Spinner{OnChanged: onChanged}
	s.base = NewSpinnerBase(s, min, max, step, decPlaces)
	return s
}

// NewSpinnerUninitialized returns a new uninitialized Spinner widget.
//
// An uninitialized Spinner widget is useful when you need to create a Spinner
// but the initial settings are unknown.
// Calling Enable on an uninitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize spinner values before enabling
// the spinner widget.
//
// Params:
//
//	decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
func NewSpinnerUninitialized(decPlaces uint) *Spinner {
	s := &Spinner{}
	s.base = NewSpinnerBaseUninitialized(s, decPlaces)
	s.ExtendBaseWidget(s)
	s.Disable()
	return s
}

// NewSpinnerWithData returns a new Spinner widget connected to the specified data source.
//
// Params:
//
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 	decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
//
//	data is the value that is bound to the spinner value.
func NewSpinnerWithData(min, max, step float64, decPlaces uint, data binding.Float) *Spinner {
	s := &Spinner{}
	s.base = NewSpinnerBaseWithData(s, min, max, step, decPlaces, data)

	return s
}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (s *Spinner) Bind(data binding.Float) {
	s.base.Bind(data)
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer.
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	box := canvas.NewRectangle(th.Color(theme.ColorNameBackground, v))
	border := canvas.NewRectangle(color.Transparent)

	text := canvas.NewText(s.base.ValueText(), th.Color(theme.ColorNameForeground, v))
	text.Alignment = fyne.TextAlignTrailing

	objects := []fyne.CanvasObject{
		box,
		border,
		text,
		s.base.UpButton(),
		s.base.DownButton(),
	}
	r := &SpinnerRenderer{
		spinner: s,
		box:     box,
		border:  border,
		text:    text,
		objects: objects,
	}
	return r
}

// Decrement decrements the spinner's value by its step value.
func (s *Spinner) Decrement() {
	s.base.Decrement()
}

// Disable disables the Spinner and its buttons.
func (s *Spinner) Disable() {
	if s.Disabled() {
		return
	}
	if s.base != nil {
		s.base.DownButton().Disable()
		s.base.UpButton().Disable()
	}
	s.DisableableWidget.Disable()
	s.Refresh()
}

// Enable enables the Spinner and its buttons as appropriate.
func (s *Spinner) Enable() {
	if s.base == nil || !s.base.Initialized() {
		return
	}

	s.DisableableWidget.Enable()
	s.SetValue(s.Value())
	s.Refresh()
}

// FocusGained is called when the spinner has been given focus.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called when the spinner has had focus removed.
//
// Implements: fyne.Focusable
func (s *Spinner) FocusLost() {
	s.focused = false
	s.Refresh()
}

// GetOnChanged returns the Spinner's OnChanged function.
//
// Implements the Spinnable interface.
func (s *Spinner) GetOnChanged() func(float64) {
	return func(float64) {
		if s.OnChanged != nil {
			s.OnChanged(s.Value())
		}
		if s.base != nil && s.base.Initialized() {
			s.Refresh()
		}
	}
}

// Increment increments the spinner's value by its step value.
func (s *Spinner) Increment() {
	s.base.Increment()
}

// MinSize returns the minimum size of the Spinner widget. The minimum size is calculated
// based on the maximum width that the value could require based on its format.
func (s *Spinner) MinSize() fyne.Size {
	th := s.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	textSize := s.textSize()
	tHeight := textSize.Height + padding.Height
	upButtonHeight := s.base.UpButton().MinSize().Height
	tWidth := textSize.Width + upButtonHeight + padding.Width
	return fyne.NewSize(tWidth, tHeight)
}

// MouseDown called on mouse click.
// This action causes the Spinner to request focus.
//
// Implements: desktop.Mouseable
func (s *Spinner) MouseDown(m *desktop.MouseEvent) {
	s.requestFocus()
	s.Refresh()
}

// MouseUp called on mouse release.
//
// Implements: desktop.Mouseable
func (s *Spinner) MouseUp(m *desktop.MouseEvent) {}

// Scrolled handles mouse scroller events.
//
// Implements fyne.Scrollable
func (s *Spinner) Scrolled(evt *fyne.ScrollEvent) {
	if s.Disabled() || !s.focused {
		return
	}
	if evt.Scrolled.DY > 0 {
		s.Increment()
	} else if evt.Scrolled.DY < 0 {
		s.Decrement()
	}
}

// SetMinMaxStep sets the widget's minimum, maximum, and step values.
//
// Params:
//
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//
// If the previously set value is less than min, then the value is set to min.
// If the previously set value is greater than max, then the value is set to max.
func (s *Spinner) SetMinMaxStep(min, max, step float64) {
	s.base.SetMinMaxStep(min, max, step)
	s.Refresh()
}

// SetValue sets the spinner value. It ensures that the value is always >= min and
// <= max.
func (s *Spinner) SetValue(val float64) {
	if s.Disabled() {
		return
	}
	s.base.SetValue(val)
	s.Refresh()
}

// TypedKey receives key input events when the spinner widget has focus.
// Increments/decrements the spinner's value when the up or down key is
// pressed.
//
// Implements: fyne.Focusable
func (s *Spinner) TypedKey(key *fyne.KeyEvent) {
	if s.Disabled() || !s.focused {
		return
	}
	switch key.Name {
	case fyne.KeyUp:
		s.Increment()
	case fyne.KeyDown:
		s.Decrement()
	default:
		return
	}
}

// TypedRune receives text input events when the spinner widget is focused.
// Increments/decrements the spinner's value when the '+' or '-' key is
// pressed.
//
// Implements: fyne.Focusable
func (s *Spinner) TypedRune(rune rune) {
	if s.Disabled() || !s.focused {
		return
	}
	switch rune {
	case '+':
		s.Increment()
	case '-':
		s.Decrement()
	default:
		return
	}
}

// Unbind disconnects any configured data source from this spinner.
// The current value will remain at the last value of the data source.
func (s *Spinner) Unbind() {
	s.base.Unbind()
}

// Value retrieves the current Spinner value.
func (s *Spinner) Value() float64 {
	return s.base.Value()
}

// requestFocus requests that this Spinner receive focus.
func (s *Spinner) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(s); c != nil {
		c.Focus(s)
	}

}

// Calculate the max size of the text that can be displayed for the spinner.
// The size cannot be larger than the larger of the sizes required to display the
// spinner's min and max values.
func (s *Spinner) textSize() fyne.Size {
	return maxTextSize(s.base.MinText(), s.base.MaxText())
}

// Validate validates the Spinner widget.
func (s *Spinner) Validate() error {
	return s.base.Validate()
}

// SpinnerRenderer is the renderer for the Spinner widget
type SpinnerRenderer struct {
	spinner *Spinner
	box     *canvas.Rectangle
	border  *canvas.Rectangle
	text    *canvas.Text
	objects []fyne.CanvasObject
}

// Destroy destroys any objects that must be destroyed when the renderer is
// destroyed.
func (r *SpinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Float64Spinner widget.
func (r *SpinnerRenderer) Layout(size fyne.Size) {
	r.spinner.Refresh()
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	padding := th.Size(theme.SizeNameInnerPadding)

	buttonSize := r.spinner.base.UpButton().MinSize()
	newSize := fyne.NewSize(size.Width-buttonSize.Width-padding/2, size.Height)
	// 0.5 is removed so on low DPI it rounds down on the trailing edge
	newSize = fyne.NewSize(newSize.Width-0.5, newSize.Height-0.5)
	topLeft := fyne.NewPos(0, 0)
	r.box.Resize(newSize)
	r.box.Move(topLeft)
	r.border.Resize(newSize)
	r.border.StrokeWidth = borderSize
	r.border.Move(topLeft)

	textSize := r.spinner.textSize()
	rMinSize := r.MinSize()
	// -2 in the line below positions the text correctly. I could not find
	// any specific reason for this value, just that it works.
	xPos := size.Width - buttonSize.Width - padding - 2
	yPos := (rMinSize.Height - textSize.Height) / 2
	r.text.Move(fyne.NewPos(xPos, yPos))

	xPos += padding
	yPos -= padding - 1
	r.spinner.base.UpButton().Resize(buttonSize)
	r.spinner.base.UpButton().Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.base.UpButton().MinSize().Height + padding/2
	r.spinner.base.DownButton().Resize(buttonSize)
	r.spinner.base.DownButton().Move(fyne.NewPos(xPos, yPos))
}

// MinSize returns the minimum size of the Flaot64Spinner widget.
func (r *SpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects associated with the Float64Spinner renderer.
func (r *SpinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh refreshes (redisplays) the Float64Spinner widget.
func (r *SpinnerRenderer) Refresh() {
	th := r.spinner.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	fgColor, bgColor, borderColor := spinnerColors(
		r.spinner.Disabled(), r.spinner.focused, r.spinner.hovered)
	r.box.FillColor = th.Color(bgColor, v)
	r.box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.border.CornerRadius = r.box.CornerRadius
	if r.spinner.Validate() == nil {
		r.border.StrokeColor = th.Color(borderColor, v)
	} else {
		r.border.StrokeColor = th.Color(theme.ColorNameError, v)
	}

	r.text.Text = r.spinner.base.ValueText()
	r.text.Color = th.Color(fgColor, v)
	r.text.Refresh()
}

// maxTextSize calculates the larger of the canvas.Text sizes for the two string params
func maxTextSize(minText, maxText string) fyne.Size {
	// color does not affect the text size, so use Black.
	minT := canvas.NewText(minText, color.Black)
	maxT := canvas.NewText(maxText, color.Black)
	minTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(minT.Text,
		minT.TextSize, minT.TextStyle, minT.FontSource)
	maxTextSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(maxT.Text,
		maxT.TextSize, maxT.TextStyle, maxT.FontSource)
	return fyne.NewSize(max(minTextSize.Width, maxTextSize.Width),
		max(minTextSize.Height, maxTextSize.Height))
}

// max returns the larger of the two arguments.
// This can/should be replaced by the appropriate go max function
// when the version of go used to build fyne-x is updated to version
// 1.21 or later.
func max(a, b float32) float32 {
	max := a
	if a < b {
		max = b
	}
	return max
}

// spinnerColors determines display colors for spinners.
func spinnerColors(disabled, focused, hovered bool) (fgColor, bgColor, borderColor fyne.ThemeColorName) {
	fgColor = theme.ColorNameForeground
	bgColor = ""
	borderColor = theme.ColorNameInputBorder
	if disabled {
		fgColor = theme.ColorNameDisabled
		borderColor = theme.ColorNameDisabled
	} else if focused {
		borderColor = theme.ColorNamePrimary
	} else if hovered {
		bgColor = theme.ColorNameHover
	}
	return fgColor, bgColor, borderColor
}
