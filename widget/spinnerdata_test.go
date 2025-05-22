package widget

import (
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

var val float64 = 0.
var dVal int = 0

type spinner struct {
	disabled  bool
	onChanged func(float64)
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
	return s.onChanged
}

func TestSpinnerData_NewSpinnerData(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 10, 2)

	assert.Equal(t, 1., d.min)
	assert.Equal(t, 10., d.max)
	assert.Equal(t, 2., d.step)
	assert.Equal(t, 1., d.value)
}

func TestSpinnerData_InvalidArgs(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 11, 10, 2)
	assert.False(t, d.initialized)

	d = NewSpinnerData(s, 1, 10, 0)
	assert.False(t, d.initialized)

	d = NewSpinnerData(s, 1, 2, 2)
	assert.False(t, d.initialized)
}

func TestSpinnerDataUninitialized(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerDataUninitialized(s)
	assert.False(t, d.initialized)
}

func TestSpinnerData_Validate(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 2, 1)
	err := d.Validate()
	assert.Nil(t, err)

	d = NewSpinnerDataUninitialized(s)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner not initialized", err.Error())

	d = NewSpinnerData(s, 2, 2, 1)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner max value must be greater than min value", err.Error())

	d = NewSpinnerData(s, 1, 2, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be greater than 0", err.Error())

	d = NewSpinnerData(s, 1, 2, 3)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be less than or equal to max - min", err.Error())
}

func TestSpinnerData_SetValue(t *testing.T) {
	s := &spinner{}
	s.onChanged = func(v float64) {
		val = v
	}
	d := NewSpinnerData(s, 1, 4, 1)
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
	d := NewSpinnerData(s, 1, 4, 1)
	d.SetValue(4.)
	d.Decrement()
	assert.Equal(t, 3., d.Value())
	d.Decrement()
	assert.Equal(t, 2., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())

	d = NewSpinnerData(s, 1, 4, 2)
	d.SetValue(4.)
	d.Decrement()
	assert.Equal(t, 2., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())

	d.SetValue(4.)
	d.s.Disable()
	d.Decrement()
	assert.Equal(t, 4., d.Value())
}

func TestSpinnerData_Increment(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 4, 1)
	d.Increment()
	assert.Equal(t, 2., d.Value())
	d.Increment()
	assert.Equal(t, 3., d.Value())
	d.Increment()
	assert.Equal(t, 4., d.Value())
	d.Increment()
	assert.Equal(t, 4., d.Value())

	d.SetValue(1.)
	d.s.Disable()
	d.Increment()
	assert.Equal(t, 1., d.Value())
}

func TestSpinnerData_ValueChanged(t *testing.T) {
	dVal = 0
	val = 0.
	s := &spinner{}
	d := NewSpinnerData(s, 1, 4, 1)
	assert.Equal(t, 0., val)
	d.onChanged = func(float64) {
		dVal++
	}
	d.SetValue(2)
	assert.Equal(t, 0., val)
	assert.Equal(t, 1, dVal)

	s.onChanged = func(v float64) {
		val = v
	}
	d.Increment()
	assert.Equal(t, 3., val)
	assert.Equal(t, 2, dVal)
}

func TestSpinnerData_SetMinMaxStep(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerDataUninitialized(s)
	assert.False(t, d.initialized)
	d.SetMinMaxStep(4, 1, 1)
	assert.False(t, d.initialized)
	assert.Equal(t, 4., d.min)
	assert.Equal(t, 1., d.max)
	assert.Equal(t, 1., d.step)
	assert.Equal(t, 0., d.Value())

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
	d := NewSpinnerDataWithData(s, 1, 12, 1, data)
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
