// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package language

import (
	"retro/emu/memory"
)

type (
	// Card (Language Card) provides additional RAM.
	Card struct {
		ram   memory.Memory
		romIN bool // true = ROM IN / false = RAM IN
		ramRW bool // true = RAM RW / false = RAM RO
		bank  byte
		last  byte
	}
)

// NewCard creates a Language Card.
func NewCard() *Card {

	// Indeed we only need a portion of it, but its less convoluted this way.
	ram := memory.NewMemory()

	return &Card{ram: ram, romIN: true, ramRW: true}
}

// Read reads a byte, if this device is sensitive to this address.
func (c *Card) Read(lo, hi byte) (byte, bool) {

	// RAM: 0xE000-0xFFFF
	if !c.romIN && hi >= 0xE0 {
		return c.ram.Read(lo, hi), true
	}
	// RAM: 0xD000-0xDFFF
	if !c.romIN && hi >= 0xD0 {
		return c.ram.Read(lo, hi-c.bank), true
	}
	// ROM: 0xD000-0xFFFF
	if c.romIN && hi >= 0xD0 {
		return 0, false
	}
	// Soft switch?
	if hi != 0xC0 || lo < 0x80 || lo > 0x8F {
		return 0, false
	}

	// Bit 3 -> Bank 0/1 offset
	c.bank = (lo << 1) & 0x10

	switch lo &= 0x03; {
	case lo == 0x00:
		c.romIN = false
		c.ramRW = false
	case lo == 0x01:
		c.romIN = true
		c.ramRW = c.last&0x01 == 0x01
	case lo == 0x02:
		c.romIN = true
		c.ramRW = false
	case lo == 0x03:
		c.romIN = false
		c.ramRW = c.last&0x01 == 0x01
	}
	c.last = lo

	// Return status? TODO: Is this required?
	// Peek into the RAM Card (by MS) book again.

	return 0, true
}

// Write writes a byte, if this device is sensitive to this address.
func (c *Card) Write(lo, hi, b byte) bool {

	// RAM: 0xE000-0xFFFF
	if c.ramRW && hi >= 0xE0 {
		c.ram.Write(lo, hi, b)
		return true
	}
	// RAM: 0xD000-0xE000
	if c.ramRW && hi >= 0xD0 {
		c.ram.Write(lo, hi-c.bank, b)
		return true
	}
	// ROM: 0xD000-0xFFFF
	if c.romIN && hi >= 0xD0 {
		return false
	}
	// Soft switch?
	if hi != 0xC0 || lo < 0x80 || lo > 0x8F {
		return false
	}

	// Bit 3 -> Bank 0/1 offset
	c.bank = (lo << 1) & 0x10

	lo &= 0x03
	c.romIN = lo == 0x01 || lo == 0x02
	c.ramRW = lo == 0x01 || lo == 0x03
	c.last = 0x00

	return true
}

// Reset resets the Language Card.
func (c *Card) Reset() {
	c.romIN = true
	c.ramRW = true
	c.last = 0x00
	c.bank = 0
}

// Slot is set by the memory Manager, depending on where this device was mounted.
func (*Card) Slot(byte) {}

// Memory returns card's memory.
func (c *Card) Memory() memory.Memory {
	return c.ram
}
