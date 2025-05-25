package widget

import (
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

var val float64 = 0.
var dVal int = 0

var _ Spinnable = (*spinner)(nil)

type spinner struct {
	disabled  bool
	OnChanged func(float64)
}

func (s *spinner) Disable() {
	s.disabled = true
}

func (s *spinner) Disabled() bool {
	return s.disabled
}

func (s *spinner) Enable() {
	s.disabled = false
}

func (s *spinner) GetOnChanged() func(float64) {
	return s.OnChanged
}

func (s *spinner) Refresh() {}

func TestSpinnerData_NewSpinnerData(t *testing.T) {
	s := &Spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 10, 2)
	assert.True(t, d.initialized)
	assert.Equal(t, 1., d.min)
	assert.Equal(t, 10., d.max)
	assert.Equal(t, 2., d.step)
	assert.Equal(t, 1., d.value)
}

func TestSpinnerData_InvalidArgs(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 11, 10, 2)
	assert.False(t, d.initialized)

	d = newSpinnerData(b, 1, 10, 0)
	assert.False(t, d.initialized)

	d = newSpinnerData(b, 1, 2, 2)
	assert.False(t, d.initialized)
}

func TestSpinnerDataUninitialized(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerDataUninitialized(b)
	assert.False(t, d.initialized)
}

func TestSpinnerData_Validate(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 2, 1)
	err := d.Validate()
	assert.Nil(t, err)

	d = newSpinnerDataUninitialized(b)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner not initialized", err.Error())

	d = newSpinnerData(b, 2, 2, 1)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner max value must be greater than min value", err.Error())

	d = newSpinnerData(b, 1, 2, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be greater than 0", err.Error())

	d = newSpinnerData(b, 1, 2, 3)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be less than or equal to max - min", err.Error())
}

func TestSpinnerData_SetValue(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 4, 1)
	d.onChanged = func(v float64) {
		val = v
	}
	assert.Equal(t, 1., d.Value())
	d.SetValue(2)
	assert.Equal(t, 2., val)
	assert.Equal(t, 2., d.Value())
	d.SetValue(5)
	assert.Equal(t, 4., val)
	assert.Equal(t, 4., d.Value())
	d.SetValue(0)
	assert.Equal(t, 1., val)
	assert.Equal(t, 1., d.Value())
}

func TestSpinnerData_Decrement(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 4, 1)
	d.SetValue(4.)
	d.Decrement()
	assert.Equal(t, 3., d.Value())
	d.Decrement()
	assert.Equal(t, 2., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())

	d = newSpinnerData(b, 1, 4, 2)
	d.SetValue(4.)
	d.Decrement()
	assert.Equal(t, 2., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())

	d.SetValue(4.)
	d.base.spinner.Disable()
	d.Decrement()
	assert.Equal(t, 4., d.Value())
}

func TestSpinnerData_Increment(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 4, 1)
	d.Increment()
	assert.Equal(t, 2., d.Value())
	d.Increment()
	assert.Equal(t, 3., d.Value())
	d.Increment()
	assert.Equal(t, 4., d.Value())
	d.Increment()
	assert.Equal(t, 4., d.Value())

	d.SetValue(1.)
	d.base.spinner.Disable()
	d.Increment()
	assert.Equal(t, 1., d.Value())
}

func TestSpinnerData_ValueChanged(t *testing.T) {
	dVal = 0
	val = 0.
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerData(b, 1, 4, 1)
	assert.Equal(t, 0., val)
	d.onChanged = func(v float64) {
		dVal++
		val = v
	}
	d.SetValue(1)
	assert.Equal(t, 1., d.value)
	assert.Equal(t, 0., val)
	assert.Equal(t, 0, dVal)
	d.SetValue(2)
	assert.Equal(t, 2., val)
	assert.Equal(t, 1, dVal)
	d.SetValue(2)
	assert.Equal(t, 2., val)
	assert.Equal(t, 1, dVal)
	d.Increment()
	assert.Equal(t, 3., val)
	assert.Equal(t, 2, dVal)
}

func TestSpinnerData_SetMinMaxStep(t *testing.T) {
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerDataUninitialized(b)
	assert.False(t, d.initialized)
	d.SetMinMaxStep(4, 1, 1)
	assert.False(t, d.initialized)
	assert.Equal(t, 4., d.min)
	assert.Equal(t, 1., d.max)
	assert.Equal(t, 1., d.step)
	assert.Equal(t, 1., d.Value())

	d.SetMinMaxStep(1, 4, 1)
	assert.True(t, d.initialized)
	assert.Equal(t, 1., d.min)
	assert.Equal(t, 4., d.max)
	assert.Equal(t, 1., d.step)
	assert.Equal(t, 1., d.Value())
}

func TestSpinnerData_NewSpinnerDataWithData(t *testing.T) {
	var data binding.Float = binding.NewFloat()
	s := &spinner{}
	b := &SpinnerBase{spinner: s}
	d := newSpinnerDataWithData(b, 1, 12, 1, data)
	data.Set(10.)
	waitForBinding()
	assert.Equal(t, 10., d.Value())

	d.SetValue(4.)
	waitForBinding()
	assert.Equal(t, 4., d.Value())

	d.Unbind()
	data.Set(6.)
	waitForBinding()
	assert.Equal(t, 4., d.Value())
}
