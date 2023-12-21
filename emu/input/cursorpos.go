// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package input

type (
	// CursorPos is emitted by the mouse.
	CursorPos struct {
		x float64
		y float64
	}
)

// NewCursorPos creates a new position event.
func NewCursorPos(x, y float64) CursorPos {
	return CursorPos{x, y}
}

// X returns the x coordinate.
func (p CursorPos) X() float64 {
	return p.x
}

// Y returns the y coordinate.
func (p CursorPos) Y() float64 {
	return p.y
}
