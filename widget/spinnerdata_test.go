package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type spinner struct {
	onChanged func(float64)
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
	d := newSpinnerDataUninitialized(s, 0)
	assert.False(t, d.initialized)
	assert.Equal(t, "%d", d.format)

	d = newSpinnerDataUninitialized(s, 1)
	assert.Equal(t, "%.1f", d.format)

	d = newSpinnerDataUninitialized(s, 5)
	assert.Equal(t, "%.5f", d.format)

	d = newSpinnerDataUninitialized(s, 10)
	assert.Equal(t, "%.6f", d.format)
}

func TestSpinnerData_Validate(t *testing.T) {
	s := &spinner{}
	d := NewSpinnerData(s, 1, 2, 1, 1)
	err := d.Validate()
	assert.Nil(t, err)

	d = newSpinnerDataUninitialized(s, 0)
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
