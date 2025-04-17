package widget

import (
	"errors"
	"strconv"
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

func TestSpinnerEntry_validate(t *testing.T) {
	testCases := []struct {
		name        string
		allowFloat  bool
		min         float64
		max         float64
		text        string
		expectedErr error
	}{
		{
			name:        "Valid integer within range",
			allowFloat:  false,
			min:         0,
			max:         10,
			text:        "5",
			expectedErr: nil,
		},
		{
			name:        "Valid float within range",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "5.5",
			expectedErr: nil,
		},
		{
			name:        "Integer equal to min",
			allowFloat:  false,
			min:         0,
			max:         10,
			text:        "0",
			expectedErr: nil,
		},
		{
			name:        "Integer equal to max",
			allowFloat:  false,
			min:         0,
			max:         10,
			text:        "10",
			expectedErr: nil,
		},
		{
			name:        "Float equal to min",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "0.0",
			expectedErr: nil,
		},
		{
			name:        "Float equal to max",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "10.0",
			expectedErr: nil,
		},
		{
			name:        "Integer below min",
			allowFloat:  false,
			min:         1,
			max:         10,
			text:        "0",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Integer above max",
			allowFloat:  false,
			min:         0,
			max:         9,
			text:        "10",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Float below min",
			allowFloat:  true,
			min:         1.0,
			max:         10.0,
			text:        "0.9",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Float above max",
			allowFloat:  true,
			min:         0.0,
			max:         9.0,
			text:        "9.1",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Invalid integer input",
			allowFloat:  false,
			min:         0,
			max:         10,
			text:        "5a",
			expectedErr: errors.New("value is not a number"),
		},
		{
			name:        "Invalid float input",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "5.5b",
			expectedErr: errors.New("value is not a number"),
		},
		{
			name:        "Empty input",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "",
			expectedErr: errors.New("value is not a number"),
		},
		{
			name:        "Large float value",
			allowFloat:  true,
			min:         0.0,
			max:         1e10,
			text:        "1e11",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Large int value",
			allowFloat:  false,
			min:         0,
			max:         1000,
			text:        strconv.Itoa(1001),
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Negative float value with min 0",
			allowFloat:  true,
			min:         0.0,
			max:         10.0,
			text:        "-1.0",
			expectedErr: errors.New("value is not between min and max"),
		},
		{
			name:        "Negative int value with min 0",
			allowFloat:  false,
			min:         0,
			max:         10,
			text:        "-1",
			expectedErr: errors.New("value is not between min and max"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := &spinnerEntry{spinner: &Spinner{min: tc.min, max: tc.max}}
			e.AllowFloat = tc.allowFloat
			err := e.validate(tc.text)

			if tc.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error: %v, but got nil", tc.expectedErr)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error: %v, but got: %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestSpinner_Enable(t *testing.T) {
	s := NewSpinner(3, 10, 1, 1)
	s.initialized = true
	s.Disable()

	s.Enable()

	if s.Disabled() {
		t.Error("Spinner should be enabled")
	}
	if s.entry.Disabled() {
		t.Error("Entry should be enabled")
	}
	if s.upButton.Disabled() {
		t.Error("Up button should be enabled")
	}
	if s.downButton.Disabled() {
		t.Error("Down button should be enabled")
	}
}

func TestSpinner_Enable_Uninitialized(t *testing.T) {
	s := NewSpinnerUninitialized(2)
	s.Enable()
	if !s.Disabled() {
		t.Error("Spinner should be disabled")
	}
	if !s.entry.Disabled() {
		t.Error("Entry should be disabled")
	}
	if !s.upButton.Disabled() {
		t.Error("Up button should be disabled")
	}
	if !s.downButton.Disabled() {
		t.Error("Down button should be disabled")
	}
	if s.initialized {
		t.Error("Spinner should not be initialized")
	}
}
