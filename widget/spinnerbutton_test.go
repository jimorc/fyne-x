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

func TestSpinnerButton_EnableDisable(t *testing.T) {
	b := newSpinnerButton(theme.Icon(theme.IconNameArrowDropUp), nil)
	b.enableDisable(false, false)
	assert.False(t, b.Disabled())

	b.enableDisable(true, false)
	assert.True(t, b.Disabled())

	b.enableDisable(false, true)
	assert.True(t, b.Disabled())

	b.enableDisable(true, true)
	assert.True(t, b.Disabled())
}
