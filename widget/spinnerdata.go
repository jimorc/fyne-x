package widget

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

// Spinnable is an interface for specifying if a widget is spinnable (i.e. is a spinner).
type Spinnable interface {
	fyne.Disableable
	// GetOnChanged retrieves the function to execute when the SpinnerData object changes
	// its value.
	GetOnChanged() func(float64)
}

// SpinnerData contains the data used by various spinner widget types.
type SpinnerData struct {
	s     Spinnable
	value float64
	min   float64
	max   float64
	step  float64

	binder basicBinder

	initialized bool
	onChanged   func(float64)
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
func NewSpinnerData(spinnable Spinnable, min, max, step float64) *SpinnerData {
	d := NewSpinnerDataUninitialized(spinnable)
	d.SetMinMaxStep(min, max, step)
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
func NewSpinnerDataUninitialized(spinnable Spinnable) *SpinnerData {
	d := &SpinnerData{
		s:           spinnable,
		initialized: false,
	}
	return d
}

func NewSpinnerDataWithData(s Spinnable, min, max, step float64,
	data binding.Float) *SpinnerData {
	d := NewSpinnerData(s, min, max, step)

	d.Bind(data)
	return d

}

// Bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (d *SpinnerData) Bind(data binding.Float) {
	d.binder.SetCallback(d.updateFromData)
	d.binder.Bind(data)
	d.onChanged = func(_ float64) {
		d.binder.CallWithData(d.writeData)
	}
}

// Decrement decrements the SpinnerData object's value by its step size.
func (d *SpinnerData) Decrement() {
	d.SetValue(d.value - d.step)
}

// Increment increments the SpinnerData object's value by its step size.
func (d *SpinnerData) Increment() {
	d.SetValue(d.value + d.step)
}

// SetMinMaxStep sets the SpinnerData's minimum, maximum, and step values. The SpinnerData
// object's value is set to minimum if the data passes validation.
//
// Params:
//
//	min is the minimum spinner value. It may be < 0.
//	max is the maximum spinner value. It must be > min.
//	step is the amount that the spinner increases or decreases by. It must be > 0 and less than or equal to max - min.
//
// If the previously set value is less than min, then the value is set to min.
// If the previously set value is greater than max, then the value is set to max.
func (d *SpinnerData) SetMinMaxStep(min, max, step float64) {
	d.min = min
	d.max = max
	d.step = step
	d.initialized = d.Validate() == nil
	if d.initialized {
		//		d.s.Enable()
		d.SetValue(min)
	}
}

// SetValue sets the value in the SpinnerData object.
// If the value is less than object's min value, the value is set to min.
// If the value is greater than object's max value, the value is set to max.
func (d *SpinnerData) SetValue(value float64) {
	if d.s.Disabled() {
		return
	}
	d.value = value
	if d.value >= d.max {
		d.value = d.max
	}
	if d.value <= d.min {
		d.value = d.min
	}
	d.valueChanged()
}

// Unbind disconnects any configured data source from this spinnerData.
// The current value will remain at the last value of the data source.
func (d *SpinnerData) Unbind() {
	d.binder.Unbind()
	d.onChanged = nil
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

// Value retrieves the value set in the SpinnerData object.
func (d *SpinnerData) Value() float64 {
	return d.value
}

// updateFromData updates the spinner to the value set in the bound data.
func (d *SpinnerData) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.Float)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	d.SetValue(val)
}

// valueChanged executes any onChanged functions in the SpinnerData and Spinnable objects.
// This method is executed every time the value changes in the SpinnerData object.
func (d *SpinnerData) valueChanged() {
	if d.onChanged != nil {
		d.onChanged(d.value)
	}
	spinnerOnChanged := d.s.GetOnChanged()
	if spinnerOnChanged != nil {
		spinnerOnChanged(d.value)
	}
}

// writeData updates the bound data item as the result of changes in the spinnerData value.
func (d *SpinnerData) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	floatTarget, ok := data.(binding.Float)
	if !ok {
		return
	}
	currentValue, err := floatTarget.Get()
	if err != nil {
		return
	}
	if currentValue != d.Value() {
		err := floatTarget.Set(d.Value())
		if err != nil {
			fyne.LogError(fmt.Sprintf("Failed to set binding value to %f", d.Value()), err)
		}
	}
}
