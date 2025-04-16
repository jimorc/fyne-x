package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
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

// spinnerEntry is the entry widget for the Spinner widget.
type spinnerEntry struct {
	NumericalEntry
}

// newSpinnerEntry creates a spinnerEntry widget.
func newSpinnerEntry() *spinnerEntry {
	e := &spinnerEntry{}
	e.ExtendBaseWidget(e)

	return e
}

// MinSize returns the minimum size of the spinnerEntry.
func (e *spinnerEntry) MinSize() fyne.Size {
	return fyne.NewSize(150, e.NumericalEntry.MinSize().Height)
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

	entry      *spinnerEntry
	upButton   *spinnerButton
	downButton *spinnerButton
}

// NewSpinner creates a new Spinner object.
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

	s := &Spinner{min: min, max: max, step: step, decimalPlaces: decPlaces}
	s.ExtendBaseWidget(s)

	s.value = s.min
	s.entry = newSpinnerEntry()
	s.upButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropUp),
		s.upButtonClicked)
	s.downButton = newSpinnerButton(s, theme.Icon(theme.IconNameArrowDropDown),
		s.downButtonClicked)

	if s.min < 0 {
		s.entry.AllowNegative = true
	}
	if s.decimalPlaces != 0 {
		s.entry.AllowFloat = true
	}
	return s
}

// CreateRenderer is a private method to fyne which links this widget to its
// renderer
func (s *Spinner) CreateRenderer() fyne.WidgetRenderer {
	r := &spinnerRenderer{spinner: s}
	r.objects = []fyne.CanvasObject{s.entry, s.upButton, s.downButton}
	return r
}

// MinSize returns the minimum size of a Spinner widget. This minimum size is
// calculated based on the maximum width that the spinner's value would require
// based on it's format.
func (s *Spinner) MinSize() fyne.Size {
	w := s.entry.MinSize().Width + s.upButton.MinSize().Width
	h := s.entry.MinSize().Height

	return fyne.NewSize(w, h)
}

// downButtonClicked handles Tap events for the spinner's down button.
func (s *Spinner) downButtonClicked() {}

// upButtonClicked handles Tap events for the spinner's up button.
func (s *Spinner) upButtonClicked() {}

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
}
