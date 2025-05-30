package widget

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

// spinnerData contains the data used by various spinner widget types.
type spinnerData struct {
	base  *SpinnerBase
	value float64
	min   float64
	max   float64
	step  float64

	binder basicBinder

	initialized bool
	onChanged   func(float64)
}

// newSpinnerData creates and initializes a new spinnerData object.
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
func newSpinnerData(base *SpinnerBase, min, max, step float64) *spinnerData {
	d := newSpinnerDataUninitialized(base)
	d.setMinMaxStep(min, max, step)
	return d
}

// newSpinnerDataUninitialized creates an uninitialized spinnerData object.
//
// Params:
//
//	spinnable is the spinner object that this data is associated with.
//	 decPlaces is the number of decimal places to display the value in. This value must be
//
// 0 <= decPlaces <= maxDecimals. If this value is greater than maxDecimals, it is set to maxDecimals.
// If decPlaces == 0, then the value is displayed as an integer.
func newSpinnerDataUninitialized(base *SpinnerBase) *spinnerData {
	d := &spinnerData{
		base:        base,
		initialized: false,
	}
	return d
}

// newSpinnerDataWithData creates and initializes a new spinnerData object
// with the data value tied to a binding.Float variable.
func newSpinnerDataWithData(base *SpinnerBase, min, max, step float64,
	data binding.Float) *spinnerData {
	d := newSpinnerData(base, min, max, step)

	d.bind(data)
	return d

}

// atMax returns true if the SpinnerData value is equal to max.
func (d *spinnerData) atMax() bool {
	return d.value >= d.max
}

// atMin returns true if the SpinnerData value is equal to min.
func (d *spinnerData) atMin() bool {
	return d.value <= d.min
}

// bind connects the specified data source to this Spinner widget.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
func (d *spinnerData) bind(data binding.Float) {
	d.binder.SetCallback(d.updateFromData)
	d.binder.Bind(data)
	d.onChanged = func(_ float64) {
		d.binder.CallWithData(d.writeData)
	}
}

// Decrement decrements the SpinnerData object's value by its step size.
func (d *spinnerData) decrement() {
	d.setValue(d.value - d.step)
}

// Increment increments the SpinnerData object's value by its step size.
func (d *spinnerData) increment() {
	d.setValue(d.value + d.step)
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
func (d *spinnerData) setMinMaxStep(min, max, step float64) {
	d.min = min
	d.max = max
	d.step = step
	d.initialized = d.validate() == nil
	if d.initialized {
		d.setValue(d.min)
	}
}

// SetValue sets the value in the SpinnerData object.
// If the value is less than object's min value, the value is set to min.
// If the value is greater than object's max value, the value is set to max.
func (d *spinnerData) setValue(value float64) {
	if d.base.spinner.Disabled() || !d.initialized {
		return
	}
	if d.value == value {
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
func (d *spinnerData) unbind() {
	d.binder.Unbind()
	d.onChanged = nil
}

// Validate validates the spinnerData settings.
func (d *spinnerData) validate() error {
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

// Value retrieves the value set in the SpinnerData object. If outside the min to max
// range, the value will be set to either min or max as appropriate.
func (d *spinnerData) Value() float64 {
	value := d.value
	if !d.initialized || value < d.min {
		d.setValue(d.min)
		value = d.min
	}
	if value > d.max {
		d.setValue(d.max)
		value = d.max
	}
	return value
}

// updateFromData updates the spinner to the value set in the bound data.
func (d *spinnerData) updateFromData(data binding.DataItem) {
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
	d.setValue(val)
}

// valueChanged executes any onChanged functions in the SpinnerData and Spinnable objects.
// This method is executed every time the value changes in the SpinnerData object.
func (d *spinnerData) valueChanged() {
	if d.onChanged != nil {
		d.onChanged(d.value)
	}
	spinnerBaseOnChanged := d.base.GetOnChanged()
	if spinnerBaseOnChanged != nil {
		spinnerBaseOnChanged(d.value)
	}
}

// writeData updates the bound data item as the result of changes in the spinnerData value.
func (d *spinnerData) writeData(data binding.DataItem) {
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
