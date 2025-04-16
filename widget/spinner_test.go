package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner(1, 10, 2, 0)
	assert.Equal(t, float64(1), s.min)
	assert.Equal(t, float64(10), s.max)
	assert.Equal(t, float64(2), s.step)
	assert.Equal(t, uint(0), s.decimalPlaces)
	assert.False(t, s.entry.AllowFloat)
	assert.False(t, s.entry.AllowNegative)

	s2 := NewSpinner(-5, 5, 1, 2)
	assert.Equal(t, float64(-5), s2.min)
	assert.Equal(t, float64(5), s2.max)
	assert.Equal(t, float64(1), s2.step)
	assert.Equal(t, uint(2), s2.decimalPlaces)
	assert.True(t, s2.entry.AllowNegative)
	assert.True(t, s2.entry.AllowFloat)
}

func TestNewSpinner_Invalid(t *testing.T) {
	s := NewSpinner(10, 1, 2, 0)
	assert.Nil(t, s)

	s2 := NewSpinner(1, 1, 2, 0)
	assert.Nil(t, s2)

	s3 := NewSpinner(1, 10, 11, 0)
	assert.Nil(t, s3)

	s4 := NewSpinner(1, 10, 2, 11)
	assert.Nil(t, s4)
}
