// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package builtin

import (
	"retro/emu/memory"
)

type (
	// Paddles handles an analog paddle devices.
	Paddles struct {
		mem memory.Memory
	}
)

// NewPaddle creates a new analog paddle device driver.
func NewPaddle(mem memory.Memory) memory.Device {
	return &Paddles{mem: mem}
}

// Read reads a byte, if this device is sensitive to this address.
func (p *Paddles) Read(lo, hi byte) (byte, bool) {
	if hi == 0xC0 && lo >= 0x61 && lo <= 0x63 { // BUTN0, BUTN1, BUTN2
		b := p.mem.DMA()[0xC000|int(lo)]
		p.mem.DMA()[0xC000|int(lo)] = 0x00
		return b, true
	}

	return 0, false
}

// Write writes a byte, if this device is sensitive to this address.
func (p *Paddles) Write(lo, hi, b byte) bool {
	if hi == 0xC0 && lo >= 0x61 && lo <= 0x63 { // BUTN0, BUTN1, BUTN2
		p.mem.DMA()[0xC000|int(lo)] = b
		return true
	}
	return false
}

// Reset does nothing here.
func (*Paddles) Reset() {}

// Slot is set by the memory Manager, depending on where this device was mounted.
func (*Paddles) Slot(byte) {}
