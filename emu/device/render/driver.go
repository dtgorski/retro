// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package render

type (
	// Renderer is the interface for all rendering modes.
	Renderer interface {
		Render(page byte, canvas []byte, flash bool)
	}

	// Mode is a Renderer type key.
	Mode byte

	// Modes is a Renderer collection indexed by Mode.
	Modes map[Mode]Renderer

	// Driver is the output rendering driver.
	Driver struct {
		modes  Modes
		canvas []byte
		mode   switchMode
	}

	switchMode byte
)

const (
	// Width is the window width.
	Width = 280

	// Height is the window height.
	Height = 192

	// ModeText identifies the default Text mode (Text 40 x 24, monochrome).
	ModeText Mode = 0x00

	// ModeLoRes identifies the low resolution mode (LoRes 40 x 48, 16 colors).
	ModeLoRes Mode = 0x01

	// ModeHiRes identifies the high resolution mode (HiRes 140 x 192, 6/8 colors).
	ModeHiRes Mode = 0x02

	// Soft switches representation.
	switchModeText  switchMode = 0x01
	switchModeMixed switchMode = 0x02
	switchModePage2 switchMode = 0x04
	switchModeHiRes switchMode = 0x08
)

// NewDriver creates a new output rendering driver.
func NewDriver(modes Modes) *Driver {
	d := &Driver{
		modes:  modes,
		canvas: make([]byte, Width*Height*4),
	}
	d.Reset()
	return d
}

// 0xC050-0xC057
func (d *Driver) setMode(lo byte) bool {

	switch lo = lo & 0x07; {
	case lo == 0x00:
		d.mode &= ^switchModeText

	case lo == 0x01:
		d.mode |= switchModeText

	// "W.Gayler - The Apple II Circuit Description", Page 102 claims:
	// "When All-Text (TEXT MODE) is selected, MIXED MODE and HIRES MODE
	// are don't cares." Did the software authors already know it back then? ;)

	case lo == 0x02: // && d.smode&switchModeText == 0:
		d.mode &= ^switchModeMixed

	case lo == 0x03: // && d.smode&switchModeText == 0:
		d.mode |= switchModeMixed

	case lo == 0x04:
		d.mode &= ^switchModePage2

	case lo == 0x05:
		d.mode |= switchModePage2

	case lo == 0x06: // && d.smode&switchModeText == 0:
		d.mode &= ^switchModeHiRes

	case lo == 0x07: // && d.smode&switchModeText == 0:
		d.mode |= switchModeHiRes

	default:
		return false
	}
	return true
}

// Read reads a byte, if this device is sensitive to this address.
func (d *Driver) Read(lo, hi byte) (byte, bool) {
	if hi != 0xC0 || lo < 0x50 || lo > 0x57 {
		return 0, false
	}
	return 0, d.setMode(lo)
}

// Write writes a byte, if this device is sensitive to this address.
func (d *Driver) Write(lo, hi, _ byte) bool {
	if hi != 0xC0 || lo < 0x50 || lo > 0x57 {
		return false
	}
	return d.setMode(lo)
}

// Reset resets the rendering driver.
func (d *Driver) Reset() {
	d.mode = switchModeText
}

// Slot is set by the memory Manager, depending on where this device was mounted.
func (*Driver) Slot(byte) {}

// Render delegates rendering to the current renderer.
func (d *Driver) Render(flash bool) []byte {

	page := byte(d.mode>>2) & 0x01

	// All text?
	if d.mode&switchModeText != 0 {
		d.modes[ModeText].Render(page, d.canvas, flash)
		return d.canvas
	}

	// Graphics?
	if d.mode&switchModeHiRes != 0 {
		d.modes[ModeHiRes].Render(page, d.canvas, false)
	} else {
		d.modes[ModeLoRes].Render(page, d.canvas, false)
	}

	// Mixed?
	if d.mode&switchModeMixed != 0 {
		d.modes[ModeText].(*Text).Mixed(page, d.canvas, flash)
	}

	return d.canvas
}
