package widget

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func waitForBinding() {
	time.Sleep(time.Millisecond * 100) // data resolves on background thread
}

func TestNewSpinner(t *testing.T) {
	s := NewSpinner(1., 5., 1.5, 0, nil)
	assert.Equal(t, 1., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestNewSpinner_BadArgs(t *testing.T) {
	s := NewSpinner(5., 5., 1., 0, nil)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when max = min")

	s = NewSpinner(5., 4., 1., 0, nil)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when min > max")

	s = NewSpinner(1., 5., 0., 0, nil)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step = 0")

	s = NewSpinner(1., 5., -5., 0, nil)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step < 0")

	s = NewSpinner(1., 5., 5., 0, nil)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step > max - min")

	s = NewSpinner(1., 5., 2., 11, nil)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.base.format)
	assert.True(t, s.base.Initialized())
}

func TestNewSpinnerWithData(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(1., 5., 2., 0, data)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)

	s.base.SetValue(1.52)
	waitForBinding()
	val, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1.52, val)

	err = data.Set(3.1)
	assert.NoError(t, err)
	waitForBinding()
	assert.Equal(t, 3.1, s.base.Value())
}

func TestSpinner_Unbind(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(1., 5., 2., 0, data)
	waitForBinding()
	s.Unbind()
	s.SetValue(2.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)
}

func TestNewSpinnerWithData_BadArgs(t *testing.T) {
	boundValue := binding.NewFloat()
	s := NewSpinnerWithData(5., 5., 1., 0, boundValue)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when max = min")

	s = NewSpinnerWithData(5., 4., 1., 0, boundValue)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when min > max")

	s = NewSpinnerWithData(1., 5., 0., 0, boundValue)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step = 0")

	s = NewSpinnerWithData(1., 5., -5., 0, boundValue)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step < 0")

	s = NewSpinnerWithData(1., 5., 5., 0, boundValue)
	assert.False(t, s.base.Initialized(), "spinner should not be initialized when step > max - min")

	s = NewSpinnerWithData(1., 5., 2., 11, boundValue)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.base.format)
	assert.True(t, s.base.Initialized())
}

func TestNewSpinnerUninitialized(t *testing.T) {
	s := NewSpinnerUninitialized(0)
	assert.False(t, s.base.Initialized())
	assert.True(t, s.Disabled())
	s.Enable()
	assert.True(t, s.Disabled())
	s.SetMinMaxStep(-1., 2., 1.1)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())

	assert.Equal(t, "%d", s.base.format)
	assert.True(t, s.base.Initialized())

	s = NewSpinnerUninitialized(4)
	assert.False(t, s.base.Initialized())
	assert.Equal(t, "%.4f", s.base.format)

	s = NewSpinnerUninitialized(maxDecimals + 2)
	assert.False(t, s.base.Initialized())
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.base.format)
}

func TestSpinner_SetValue(t *testing.T) {
	s := NewSpinner(1, 5, 2, 0, nil)
	s.SetValue(2)
	assert.Equal(t, 2., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_SetValue_LessThanMin(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.SetValue(3)
	assert.Equal(t, 4., s.Value())
	assert.True(t, s.base.DownButton().Disabled())
	assert.False(t, s.base.UpButton().Disabled())
}

func TestSpinner_SetValue_GreaterThanMax(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.SetValue(23.)
	assert.Equal(t, 22., s.Value())
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_SetValue_DisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.Disable()
	s.SetValue(10.)
	assert.Equal(t, 4., s.Value())
	assert.True(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_SetMinMaxStep(t *testing.T) {
	s := NewSpinner(1., 6., 2., 0, nil)
	s.SetMinMaxStep(0., 10., 1.)
	assert.Equal(t, 0., s.base.data.min)
	assert.Equal(t, 10., s.base.data.max)
	assert.Equal(t, 1., s.base.data.step)
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_SetMinMaxStep_BadArgs(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	s.SetMinMaxStep(11, 10, 2)
	assert.NotNil(t, s.base.Validate())
	assert.Equal(t, 10., s.base.Value())
	s.SetMinMaxStep(1, 10, 10)
	assert.NotNil(t, s.base.Validate())
	assert.Equal(t, 1., s.base.Value())
	s.SetMinMaxStep(1, 10, -1)
	assert.NotNil(t, s.base.Validate())
	assert.Equal(t, 1., s.base.Value())
}

func TestSpinner_SetMinMaxStep_OutsideRange(t *testing.T) {
	s := NewSpinner(-2, 20, 1, 0, nil)
	s.SetValue(19.)
	s.SetMinMaxStep(-1., 10., 1.2)
	assert.Equal(t, -1., s.base.Value())
	s.SetValue(-1.)
	s.SetMinMaxStep(1., 10., 1.)
	assert.Equal(t, 1., s.base.Value())
}

func TestNewSpinnerSpinner_SetMinMaxStep_DataAboveRange(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(-2, 20, 1, 0, data)
	data.Set(19.)
	waitForBinding()
	assert.Equal(t, 19., s.Value())
	s.SetMinMaxStep(-1., 10., 1.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1., val)
}

func TestSpinner_SetMinMaxStep_DataBelowRange(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(-2, 20, 1, 0, data)
	data.Set(-2.)
	waitForBinding()
	assert.Equal(t, -2., s.Value())
	s.SetMinMaxStep(-1., 10., 1.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1., val)
}

func TestSpinner_UpButtonTapped(t *testing.T) {
	s := NewSpinner(4., 10., 5., 0, nil)
	s.base.UpButton().Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9., s.Value())
	s.base.UpButton().Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10., s.Value())
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_DownButtonTapped(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(10.)
	s.base.DownButton().Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5., s.Value())
	s.base.DownButton().Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_EnableDisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(7.)
	s.Disable()
	assert.True(t, s.Disabled())
	assert.True(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_EnableDisabledSpinner_UpButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(10.)
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
	s.Disable()
	assert.True(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())

	s.Enable()
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_EnableDisabledSpinner_DownButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(4.)
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
	s.Disable()
	assert.True(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())

	s.Enable()
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_RunePlus(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.TypedRune('+')
	assert.Equal(t, 9., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())

	s.TypedRune('+')
	assert.Equal(t, 10., s.Value())
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_RuneMinus(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(10.)
	s.TypedRune('-')
	assert.Equal(t, 5., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())

	s.TypedRune('-')
	assert.Equal(t, 4., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_RunePlus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.Disable()
	s.TypedRune('+')
	assert.Equal(t, 4., s.Value())
}

func TestSpinner_RuneMinus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	s.TypedRune('-')
	assert.Equal(t, 8., s.Value())
}

func TestSpinner_RunePlus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = false
	s.TypedRune('+')
	assert.Equal(t, 4., s.Value())
}

func TestSpinner_RuneMinus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	s.TypedRune('-')
	assert.Equal(t, 8., s.Value())
}

func TestSpinner_KeyUp(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 9., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 10., s.Value())
	assert.True(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())
}

func TestSpinner_KeyDown(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(10)
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 5., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.False(t, s.base.DownButton().Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 4., s.Value())
	assert.False(t, s.base.UpButton().Disabled())
	assert.True(t, s.base.DownButton().Disabled())
}

func TestSpinner_KeyUp_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.Value())
}

func TestSpinner_KeyDown_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.Value())
}

func TestSpinner_KeyUp_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.Value())
}

func TestSpinner_KeyDown_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.Value())
}

func TestSpinner_Scrolled(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	s.focused = true
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.Value())
	s.Scrolled(&e)
	assert.Equal(t, 3., s.Value())
	delta = fyne.Delta{DX: 0, DY: -25}
	e = fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.Value())
}

func TestSpinner_Scrolled_Disabled(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	s.focused = true
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.Value())
}

func TestSpinner_Scrolled_NotFocused(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.Value())
}

/*
	func TestSpinner_OnChanged(t *testing.T) {
		var v float64
		s := NewSpinner(1, 10, 1, 0, func(newVal float64) {
			v = newVal
		})
		s.SetValue(3.)
		assert.Equal(t, 3., v)
	}
*/
/*
func TestSpinner_OnChanged_Disabled(t *testing.T) {
	var v float64
	s := NewSpinner(1, 10, 1, 0, func(newVal float64) {
		v = newVal
	})
	s.Disable()
	s.SetValue(3.)
	assert.Equal(t, 1., v)
}
*/
func TestSpinner_Binding_OutsideRange(t *testing.T) {
	val := binding.NewFloat()
	s := NewSpinnerWithData(1, 5, 2, 0, val)
	waitForBinding()
	err := val.Set(7)
	assert.NoError(t, err)
	waitForBinding()

	assert.Equal(t, 5., s.Value())

	v, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5., v)
}
