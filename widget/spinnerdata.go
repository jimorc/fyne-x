package widget

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
)

// maxDecimals is the maximum number of decimal places that can be displayed.
var maxDecimals uint = 6

// Spinnable is an interface for specifying if a widget is spinnable (i.e. is a spinner).
type Spinnable interface {
	GetOnChanged() func(float64)
}

// SpinnerData contains the data used by various spinner widget types.
type SpinnerData struct {
	s           Spinnable
	value       float64
	min         float64
	max         float64
	step        float64
	format      string
	initialized bool
}

// NewSpinnerData creates and initializes a new spinnerData object.
//
// Params:
//
//	spinnable is the spinner object that this data is associated with.
//		min is the minimum spinner value. It may be < 0.
//		max is the maximum spinner value. It must be > min.
//		step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//	 	decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
func newSpinnerData(spinnable Spinnable, min, max, step float64, decPlaces uint) *SpinnerData {
	d := newSpinnerDataUninitialized(spinnable, decPlaces)
	d.min = min
	d.max = max
	d.step = step
	d.initialized = d.Validate() == nil

	if d.initialized {
		d.value = min
	}
	return d
}

// NewSpinnerDataUninitialized creates an uninitialized spinnerData object.
//
// Params:
//
//	spinnable is the spinner object that this data is associated with.
//	 decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
func newSpinnerDataUninitialized(spinnable Spinnable, decPlaces uint) *SpinnerData {
	d := &SpinnerData{
		s:           spinnable,
		initialized: false,
	}
	if decPlaces > maxDecimals {
		fyne.LogError(fmt.Sprintf("spinner decPlaces: %d too large. Set to %d", decPlaces, maxDecimals), nil)
		decPlaces = maxDecimals
	}
	if decPlaces == 0 {
		d.format = "%d"
	} else {
		d.format = fmt.Sprintf("%%.%df", decPlaces)
	}
	return d
}

// Validate validates the spinnerData settings.
func (d *SpinnerData) Validate() error {
	if d.min == 0. && d.max == 0. && d.step == 0. {
		return errors.New("spinner not initialized")
	}
	if d.min >= d.max {
		return errors.New("spinner max value must be greater than min value")
	}
	if d.step <= 0 {
		return errors.New("spinner step must be greater than 0")
	}
	if d.step > d.max-d.min {
		return errors.New("spinner step must be less than or equal to max - min")
	}
	return nil
}
