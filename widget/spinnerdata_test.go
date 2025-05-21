package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type spinner struct {
	onChanged func(float64)
}

func (s *spinner) getOnChanged() func(float64) {
	return s.onChanged
}

func TestSpinnerData_NewSpinnerData(t *testing.T) {
	s := &spinner{}
	d := newSpinnerData(s, 1, 10, 2, 0)

	assert.Equal(t, 1., d.min)
	assert.Equal(t, 10., d.max)
	assert.Equal(t, 2., d.step)
	assert.Equal(t, 1., d.value)
	assert.Equal(t, "%d", d.format)
}

func TestSpinnerData_Format(t *testing.T) {
	s := &spinner{}
	d := newSpinnerData(s, 1, 10, 2, 0)
	assert.Equal(t, "%d", d.format)

	d = newSpinnerData(s, 1, 10, 2, 1)
	assert.Equal(t, "%.1f", d.format)

	d = newSpinnerData(s, 1, 10, 2, 5)
	assert.Equal(t, "%.5f", d.format)

	d = newSpinnerData(s, 1, 10, 2, 10)
	assert.Equal(t, "%.6f", d.format)
}

func TestSpinnerData_InvalidArgs(t *testing.T) {
	s := &spinner{}
	d := newSpinnerData(s, 11, 10, 2, 0)
	assert.False(t, d.initialized)

	d = newSpinnerData(s, 1, 10, 0, 0)
	assert.False(t, d.initialized)

	d = newSpinnerData(s, 1, 2, 2, 0)
	assert.False(t, d.initialized)
}
