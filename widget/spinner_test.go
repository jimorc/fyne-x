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
	assert.Equal(t, 1., s.data.min)
	assert.Equal(t, 5., s.data.max)
	assert.Equal(t, 1.5, s.data.step)
	assert.Equal(t, "%d", s.data.format)
	assert.Equal(t, 1., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s = NewSpinner(1., 5., 1.5, 3, nil)
	assert.Equal(t, "%.3f", s.data.format)

}

func TestNewSpinner_BadArgs(t *testing.T) {
	s := NewSpinner(5., 5., 1., 0, nil)
	assert.False(t, s.data.initialized, "spinner should not be initialized when max = min")

	s = NewSpinner(5., 4., 1., 0, nil)
	assert.False(t, s.data.initialized, "spinner should not be initialized when min > max")

	s = NewSpinner(1., 5., 0., 0, nil)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step = 0")

	s = NewSpinner(1., 5., -5., 0, nil)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step < 0")

	s = NewSpinner(1., 5., 5., 0, nil)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step > max - min")

	s = NewSpinner(1., 5., 2., 11, nil)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.data.format)
	assert.True(t, s.data.initialized)
}

func TestNewSpinnerWithData(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(1., 5., 2., 0, data)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1., val)

	s.data.SetValue(1.52)
	waitForBinding()
	val, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1.52, val)

	err = data.Set(3.1)
	assert.NoError(t, err)
	waitForBinding()
	assert.Equal(t, 3.1, s.data.Value())
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
	assert.False(t, s.data.initialized, "spinner should not be initialized when max = min")

	s = NewSpinnerWithData(5., 4., 1., 0, boundValue)
	assert.False(t, s.data.initialized, "spinner should not be initialized when min > max")

	s = NewSpinnerWithData(1., 5., 0., 0, boundValue)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step = 0")

	s = NewSpinnerWithData(1., 5., -5., 0, boundValue)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step < 0")

	s = NewSpinnerWithData(1., 5., 5., 0, boundValue)
	assert.False(t, s.data.initialized, "spinner should not be initialized when step > max - min")

	s = NewSpinnerWithData(1., 5., 2., 11, boundValue)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.data.format)
	assert.True(t, s.data.initialized)
}

func TestNewSpinnerUninitialized(t *testing.T) {
	s := NewSpinnerUninitialized(0)
	assert.False(t, s.data.initialized)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.True(t, s.Disabled())
	s.SetMinMaxStep(-1., 2., 1.1)
	assert.True(t, s.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())

	assert.Equal(t, "%d", s.data.format)
	assert.True(t, s.data.initialized)

	s = NewSpinnerUninitialized(4)
	assert.False(t, s.data.initialized)
	assert.Equal(t, "%.4f", s.data.format)

	s = NewSpinnerUninitialized(maxDecimals + 2)
	assert.False(t, s.data.initialized)
	assert.Equal(t, fmt.Sprintf("%%.%df", maxDecimals), s.data.format)
}

func TestSpinner_SetValue(t *testing.T) {
	s := NewSpinner(1, 5, 2, 0, nil)
	s.SetValue(2)
	assert.Equal(t, 2., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_SetValue_LessThanMin(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.SetValue(3)
	assert.Equal(t, 4., s.GetValue())
	assert.True(t, s.downButton.Disabled())
	assert.False(t, s.upButton.Disabled())
}

func TestSpinner_SetValue_GreaterThanMax(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.SetValue(23.)
	assert.Equal(t, 22., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_SetValue_DisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 22, 5, 0, nil)
	s.Disable()
	s.SetValue(10.)
	assert.Equal(t, 4., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_SetMinMaxStep(t *testing.T) {
	s := NewSpinner(1., 6., 2., 0, nil)
	s.SetMinMaxStep(0., 10., 1.)
	assert.Equal(t, 0., s.data.min)
	assert.Equal(t, 10., s.data.max)
	assert.Equal(t, 1., s.data.step)
}

func TestSpinner_SetMinMaxStep_BadArgs(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	assert.Panics(t, func() { s.SetMinMaxStep(11, 10, 2) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, 10) })
	assert.Panics(t, func() { s.SetMinMaxStep(1, 10, -1) })
}

func TestSpinner_SetMinMaxStep_OutsideRange(t *testing.T) {
	s := NewSpinner(-2, 20, 1, 0, nil)
	s.SetValue(19.)
	s.SetMinMaxStep(-1., 10., 1.2)
	assert.Equal(t, 10., s.GetValue())
	s.SetValue(-1.)
	s.SetMinMaxStep(1., 10., 1.)
	assert.Equal(t, 1., s.GetValue())
}

/*
	func TestNewSpinnerSpinner_SetMinMaxStep_DataAboveRange(t *testing.T) {
		data := binding.NewFloat()
		s := NewSpinnerWithData(-2, 20, 1, 0, data)
		data.Set(19.)
		waitForBinding()
		assert.Equal(t, 19., s.GetValue())
		s.SetMinMaxStep(-1., 10., 1.)
		waitForBinding()
		val, err := data.Get()
		assert.NoError(t, err)
		assert.Equal(t, 10., val)
	}
*/
/*
func TestSpinner_SetMinMaxStep_DataBelowRange(t *testing.T) {
	data := binding.NewFloat()
	s := NewSpinnerWithData(-2, 20, 1, 0, data)
	data.Set(-2.)
	waitForBinding()
	assert.Equal(t, -2., s.GetValue())
	s.SetMinMaxStep(-1., 10., 1.)
	waitForBinding()
	val, err := data.Get()
	assert.NoError(t, err)
	assert.Equal(t, -1., val)
}
*/
func TestSpinner_UpButtonTapped(t *testing.T) {
	s := NewSpinner(4., 10., 5., 0, nil)
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 9., s.GetValue())
	s.upButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_DownButtonTapped(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(10.)
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 5., s.GetValue())
	s.downButton.Tapped(&fyne.PointEvent{})
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_EnableDisabledSpinner(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(7.)
	s.Disable()
	assert.True(t, s.Disabled())
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Enable()
	assert.False(t, s.Disabled())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_EnableDisabledSpinner_UpButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(10.)
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_EnableDisabledSpinner_DownButtonDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.SetValue(4.)
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
	s.Disable()
	assert.True(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())

	s.Enable()
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_RunePlus(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.TypedRune('+')
	assert.Equal(t, 9., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('+')
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_RuneMinus(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(10.)
	s.TypedRune('-')
	assert.Equal(t, 5., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedRune('-')
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_RunePlus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.Disable()
	s.TypedRune('+')
	assert.Equal(t, 4., s.GetValue())
}

func TestSpinner_RuneMinus_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	s.TypedRune('-')
	assert.Equal(t, 8., s.GetValue())
}

func TestSpinner_RunePlus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = false
	s.TypedRune('+')
	assert.Equal(t, 4., s.GetValue())
}

func TestSpinner_RuneMinus_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8)
	s.focused = false
	s.TypedRune('-')
	assert.Equal(t, 8., s.GetValue())
}

func TestSpinner_KeyUp(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 9., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 10., s.GetValue())
	assert.True(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())
}

func TestSpinner_KeyDown(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(10)
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 5., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.False(t, s.downButton.Disabled())

	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
	assert.False(t, s.upButton.Disabled())
	assert.True(t, s.downButton.Disabled())
}

func TestSpinner_KeyUp_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
}

func TestSpinner_KeyDown_SpinnerDisabled(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.Disable()
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.GetValue())
}

func TestSpinner_KeyUp_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyUp}
	s.TypedKey(&key)
	assert.Equal(t, 4., s.GetValue())
}

func TestSpinner_KeyDown_SpinnerNotFocused(t *testing.T) {
	s := NewSpinner(4, 10, 5, 0, nil)
	s.focused = true
	s.SetValue(8.)
	s.focused = false
	key := fyne.KeyEvent{Name: fyne.KeyDown}
	s.TypedKey(&key)
	assert.Equal(t, 8., s.GetValue())
}

func TestSpinner_Scrolled(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	s.focused = true
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.GetValue())
	s.Scrolled(&e)
	assert.Equal(t, 3., s.GetValue())
	delta = fyne.Delta{DX: 0, DY: -25}
	e = fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 2., s.GetValue())
}

func TestSpinner_Scrolled_Disabled(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	s.focused = true
	s.Disable()
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.GetValue())
}

func TestSpinner_Scrolled_NotFocused(t *testing.T) {
	s := NewSpinner(1, 10, 1, 0, nil)
	delta := fyne.Delta{DX: 0, DY: 25}
	e := fyne.ScrollEvent{Scrolled: delta}
	s.Scrolled(&e)
	assert.Equal(t, 1., s.GetValue())
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

	assert.Equal(t, 5., s.GetValue())

	v, err := val.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5., v)
}
