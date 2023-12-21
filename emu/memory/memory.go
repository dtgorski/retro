// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package memory

import (
	"io"
)

type (
	// Bus is a 8-bit data bus with a 16-bit little-endian address width.
	Bus interface {
		Read(lo, hi byte) byte
		Write(lo, hi, b byte)
	}

	// DMA allows to retrieve raw memory.
	DMA interface {
		DMA() []byte
	}

	// Memory is a 64KB, 16-bit CPU addressable portion of memory.
	Memory interface {
		Bus
		DMA
		Load(addr uint16, r io.Reader) error
		MustLoad(addr uint16, r io.Reader)
	}

	memory struct {
		bytes [0x10000]byte
	}
)

// NewMemory creates a new memory.
func NewMemory() Memory {
	return &memory{}
}

// Read reads a byte from address space.
func (m *memory) Read(lo, hi byte) byte {
	return m.bytes[uint16(hi)<<8|uint16(lo)]
}

// Writes a byte to address space.
func (m *memory) Write(lo, hi, b byte) {
	m.bytes[uint16(hi)<<8|uint16(lo)] = b
}

// DMA allows to directly access memory.
func (m *memory) DMA() []byte {
	return m.bytes[:]
}

// Load loads payload into the Memory.
func (m *memory) Load(addr uint16, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	copy(m.bytes[addr:], data)
	return nil
}

// MustLoad panics if the payload can not be loaded.
func (m *memory) MustLoad(addr uint16, r io.Reader) {
	if err := m.Load(addr, r); err != nil {
		panic(err)
	}
}
