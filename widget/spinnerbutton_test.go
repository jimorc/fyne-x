package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestSpinnerButton_NewSpinnerButton(t *testing.T) {
	var count uint32 = 0
	b := newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), func() { count++ })
	assert.Less(t, float32(0), b.MinSize().Width)

	b.Tapped(nil)
	assert.Equal(t, uint32(1), count)
	b.Tapped(nil)
	assert.Equal(t, uint32(2), count)
}

func TestSpinnerButton_Move(t *testing.T) {
	b := newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), nil)
	_ = container.NewWithoutLayout(b)
	b.Move(fyne.NewPos(10, 20))

	assert.Equal(t, float32(10), b.position.X)
	assert.Equal(t, float32(10), b.Position().X)
	assert.Equal(t, float32(20), b.position.Y)
	assert.Equal(t, float32(20), b.Position().Y)
}

func TestSpinnerButton_ContainsPoint(t *testing.T) {
	b := newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), nil)
	_ = container.NewWithoutLayout(b)
	b.Move(fyne.NewPos(10, 20))

	assert.True(t, b.ContainsPoint(fyne.NewPos(10, 20)))
	assert.True(t, b.ContainsPoint(fyne.NewPos(10+b.size.Width, 20+b.size.Height)))
	assert.False(t, b.ContainsPoint(fyne.NewPos(9, 20)))
	assert.False(t, b.ContainsPoint(fyne.NewPos(10, 20+b.size.Height+1)))
}

func TestSpinnerButton_EnableDisable(t *testing.T) {
	b := newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), nil)
	b.EnableDisable(false, false)
	assert.False(t, b.Disabled())

	b.EnableDisable(true, false)
	assert.True(t, b.Disabled())

	b.EnableDisable(false, true)
	assert.True(t, b.Disabled())

	b.EnableDisable(true, true)
	assert.True(t, b.Disabled())
}
