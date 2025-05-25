package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func TestNewSpinnerBase(t *testing.T) {
	s := &Spinner{}
	b := NewSpinnerBase(s, 1., 5., 1.5, 0)
	assert.Equal(t, 1., b.data.min)
	assert.Equal(t, 5., b.data.max)
	assert.Equal(t, 1.5, b.data.step)
	assert.Equal(t, 1., b.Value())
	assert.False(t, b.upButton.Disabled())
	assert.True(t, b.downButton.Disabled())
}

func TestNewSpinnerBase_BadArgs(t *testing.T) {
	s := &Spinner{}
	b := NewSpinnerBase(s, 5., 5., 1., 0)
	assert.False(t, b.data.initialized, "spinner should not be initialized when max = min")

	b = NewSpinnerBase(s, 5., 4., 1., 0)
	assert.False(t, b.data.initialized, "spinner should not be initialized when min > max")

	b = NewSpinnerBase(s, 1., 5., 0., 0)
	assert.False(t, b.data.initialized, "spinner should not be initialized when step = 0")

	b = NewSpinnerBase(s, 1., 5., -5., 0)
	assert.False(t, b.data.initialized, "spinner should not be initialized when step < 0")

	b = NewSpinnerBase(s, 1., 5., 5., 0)
	assert.False(t, b.data.initialized, "spinner should not be initialized when step > max - min")

	b = NewSpinnerBase(s, 1., 5., 2., 11)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), b.format)
	assert.True(t, b.data.initialized)
}

func TestNewSpinnerBaseWithData(t *testing.T) {
	data := binding.NewFloat()
	s := &Spinner{}
	b := NewSpinnerBaseWithData(s, 1., 5., 2., 0, data)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)

	b.SetValue(1.52)
	waitForBinding()
	val, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1.52, val)

	err = data.Set(3.1)
	assert.NoError(t, err)
	waitForBinding()
	assert.Equal(t, 3.1, b.Value())
}

func TestSpinnerBase_SetValue(t *testing.T) {
	s := &Spinner{}
	b := NewSpinnerBase(s, 1., 5., 1., 0)
	assert.Equal(t, 1., b.Value())
	assert.True(t, b.DownButton().Disabled())
	assert.False(t, b.UpButton().Disabled())

	b.SetValue(3.)
	assert.Equal(t, 3., b.Value())
	assert.False(t, b.DownButton().Disabled())
	assert.False(t, b.UpButton().Disabled())

	b.SetValue(5.)
	assert.Equal(t, 5., b.Value())
	assert.False(t, b.DownButton().Disabled())
	assert.True(t, b.UpButton().Disabled())
}

func TestSpinnerBase_SetMinMaxStep(t *testing.T) {
	s := &Spinner{}
	b := NewSpinnerBase(s, 1., 5., 1., 0)
	b.SetMinMaxStep(2., 10., 2.)
	assert.Equal(t, 2., b.Value())
	assert.Equal(t, 2., b.data.min)
	assert.Equal(t, 10., b.data.max)
	assert.Equal(t, 2., b.data.step)
	assert.True(t, b.downButton.Disabled())
	assert.False(t, b.upButton.Disabled())

	b.SetMinMaxStep(0., 12., 2.)
	assert.Equal(t, 0., b.Value())
	assert.Equal(t, 0., b.data.min)
	assert.Equal(t, 12., b.data.max)
	assert.Equal(t, 2., b.data.step)
	assert.True(t, b.downButton.Disabled())
	assert.False(t, b.upButton.Disabled())
}
