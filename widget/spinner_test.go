package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
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
	assert.Equal(t, float64(1), s.value)
	assert.True(t, s.initialized)

	s2 := NewSpinner(-5, 5, 1, 2)
	assert.Equal(t, float64(-5), s2.min)
	assert.Equal(t, float64(5), s2.max)
	assert.Equal(t, float64(1), s2.step)
	assert.Equal(t, uint(2), s2.decimalPlaces)
	assert.True(t, s2.entry.AllowNegative)
	assert.True(t, s2.entry.AllowFloat)
	assert.Equal(t, float64(-5), s2.value)
	assert.True(t, s2.initialized)
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

func TestNewSpinnerUninitialized(t *testing.T) {
	s := NewSpinnerUninitialized(0)
	assert.Equal(t, uint(0), s.decimalPlaces)
	assert.False(t, s.initialized)
	assert.True(t, s.Disabled())

	s2 := NewSpinnerUninitialized(2)
	assert.Equal(t, uint(2), s2.decimalPlaces)
	assert.False(t, s2.initialized)
	assert.True(t, s2.Disabled())
}

func TestSpinnerSetValue(t *testing.T) {
	s := NewSpinner(1, 10, 2, 0)
	assert.Equal(t, float64(1), s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.SetValue(5)
	assert.Equal(t, float64(5), s.value)
	assert.Equal(t, "5", s.entry.Text)
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
	s.SetValue(0)
	assert.Equal(t, float64(1), s.value)
	assert.Equal(t, "1", s.entry.Text)
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.SetValue(11)
	assert.Equal(t, float64(10), s.value)
	assert.Equal(t, "10", s.entry.Text)
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
	s.SetValue(-1)
	assert.Equal(t, float64(1), s.value)
	assert.Equal(t, "1", s.entry.Text)
	assert.True(t, s.downButton.Disabled())
	assert.False(t, s.upButton.Disabled())
}

func TestSpinnerSetValue_Disabled(t *testing.T) {
	s := NewSpinner(1, 10, 2, 0)
	s.Disable()
	s.SetValue(5)
	assert.Equal(t, float64(1), s.value)
	assert.Equal(t, "1", s.entry.Text)
}

func TestSpinner_UpButtonTapped(t *testing.T) {
	s := NewSpinner(4., 10., 5., 0)
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9., s.GetValue())
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_DownButtonTapped(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0)
	s.SetValue(10.)
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5., s.GetValue())
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_Disable(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0)
	s.Disable()
	assert.True(t, s.Disabled())
	assert.True(t, s.entry.Disabled())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_Enable(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0)
	s.Disable()
	s.Enable()
	assert.False(t, s.Disabled())
	assert.False(t, s.entry.Disabled())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinnerUninitialized_Enable(t *testing.T) {
	s := NewSpinnerUninitialized(0)
	s.Enable()
	assert.True(t, s.Disabled())
	assert.True(t, s.entry.Disabled())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinnerUninitialized_SetMinMaxStep(t *testing.T) {
	s := NewSpinnerUninitialized(0)
	s.SetMinMaxStep(4, 10, 5)
	assert.Equal(t, float64(4), s.min)
	assert.Equal(t, float64(10), s.max)
	assert.Equal(t, float64(5), s.step)
	assert.Equal(t, uint(0), s.decimalPlaces)
	assert.False(t, s.entry.AllowFloat)
	assert.False(t, s.entry.AllowNegative)
	assert.Equal(t, float64(4), s.value)
	assert.True(t, s.initialized)
	assert.False(t, s.Disabled())
}

func TestSpinner_SetMinMaxStep(t *testing.T) {
	s := NewSpinner(4, 10, 5, 1)
	s.SetMinMaxStep(-1, 5, 0.75)
	assert.Equal(t, float64(-1), s.min)
	assert.Equal(t, float64(5), s.max)
	assert.Equal(t, float64(0.75), s.step)
	assert.Equal(t, uint(1), s.decimalPlaces)
	assert.True(t, s.entry.AllowFloat)
	assert.True(t, s.entry.AllowNegative)
	assert.Equal(t, float64(-1), s.value)
	assert.True(t, s.initialized)
	assert.False(t, s.Disabled())
}

func TestSpinnerEntry_Validator(t *testing.T) {
	s := NewSpinner(4, 10, 5, 1)
	assert.Nil(t, s.entry.Validator("4"))
	assert.Nil(t, s.entry.Validator("4.5"))
	assert.Nil(t, s.entry.Validator("5.123456"))
	err := s.entry.Validator("-s.5")
	assert.NotNil(t, err)
	assert.Equal(t, "value is not a number", err.Error())
	err = s.entry.Validator("11")
	assert.NotNil(t, err)
	assert.Equal(t, "value is not between min and max", err.Error())
	err = s.entry.Validator("3")
	assert.NotNil(t, err)
	assert.Equal(t, "value is not between min and max", err.Error())
}

func TestSpinnerEntryUpKey(t *testing.T) {
	s := NewSpinner(4, 10, 5, 1)
	s.entry.Tapped(&fyne.PointEvent{})
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Equal(t, float64(9), s.GetValue())
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Equal(t, float64(10), s.GetValue())
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	assert.Equal(t, float64(10), s.GetValue())
}

func TestSpinnerEntryDownKey(t *testing.T) {
	s := NewSpinner(1, 10, 5, 1)
	s.SetValue(8)
	s.entry.Tapped(&fyne.PointEvent{})
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Equal(t, float64(3), s.GetValue())
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Equal(t, float64(1), s.GetValue())
	s.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	assert.Equal(t, float64(1), s.GetValue())
}

func TestSpinnerEntryTypedShortcut(t *testing.T) {
	s := NewSpinner(1, 10, 5, 1)
	s.entry.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: fyne.KeyModifierControl})
	assert.Equal(t, float64(1), s.GetValue())
	s.entry.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyUp, Modifier: fyne.KeyModifierControl})
	assert.Equal(t, float64(10), s.GetValue())
}
