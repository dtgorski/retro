// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package builtin

import (
	"retro/emu/memory"
)

type (
	// Keyboard handles I/O soft switches 0xC000 and 0xC010.
	Keyboard struct {
		mem memory.Memory
	}
)

// NewKeyboard creates a new Keyboard device.
func NewKeyboard(mem memory.Memory) memory.Device {
	return &Keyboard{mem: mem}
}

// Read reads a byte, if this device is sensitive to this address.
func (p *Keyboard) Read(lo, hi byte) (byte, bool) {
	if hi == 0xC0 && lo == 0x10 { // KBDSTRB
		p.mem.DMA()[0xC000] = p.mem.DMA()[0xC000] & 0x7F
		return 0, true
	}
	return 0, false
}

// Write writes a byte, if this device is sensitive to this address.
func (p *Keyboard) Write(lo, hi, _ byte) bool {
	if hi == 0xC0 && lo == 0x10 { // KBDSTRB
		p.mem.DMA()[0xC000] = p.mem.DMA()[0xC000] & 0x7F
		return true
	}
	return false
}

// Reset does nothing here.
func (*Keyboard) Reset() {}

// Slot is set by the memory Manager, depending on where this device was mounted.
func (*Keyboard) Slot(byte) {}
