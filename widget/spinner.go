package widget

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// spinnerButton is the widget used for buttons in the Spinner widget.
type spinnerButton struct {
	widget.Button

	spinner *Spinner

	position fyne.Position
	size     fyne.Size
}

// newSpinnerButton creates a spinerButton for use in Spinner widgets.
//
// Params:
//
//	resource is the resource used as the button icon.
//	onTapped is the callback function for button clicks.
func newSpinnerButton(s *Spinner, resource fyne.Resource, onTapped func()) *spinnerButton {
	b := &spinnerButton{spinner: s}
	b.ExtendBaseWidget(b)

	b.setButtonProperties(resource, onTapped)
	return b
}

// MinSize returns the minimum size of the button. Because the minimum size is a constant
// based on the spinner height and theme properties, the minimum size is calculated when
// the button is created.
func (b *spinnerButton) MinSize() fyne.Size {
	return b.size
}

func (b *spinnerButton) Move(pos fyne.Position) {
	b.position = pos
	b.BaseWidget.Move(pos)
}

// setButtonProperties sets the button properties;
//
// Params:
//
//	resource is the resource for the button icon.
//	onTapped is the callback function for button clicks.
func (b *spinnerButton) setButtonProperties(resource fyne.Resource, onTapped func()) {
	b.Icon = resource
	b.OnTapped = onTapped
	b.Text = ""

	// calculate minimum button size (really just its height).
	th := b.Theme()
	tHeight := b.spinner.entry.MinSize().Height
	h := tHeight/2 - th.Size(theme.SizeNameInputBorder) - 1
	b.size = fyne.NewSize(h, h)
}

var _ fyne.Disableable = (*spinnerButton)(nil)
var _ fyne.Focusable = (*spinnerButton)(nil)

// spinnerEntry is the entry widget for the Spinner widget.
type spinnerEntry struct {
	NumericalEntry

	shortcut fyne.ShortcutHandler
	spinner  *Spinner
}

// newSpinnerEntry creates a spinnerEntry widget.
func newSpinnerEntry(s *Spinner) *spinnerEntry {
	e := &spinnerEntry{spinner: s}
	e.registerShortcuts()
	e.Validator = e.validate
	e.ExtendBaseWidget(e)

	return e
}

// MinSize calculates the minimum size required for the spinner entry.
//
// It determines the minimum size based on the minimum and maximum values
// of the spinner, taking into account the number of decimal places.
// If the spinner is not initialized, it returns the minimum size of the
// underlying NumericalEntry.
func (e *spinnerEntry) MinSize() fyne.Size {
	size := e.NumericalEntry.MinSize()
	th := e.spinner.Theme()
	iconSpace := th.Size(theme.SizeNameInlineIcon)
	padding := th.Size(theme.SizeNameInnerPadding)
	borderSize := th.Size(theme.SizeNameInputBorder)
	minText := fmt.Sprintf("%d", int(e.spinner.min))
	maxText := fmt.Sprintf("%d", int(e.spinner.max))
	if e.spinner.decimalPlaces != 0 {
		format := fmt.Sprintf("%%.%df", e.spinner.decimalPlaces)
		minText = fmt.Sprintf(format, e.spinner.min)
		maxText = fmt.Sprintf(format, e.spinner.max)
	}
	tSize := e.minSizeForTextSize(e.Text)
	minSize := e.minSizeForTextSize(minText)
	maxSize := e.minSizeForTextSize(maxText)
	wSz := tSize.Width
	if minSize.Width > wSz {
		wSz = minSize.Width
	}
	if maxSize.Width > wSz {
		wSz = maxSize.Width
	}
	return fyne.NewSize(wSz+iconSpace+padding*2+borderSize, size.Height)
}

// TypedShortcut handles the entry's shortcut keys.
func (e *spinnerEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if !e.spinner.Disabled() {
		switch shortcut.(type) {
		case *desktop.CustomShortcut:
			e.shortcut.TypedShortcut(shortcut)
		default:
			e.NumericalEntry.TypedShortcut(shortcut)
		}
	}
}

// TypedKey receives key input events when the spinner's entry widget has focus.
// Increments/decrements the spinner's value when  the up or down key is pressed.
//
// Implements: fyne.Focusable
func (e *spinnerEntry) TypedKey(key *fyne.KeyEvent) {
	if !e.Disabled() {
		switch key.Name {
		case fyne.KeyUp:
			e.spinner.upButtonClicked()
		case fyne.KeyDown:
			e.spinner.downButtonClicked()
		default:
			e.NumericalEntry.TypedKey(key)
		}
	}
}

func (e *spinnerEntry) minSizeForTextSize(text string) fyne.Size {
	t := canvas.NewText(text, color.Black)
	textSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(t.Text,
		t.TextSize, t.TextStyle, t.FontSource)
	return textSize
}

// registerShortcuts registers the shortcuts for the spinnerEntry widget.
func (e *spinnerEntry) registerShortcuts() {
	keyDown := &desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: fyne.KeyModifierControl}
	keyUp := &desktop.CustomShortcut{KeyName: fyne.KeyUp, Modifier: fyne.KeyModifierControl}
	e.shortcut.AddShortcut(keyDown, func(shortcut fyne.Shortcut) {
		e.spinner.SetValue(e.spinner.min)
	})
	e.shortcut.AddShortcut(keyUp, func(shortcut fyne.Shortcut) {
		e.spinner.SetValue(e.spinner.max)
	})
}

// validate tests the text in the entry widget of the spinner to ensure that
// it represents a number and is between the spinner's min and max values.
func (e *spinnerEntry) validate(text string) error {
	if !e.isNumber(text) {
		return errors.New("value is not a number")
	}
	if e.AllowFloat {
		v, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		if v < e.spinner.min || v > e.spinner.max {
			return errors.New("value is not between min and max")
		}
		return nil

	} else {
		v, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		if v < int(e.spinner.min) || v > int(e.spinner.max) {
			return errors.New("value is not between min and max")
		}
		return nil
	}
}

var _ fyne.Disableable = (*Spinner)(nil)

// Spinner is the Spinner widget.
type Spinner struct {
	widget.DisableableWidget

	value         float64
	min           float64
	max           float64
	step          float64
	decimalPlaces uint
	initialized   bool

	entry      *spinnerEntry
	upButton   *spinnerButton
	downButton *spinnerButton
}

// NewSpinner creates a new Spinner object. The initrial value is set to the
// value of the min argument.
//
// Params:
//
//	min is the minimum value that the spinner can be set to. It may be < 0.
//	max is the maximum value that the spinner can be set to. It must be > min.
//	step is the amount that the spinner value increases or decreases by. It
//
// must be > 0 and <= max-min.
//
//	decPlaces is the number of decimal places to display the value in. This
//
// value must be <= 10. If this
//
// value is 0, then the spinner displays integer values.
//
// Returns a Spinner object if the arguments passed are valid. Otherwise, nil
// is returned.
func NewSpinner(min, max, step float64, decPlaces uint) *Spinner {
	if min >= max {
		fyne.LogError("spinner max must be > min", nil)
		return nil
	}
	if step <= 0 {
		fyne.LogError("step must be > 0", nil)
		return nil
	}
	if step > max-min {
		fyne.LogError("step must be <= max-min", nil)
		return nil
	}
	if decPlaces > 10 {
		fyne.LogError("decimal places must be <= 10", nil)
		return nil
	}

	s := &Spinner{min: min, max: max, step: step,
		decimalPlaces: decPlaces, initialized: true}
	s.ExtendBaseWidget(s)

	s.entry = newSpinnerEntry(s)
	s.upButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp),
		s.upButtonClicked)
	s.downButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown),
		s.downButtonClicked)
	s.Enable()

	if s.min < 0 {
		s.entry.AllowNegative = true
	}
	if s.decimalPlaces != 0 {
		s.entry.AllowFloat = true
	}
	s.SetValue(s.min)
	return s
}

// NewSpinnerUnitialized returns an uninitialized Spinner object.
//
// An uninitialized Spinner object is useful when you need to create a Spinner
// but the initial settings are unknown.
// Calling Enable on an unitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize the spinner values before enabling
// the spinner widget.
func NewSpinnerUninitialized(decPlaces uint) *Spinner {
	s := &Spinner{decimalPlaces: decPlaces, initialized: false}
	s.ExtendBaseWidget(s)
	s.entry = newSpinnerEntry(s)
	s.upButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp),
		s.upButtonClicked)
	s.downButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown),
		s.downButtonClicked)
	if s.decimalPlaces != 0 {
		s.entry.AllowFloat = true
	}
	s.Disable()
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	r := &spinnerRenderer{spinner: s}
	r.objects = []fyne.CanvasObject{s.entry, s.upButton, s.downButton}
	return r
}

// Disable disables the spinner and all of its components.
func (s *Spinner) Disable() {
	s.DisableableWidget.Disable()
	s.entry.Disable()
	s.upButton.Disable()
	s.downButton.Disable()
	s.Refresh()
}

// Enable enables the spinner and all of its components.
func (s *Spinner) Enable() {
	if !s.initialized {
		fyne.LogError("Trying to enable uninitialized spinner", nil)
		return
	}
	s.DisableableWidget.Enable()
	s.entry.Enable()
	s.upButton.Enable()
	s.downButton.Enable()
	s.Refresh()
}

// GetValue returns the value of the spinner.
func (s *Spinner) GetValue() float64 {
	return s.value
}

// MinSize returns the minimum size of a Spinner widget. This minimum size is
// calculated based on the maximum width that the spinner's value would require
// based on it's format.
func (s *Spinner) MinSize() fyne.Size {
	w := s.entry.MinSize().Width + s.upButton.MinSize().Width
	h := s.entry.MinSize().Height

	return fyne.NewSize(w, h)
}

// SetMinMaxStep sets the spinner's min, max and step values. It also
// sets the value to min, and enables the spinner.
func (s *Spinner) SetMinMaxStep(min, max, step float64) {
	s.min = min
	s.max = max
	s.step = step
	s.initialized = true
	s.entry.AllowNegative = s.min < 0
	s.Enable()
	s.SetValue(s.min)
	s.Refresh()
}

// SetValue sets the spinner value. If the value is < min, then the
// value is set to min. If the value is > max, then the value is set
// to max. The spinner's buttons are enabled or disabled as appropriate.
func (s *Spinner) SetValue(value float64) {
	if s.Disabled() {
		return
	}
	if !s.initialized {
		fyne.LogError("Trying to set value of uninitialized spinner", nil)
		return
	}
	s.value = value
	if value <= s.min {
		s.value = s.min
		s.downButton.Disable()
	} else {
		s.downButton.Enable()
	}

	if value >= s.max {
		s.value = s.max
		s.upButton.Disable()
	} else {
		s.upButton.Enable()
	}
	s.Refresh()
}

// downButtonClicked handles Tap events for the spinner's down button.
func (s *Spinner) downButtonClicked() {
	if s.Disabled() {
		return
	}
	s.SetValue(s.value - s.step)
}

// upButtonClicked handles Tap events for the spinner's up button.
func (s *Spinner) upButtonClicked() {
	if s.Disabled() {
		return
	}
	s.SetValue(s.value + s.step)
}

// spinnerRenderer is the renderer for the Spinner widget.
type spinnerRenderer struct {
	spinner *Spinner

	objects []fyne.CanvasObject
}

// Destroy destroys any objects that should be destroyed when the renderer is destroyed.
func (r *spinnerRenderer) Destroy() {}

// Layout positions and sizes all of the objects that make up the Spinner widget.
func (r *spinnerRenderer) Layout(size fyne.Size) {
	th := r.spinner.Theme()
	borderSize := th.Size(theme.SizeNameInputBorder)
	padding := th.Size(theme.SizeNameInnerPadding)
	buttonSize := r.spinner.upButton.MinSize()

	xPos := float32(0)
	yPos := float32(0)
	r.spinner.entry.Resize(r.spinner.entry.MinSize())
	r.spinner.entry.Move(fyne.NewPos(xPos, yPos))

	xPos += r.spinner.entry.Size().Width + padding/4
	yPos += borderSize
	r.spinner.upButton.Resize(buttonSize)
	r.spinner.upButton.Move(fyne.NewPos(xPos, yPos))

	yPos = r.spinner.entry.Size().Height - r.spinner.downButton.Size().Height - borderSize
	r.spinner.downButton.Resize(buttonSize)
	r.spinner.downButton.Move(fyne.NewPos(xPos, yPos))

	r.spinner.Refresh()
}

// MinSize returns the minimum size that the Spinner widget can be rendered to.
func (r *spinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

// Objects returns the objects that are rendered by the Spinner renderer.
func (r *spinnerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh redisplays the Spinner widget.
func (r *spinnerRenderer) Refresh() {
	if r.spinner.decimalPlaces == 0 {
		r.spinner.entry.SetText(fmt.Sprintf("%d", int(r.spinner.value)))
	} else {
		format := fmt.Sprintf("%%.%df", r.spinner.decimalPlaces)
		r.spinner.entry.SetText(fmt.Sprintf(format, r.spinner.value))
	}
	r.spinner.entry.Refresh()
	r.spinner.upButton.Refresh()
	r.spinner.downButton.Refresh()
}
