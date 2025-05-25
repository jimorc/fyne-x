package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// spinnerButton is the widget used for buttons in Spinners.
type SpinnerButton struct {
	widget.Button

	position fyne.Position
	size     fyne.Size
}

// newSpinnerButton creates a spinnerButton for use in Spinner widgets.
//
// Params:
//
//	resource is the resource to be used as the button icon.
//	onTapped is the callback function for button clicks.
func newSpinnerButton(resource fyne.Resource, onTapped func()) *SpinnerButton {
	b := &SpinnerButton{}
	b.ExtendBaseWidget(b)
	b.setButtonProperties(resource, onTapped)
	return b
}

// MinSize returns the minimum size of the button. Because the minimum size is a constant
// based on the spinner height and theme properties, the minimum size is calculated when
// the button is created.
func (b *SpinnerButton) MinSize() fyne.Size {
	return fyne.NewSize(b.size.Height, b.size.Height)
}

// Move moves the button.
func (b *SpinnerButton) Move(pos fyne.Position) {
	b.position = pos
	b.BaseWidget.Move(pos)
}

// enableDisable enables or disables the button based on whether the button's
// parent spinner widget is disabled, and whether the spinner's value is at its limit.
//
// Params:
//
//		parentDisabled indicates whether the button's parent spinner is disabled or not.
//		limit indicates whether the spinner's value is at the corresponding limit for this
//	 button. For example, for an up button, the  limit should be true if value == max, and
//	 for a down button, the limit should be true if value == min.
func (b *SpinnerButton) EnableDisable(parentDisabled, limit bool) {
	if parentDisabled {
		b.Disable()
	} else {
		b.Enable()
		if limit {
			b.Disable()
		}
	}
}

// setButtonProperties sets the button properties.
//
// Params:
//
//	resource is the Resource for the button icon.
//	onTapped is the function to be called when the button is tapped.
func (b *SpinnerButton) setButtonProperties(resource fyne.Resource, onTapped func()) {
	b.Icon = resource
	b.Text = ""
	b.OnTapped = onTapped

	// calculate the minimum button size (really just its height).
	th := b.Theme()
	padding := fyne.NewSquareSize(th.Size(theme.SizeNameInnerPadding) * 2)
	text := canvas.NewText("0", color.Black)
	textSize, _ := fyne.CurrentApp().Driver().RenderedTextSize(text.Text,
		text.TextSize, text.TextStyle, text.FontSource)
	tHeight := textSize.Height + padding.Height

	h := tHeight/2 - th.Size(theme.SizeNameInputBorder) - 2
	b.size = fyne.NewSize(h, h)

}
