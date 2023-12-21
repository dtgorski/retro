// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package memory

type (
	// Device is a peripheral driver.
	Device interface {
		Read(lo, hi byte) (byte, bool)
		Write(lo, hi, b byte) bool
		Reset()
		Slot(num byte)
	}

	// Manager delegates memory access.
	Manager struct {
		mem  Memory
		dev  []Device
		list []Device
	}
)

// NewManager creates a new system memory manager unit.
func NewManager(mem Memory, devices ...Device) *Manager {
	m := &Manager{mem: mem}

	// 0-7 reserved for slotted devices.
	m.dev = make([]Device, 8)
	m.dev = append(m.dev, devices...)

	return m
}

// Mount sets device in slot.
func (m *Manager) Mount(slot byte, device Device) {
	device.Slot(slot & 0x07)
	device.Reset()

	m.dev[slot&0x07] = device

	// Prepare a list for a little faster access later on.
	m.list = m.list[:]
	for i := len(m.dev) - 1; i >= 0; i-- {
		if m.dev[i] != nil {
			m.list = append(m.list, m.dev[i])
		}
	}
}

// Slot returns device in slot or nil.
func (m *Manager) Slot(num byte) Device {
	return m.dev[num&0x07]
}

// Read reads a byte from address space.
func (m *Manager) Read(lo, hi byte) byte {
	for _, dev := range m.list {
		if b, ok := dev.Read(lo, hi); ok {
			return b
		}
	}
	return m.mem.Read(lo, hi)
}

// Writes a byte to address space.
func (m *Manager) Write(lo, hi, b byte) {
	for _, dev := range m.list {
		if dev.Write(lo, hi, b) {
			return
		}
	}
	if hi < 0xC1 {
		m.mem.Write(lo, hi, b)
	}
}

// Reset resets all devices.
func (m *Manager) Reset() {
	for _, dev := range m.list {
		dev.Reset()
	}
}
