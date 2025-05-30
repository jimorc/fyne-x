package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// maxDecimals is the maximum number of decimal places that can be displayed.
var maxDecimals uint = 6

// Spinnable is an interface for specifying if a widget is spinnable (i.e. is a spinner).
type Spinnable interface {
	fyne.Disableable
	// GetOnChanged retrieves the function to execute when the SpinnerData object changes
	// its value.
	GetOnChanged() func(float64)
	// Refresh redisplays the Spinnable.
	Refresh()
}

// SpinnerBase contains functionality that is common to all spinner widgets. It has a minimum,
// maximum, step and current values along with spinnerButtons
// to increment and decrement the spinner value.
type SpinnerBase struct {
	spinner    Spinnable
	data       *spinnerData
	upButton   *spinnerButton
	downButton *spinnerButton

	format string
	mPr    *message.Printer
}

// NewSpinnerBase creates and initializes a new SpinnerBase object.
//
// Params:
//
//	s is the parent spinner object.
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
//
//	onChanged is the callback function that is called whenever the spinner value changes.
func NewSpinnerBase(s Spinnable, min, max, step float64, decPlaces uint) *SpinnerBase {
	base := &SpinnerBase{spinner: s}

	base.upButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), base.Increment)
	base.downButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropDown), base.Decrement)
	base.data = newSpinnerData(base, min, max, step)
	if base.data.initialized {
		s.Enable()
		base.upButton.Enable()
		base.downButton.Disable()
	}
	base.setFormat(decPlaces)
	return base
}

// NewSpinnerBaseUninitialized returns a new uninitialized SpinnerBase object.
//
// An uninitialized Spinner widget is useful when you need to create a Spinner
// but the initial settings are unknown.
// Calling Enable on an uninitialized spinner will not enable the spinner; you
// must first call SetMinMaxStep to initialize spinner values before enabling
// the spinner widget.
//
// Params:
//
//	s is the parent spinner object.
//	decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
func NewSpinnerBaseUninitialized(s Spinnable, decPlaces uint) *SpinnerBase {
	base := &SpinnerBase{spinner: s}
	base.setFormat(decPlaces)
	base.data = newSpinnerData(base, 0, 0, 0)
	base.upButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), base.Increment)
	base.downButton = newSpinnerButton(theme.Icon(theme.IconNameArrowDropDown), base.Decrement)
	base.upButton.enableDisable(base.spinner.Disabled(), base.AtMax())
	base.downButton.enableDisable(base.spinner.Disabled(), base.AtMin())
	s.Disable()
	base.upButton.Disable()
	base.downButton.Disable()
	return base
}

// NewSpinnerBaseWithData returns a new Spinner widget connected to the specified data source.
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
func NewSpinnerBaseWithData(s Spinnable, min, max, step float64,
	decPlaces uint, data binding.Float) *SpinnerBase {
	base := NewSpinnerBase(s, min, max, step, decPlaces)
	base.Bind(data)
	return base
}

// AtMax returns true if the spinner data value is at its max value.
func (s *SpinnerBase) AtMax() bool {
	return s.data.atMax()
}

// AtMin returns true if the spinner data value is at its min value.
func (s *SpinnerBase) AtMin() bool {
	return s.data.atMin()
}

// Bind connects the specified data source to the Spinnable object.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (s *SpinnerBase) Bind(data binding.Float) {
	s.data.bind(data)
}

// DownButton returns a pointer to the SpinnerBase downButton.
func (s *SpinnerBase) DownButton() *spinnerButton {
	return s.downButton
}

// EnableDisableButtons enables or disables the up and down buttons based on whether
// the parent spinner is disabled and on whether the data value is equal to max or min.
func (s *SpinnerBase) EnableDisableButtons(spinnerDisabled bool) {
	s.upButton.enableDisable(spinnerDisabled, s.data.atMax())
	s.downButton.enableDisable(spinnerDisabled, s.data.atMin())
}

// GetOnChanged retrieves the onChanged function for the spinner.
//
// Implements the Spinnable interface.
func (s *SpinnerBase) GetOnChanged() func(float64) {
	if s.data != nil {
		return func(float64) {
			s.downButton.enableDisable(false, s.data.atMin())
			s.upButton.enableDisable(false, s.data.atMax())
			spinnerOnChanged := s.spinner.GetOnChanged()
			if spinnerOnChanged != nil {
				spinnerOnChanged(s.data.Value())
			}
		}

	}
	return func(float64) {}
}

// Initialized returns true if the SpinnerBase's SpinnerData object has been initialized.
func (s *SpinnerBase) Initialized() bool {
	if s.data == nil {
		return false
	}
	return s.data.initialized
}

// MaxText returns the max value as a formatted string.
// This method is useful for determining the minimum required widget size.
func (s *SpinnerBase) MaxText() string {
	return s.formatAsText(s.data.max, s.format)
}

// MaxValue returns the spinnerData max value.
func (s *SpinnerBase) MaxValue() float64 {
	return s.data.max
}

// MinText returns the min value as a formatted string.
// This method is useful for determining the minimum required widget size
func (s *SpinnerBase) MinText() string {
	return s.formatAsText(s.data.min, s.format)
}

// MinValue returns the spinnerData min value.
func (s *SpinnerBase) MinValue() float64 {
	return s.data.min
}

func (s *SpinnerBase) SetMinMaxStep(min, max, step float64) {
	if s.data == nil {
		s.data = newSpinnerData(s, min, max, step)
		return
	}
	s.data.setMinMaxStep(min, max, step)
	s.upButton.enableDisable(s.spinner.Disabled(), s.data.atMax())
	s.downButton.enableDisable(s.spinner.Disabled(), s.data.atMin())
}

func (s *SpinnerBase) SetValue(value float64) {
	s.data.setValue(value)
	if s.spinner.Disabled() {
		return
	}
	s.upButton.enableDisable(false, s.data.atMax())
	s.downButton.enableDisable(false, s.data.atMin())
}

func (s *SpinnerBase) StepValue() float64 {
	return s.data.step
}

// Unbind removes the binding from the spinner data.
func (s *SpinnerBase) Unbind() {
	s.data.unbind()
}

// UpButton returns a pointer to the SpinnerBase upButton.
func (s *SpinnerBase) UpButton() *spinnerButton {
	return s.upButton
}

// Validate validates the spinnerData values.
func (s *SpinnerBase) Validate() error {
	return s.data.validate()
}

// Value returns the spinnerData value.
func (s *SpinnerBase) Value() float64 {
	return s.data.Value()
}

// ValueText retrieves the spinner value as formatted text.
func (s *SpinnerBase) ValueText() string {
	return s.formatAsText(s.data.value, s.format)
}

// Decrement decrements the data's value by step amount, or to min if that is larger.
func (s *SpinnerBase) Decrement() {
	s.data.decrement()
	if s.spinner.Disabled() {
		return
	}
	s.downButton.enableDisable(false, s.data.atMin())
	s.upButton.Enable()
	s.spinner.Refresh()
}

// Increment icrements the data's value by step amount, or to max if that is less.
func (s *SpinnerBase) Increment() {
	s.data.increment()
	if s.spinner.Disabled() {
		return
	}
	s.upButton.enableDisable(false, s.data.atMax())
	s.downButton.Enable()
	s.spinner.Refresh()
}

// formatAsText formats the value according to the specified format.
//
// Params:
//
//	value is the value to format.
//	format is the format to use. This format should be either "%d", or "%.nf"
//	where n is either an empty string or an integer.
func (s *SpinnerBase) formatAsText(value float64, format string) string {
	if format == "%d" {
		return s.mPr.Sprintf(format, int(value))
	} else {
		return s.mPr.Sprintf(format, value)
	}
}

// setFormat determines the format to display the value in.
//
// Params:
//
//	decPlaces is the number of decimal places to display the value with.
//	If decPlaces == 0, the value is displayed as an integer.
//	If decPlaces > maxDecimals, it is set to maxDecimals.
func (s *SpinnerBase) setFormat(decPlaces uint) {
	if decPlaces > maxDecimals {
		fyne.LogError(fmt.Sprintf("spinner decPlaces: %d too large. Set to %d", decPlaces, maxDecimals), nil)
		decPlaces = maxDecimals
	}
	if decPlaces == 0 {
		s.format = "%d"
	} else {
		s.format = fmt.Sprintf("%%.%df", decPlaces)
	}
	locale := lang.SystemLocale().String()
	lang, err := language.Parse(locale)
	if err != nil {
		fyne.LogError(fmt.Sprintf("'%s' parse error: ", locale), err)
		lang = language.English // Fallback to English
		locale = "en-US"
	}
	s.mPr = message.NewPrinter(lang)
}
