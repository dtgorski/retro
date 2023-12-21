// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package input

type (
	// MouseButton is emitted by the mouse.
	MouseButton struct {
		b int
		a int
		m int
	}
)

// NewMouseButton creates a new mouse event.
func NewMouseButton(button, action, modifier int) MouseButton {
	if button < 0 || button > 2 {
		button = 0
	}
	return MouseButton{button, action, modifier}
}

// Button returns [0..2] according to the button used.
func (b MouseButton) Button() int {
	return b.b
}

// IsButton0 signals if left button is used.
func (b MouseButton) IsButton0() bool {
	return b.Button() == 0
}

// IsPressed signals if this button is pressed (released otherwise).
func (b MouseButton) IsPressed() bool {
	return b.a == 1
}
