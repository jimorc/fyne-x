package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// maxDecimals is the maximum number of decimal places that can be displayed.
var maxDecimals uint = 6

// Spinnable is an interface for specifying if a widget is spinnable (i.e. is a spinner).
type Spinnable interface {
	getOnChanged() func(float64)
}

// SpinnerData contains the data used by various spinner widget types.
type spinnerData struct {
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
func newSpinnerData(spinnable Spinnable, min, max, step float64, decPlaces uint) *spinnerData {
	d := newSpinnerDataUninitialized(spinnable, decPlaces)
	d.min = min
	d.max = max
	d.step = step
	d.initialized = true
	if min >= max {
		fyne.LogError("Spinner max value must be greater than min value", nil)
		d.initialized = false
	}
	if step < 1 {
		fyne.LogError("Spinner step must be greater than 0", nil)
		d.initialized = false
	}
	if step > max-min {
		fyne.LogError("Spinner step must be less than or equal to max - min", nil)
		d.initialized = false
	}
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
func newSpinnerDataUninitialized(spinnable Spinnable, decPlaces uint) *spinnerData {
	d := &spinnerData{
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
