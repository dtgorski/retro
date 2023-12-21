// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

type (
	// Card is an Apple Disk II Interface Card.
	Card struct {
		rom  []byte
		drv1 *Drive
		drv2 *Drive
		hot  *Drive
		slot byte
	}
)

// NewCard creates a new Apple Disk II card with two diskette drives.
func NewCard(rom []byte) *Card {

	patched := make([]byte, len(rom))
	copy(patched, rom)

	// Skip wait routine.
	patched[0x4C] = 0xA9 // LDA
	patched[0x4D] = 0x00 // #0
	patched[0x4E] = 0xEA // NOP

	drv1 := NewDrive()
	drv2 := NewDrive()

	card := &Card{
		rom:  patched,
		drv1: drv1,
		drv2: drv2,
		hot:  drv1,
	}
	card.Reset()

	return card
}

func (c *Card) switches() map[byte]func(b byte) byte {
	return map[byte]func(b byte) byte{
		0x00: func(b byte) byte { c.hot.Phase(0>>1, false); return 0 }, // DRV_P0_OFF
		0x01: func(b byte) byte { c.hot.Phase(1>>1, true); return 0 },  // DRV_P0_ON
		0x02: func(b byte) byte { c.hot.Phase(2>>1, false); return 0 }, // DRV_P1_OFF
		0x03: func(b byte) byte { c.hot.Phase(3>>1, true); return 0 },  // DRV_P1_ON
		0x04: func(b byte) byte { c.hot.Phase(4>>1, false); return 0 }, // DRV_P2_OFF
		0x05: func(b byte) byte { c.hot.Phase(5>>1, true); return 0 },  // DRV_P2_ON
		0x06: func(b byte) byte { c.hot.Phase(6>>1, false); return 0 }, // DRV_P3_OFF
		0x07: func(b byte) byte { c.hot.Phase(7>>1, true); return 0 },  // DRV_P3_ON
		0x08: func(b byte) byte { c.hot.Motor(false); return 0 },       // DRV_OFF
		0x09: func(b byte) byte { c.hot.Motor(true); return 0 },        // DRV_ON
		0x0A: func(b byte) byte { c.hot = c.drv1; return 0 },           // DRV_SEL1
		0x0B: func(b byte) byte { c.hot = c.drv2; return 0 },           // DRV_SEL2
		0x0C: func(b byte) byte { return c.hot.TrackReader().Read() },  // DRV_SHIFT / Q6L
		0x0D: func(b byte) byte { return 0 },                           // DRV_LOAD / Q6H
		0x0E: func(b byte) byte { return 0 },                           // DRV_READ / Q7L
		0x0F: func(b byte) byte { return 0 },                           // DRV_WRITE / Q7H
	}
}

// Read reads a byte, if this device is sensitive to this address.
func (c *Card) Read(lo, hi byte) (byte, bool) {

	// Not interested?
	if c.slot == 0 || c.slot > 7 {
		return 0, false
	}
	// Read Card ROM? 0xCn00-0xCnFF?
	if hi == 0xC0|c.slot {
		return c.rom[lo], true
	}
	// I/O switches?
	if hi == 0xC0 && lo >= 0x80|(c.slot<<4) && lo <= 0x8F|(c.slot<<4) {
		return c.switches()[lo&0x0F](0), true
	}
	return 0, false
}

// Write writes a byte, if this device is sensitive to this address.
func (c *Card) Write(lo, hi, _ byte) bool {

	// Not interested?
	if c.slot == 0 || c.slot > 7 {
		return false
	}
	// Write Card ROM?! 0xCn00-0xCnFF?
	if hi == 0xC0|c.slot {
		return true
	}
	// I/O switches?
	if hi == 0xC0 && lo >= 0x80|(c.slot<<4) && lo <= 0x8F|(c.slot<<4) {
		c.switches()[lo&0x0F](0)
		return true
	}
	return false
}

// Reset resets the Disk Card.
func (c *Card) Reset() {
	c.drv1.Motor(false)
	c.drv2.Motor(false)
}

// Slot is set by the memory Manager, depending on where this device was mounted.
func (c *Card) Slot(num byte) {
	c.slot = num & 0x07
}

// DMA allows to directly access memory.
func (c *Card) DMA() []byte {
	return c.rom
}

// Drive returns the drive by number (0/1).
func (c *Card) Drive(num int) *Drive {
	if num&0x01 == 0x00 {
		return c.drv1
	}
	return c.drv2
}
