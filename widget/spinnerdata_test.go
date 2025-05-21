package widget

import (
	"testing"

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
	d := NewSpinnerData(s, 1, 10, 2, 0)

	assert.Equal(t, 1., d.min)
	assert.Equal(t, 10., d.max)
	assert.Equal(t, 2., d.step)
	assert.Equal(t, 1., d.value)
	assert.Equal(t, "%d", d.format)
}

func TestSpinnerData_Format(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 10, 2, 0)
	assert.Equal(t, "%d", d.format)

	d = NewSpinnerData(s, 1, 10, 2, 1)
	assert.Equal(t, "%.1f", d.format)

	d = NewSpinnerData(s, 1, 10, 2, 5)
	assert.Equal(t, "%.5f", d.format)

	d = NewSpinnerData(s, 1, 10, 2, 10)
	assert.Equal(t, "%.6f", d.format)
}

func TestSpinnerData_InvalidArgs(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 11, 10, 2, 0)
	assert.False(t, d.initialized)

	d = NewSpinnerData(s, 1, 10, 0, 0)
	assert.False(t, d.initialized)

	d = NewSpinnerData(s, 1, 2, 2, 0)
	assert.False(t, d.initialized)
}

func TestSpinnerDataUninitialized(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerDataUninitialized(s, 0)
	assert.False(t, d.initialized)
	assert.Equal(t, "%d", d.format)

	d = NewSpinnerDataUninitialized(s, 1)
	assert.Equal(t, "%.1f", d.format)

	d = NewSpinnerDataUninitialized(s, 5)
	assert.Equal(t, "%.5f", d.format)

	d = NewSpinnerDataUninitialized(s, 10)
	assert.Equal(t, "%.6f", d.format)
}

func TestSpinnerData_Validate(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 2, 1, 1)
	err := d.Validate()
	assert.Nil(t, err)

	d = NewSpinnerDataUninitialized(s, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner not initialized", err.Error())

	d = NewSpinnerData(s, 2, 2, 1, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner max value must be greater than min value", err.Error())

	d = NewSpinnerData(s, 1, 2, 0, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be greater than 0", err.Error())

	d = NewSpinnerData(s, 1, 2, 3, 0)
	err = d.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "spinner step must be less than or equal to max - min", err.Error())
}

func TestSpinnerData_SetValue(t *testing.T) {
	s := &spinner{}
	s.onChanged = func(v float64) {
		val = v
	}
	d := NewSpinnerData(s, 1, 4, 1, 0)
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
	d := NewSpinnerData(s, 1, 4, 1, 0)
	d.SetValue(4.)
	d.Decrement()
	assert.Equal(t, 3., d.Value())
	d.Decrement()
	assert.Equal(t, 2., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())
	d.Decrement()
	assert.Equal(t, 1., d.Value())

	d = NewSpinnerData(s, 1, 4, 2, 0)
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
	d := NewSpinnerData(s, 1, 4, 1, 0)
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
	d := NewSpinnerData(s, 1, 4, 1, 0)
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
