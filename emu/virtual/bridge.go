// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package virtual

import (
	"retro/emu/device/render"
	"retro/emu/input"
	"retro/emu/memory"
)

type (
	// Bridge brings memory, cards and graphics rendering together.
	Bridge struct {
		manager  *memory.Manager
		driver   *render.Driver
		keyMap   *input.KeyMap
		channels *Channels
	}
)

// NewBridge creates a new Bridge.
func NewBridge(
	manager *memory.Manager,
	driver *render.Driver,
	keyMap *input.KeyMap,
	channels *Channels,
) *Bridge {
	return &Bridge{manager, driver, keyMap, channels}
}

// Memory is system memory manager unit.
func (b *Bridge) Memory() *memory.Manager {
	return b.manager
}

// Renderer returns the display rendering driver.
func (b *Bridge) Renderer() *render.Driver {
	return b.driver
}

// Reset resets all peripheral cards.
func (b *Bridge) Reset() {
	b.manager.Reset()
}

// KeyMap return the key translation map.
func (b *Bridge) KeyMap() *input.KeyMap {
	return b.keyMap
}

// Channels returns the I/O channels transport.
func (b *Bridge) Channels() *Channels {
	return b.channels
}
