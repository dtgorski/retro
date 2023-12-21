// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package virtual

import (
	"retro/emu/input"
)

type (
	// Channels transports peripheral I/O event channels.
	Channels struct {
		keyCh chan input.KeyInput
		bufCh chan input.KeyInput
		butCh chan input.MouseButton
		posCh chan input.CursorPos
	}
)

// NewChannels create a new I/O channel transport.
func NewChannels(
	keyCh chan input.KeyInput,
	bufCh chan input.KeyInput,
	butCh chan input.MouseButton,
	posCh chan input.CursorPos,
) *Channels {
	return &Channels{keyCh, bufCh, butCh, posCh}
}

// KeyInput returns the key press/release event channel.
func (c *Channels) KeyInput() chan input.KeyInput {
	return c.keyCh
}

// KeyBuffer returns the key input buffer.
func (c *Channels) KeyBuffer() chan input.KeyInput {
	return c.bufCh
}

// MouseButton returns the mouse button channel.
func (c *Channels) MouseButton() chan input.MouseButton {
	return c.butCh
}

// CursorPos return the  mouse cursor position channel.
func (c *Channels) CursorPos() chan input.CursorPos {
	return c.posCh
}
