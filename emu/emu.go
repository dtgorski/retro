// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package emu

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"retro/emu/config"
	"retro/emu/input"
	"retro/emu/virtual"
	"retro/gui"
	"strings"
	"syscall"
)

// Run runs an Apple II emulation.
func Run(ctx context.Context, conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				err = errors.New(t.Error())
			default:
				err = errors.New(t.(string))
			}
		}
	}()

	// POSIX signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// ---

	// Keyboard buffer, window I/O, input translation.
	channels := virtual.NewChannels(
		make(chan input.KeyInput),
		make(chan input.KeyInput, 0x1000),
		make(chan input.MouseButton),
		make(chan input.CursorPos),
	)
	keyMap := input.NewKeyMap()

	// The emulator.
	machine := virtual.NewAppleTwo(conf, keyMap, channels)
	bridge := machine.Bridge()

	if err = insertDisks(conf, bridge); err != nil {
		return err
	}

	// Emulator power on.
	errCh := make(chan error)
	go func() { errCh <- machine.PowerOn(ctx) }()

	// ---

	// New output window.
	props := gui.Properties{
		Width:  conf.Window.Zoom * 280,
		Height: conf.Window.Zoom * 192,
		Title:  conf.Window.Title,
	}
	win := gui.NewWindow(props, bridge.Renderer(), channels)

	if err = win.Open(); err != nil {
		return err
	}
	go func() { errCh <- win.RenderAndListen() }()

	// ---

	// Helper, to be run concurrently on demand.
	paste := func(s string) {
		for _, b := range []byte(strings.ToUpper(s)) {
			channels.KeyBuffer() <- keyMap.FromASCII(b)
		}
	}

	// Paddle positioning.
	aspectW := 255 / float64(props.Width)
	aspectH := 255 / float64(props.Height)

	mem := bridge.Memory()

	// Main loop.
	for {
		select {

		// Keyboard in.
		case key := <-channels.KeyInput():
			switch {
			case key.IsCtrlShiftR():
				machine.Reset()
			case key.IsCtrlV():
				go paste(win.Clipboard())
			default:
				channels.KeyBuffer() <- key
			}

		// Paddle move, requires tweaked ROM at 0xFB28 (see garage.go).
		case pos := <-channels.CursorPos():
			x := byte(pos.X() * aspectW)
			y := byte(pos.Y() * aspectH)
			mem.Write(0x64, 0xC0, x)
			mem.Write(0x65, 0xC0, y)

		// Mouse button.
		case but := <-channels.MouseButton():
			no := byte(but.Button())
			on := map[bool]byte{true: 0x80, false: 0x00}
			mem.Write(0x61+no, 0xC0, on[but.IsPressed()])

		// Machine or window error.
		case err = <-errCh:
			if err != nil && err.Error() != "context canceled" {
				return fmt.Errorf("internal fault: %s", err)
			}
			return nil

		// POSIX signal.
		case <-sigCh:
			return fmt.Errorf("interrupted")
		}
	}
}
