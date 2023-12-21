// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package input

type (
	// KeyInput is a real or simulated keyboard stroke.
	KeyInput struct {
		key int
		act int
		mod int
	}
)

// NewKeyInput create a new keyboard input event.
func NewKeyInput(key, action, modifier int) KeyInput {
	return KeyInput{key, action, modifier}
}

// IsCtrlShiftR signals when CTRL-SHIFT-R is pressed.
func (e KeyInput) IsCtrlShiftR() bool {
	return e.key == 0x52 && (e.act == 1 || e.act == 2) && e.mod == 3
}

// IsCtrlV signals when CTRL-V is pressed.
func (e KeyInput) IsCtrlV() bool {
	return e.key == 0x56 && (e.act == 1 || e.act == 2) && e.mod == 2
}
