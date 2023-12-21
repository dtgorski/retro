// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package virtual

/*
#include <time.h>
*/
import "C"
import (
	"context"
	"errors"
	"time"
)

type (
	// The CPU interface abstracts the processor implementation.
	CPU interface {
		PC(lo, hi byte)
		PCL() byte
		PCH() byte
		Reset()
		Step() (cycles uint, err error)
	}

	// Machine represents the Apple II computer itself.
	Machine struct {
		bridge *Bridge
		cpu    CPU
		clock  int
	}
)

// NewMachine creates a new Machine.
func NewMachine(bridge *Bridge, cpu CPU, hz int) *Machine {
	return &Machine{bridge, cpu, hz}
}

// Bridge returns the Machine's Bridge.
func (m *Machine) Bridge() *Bridge {
	return m.bridge
}

// PowerOn runs the machine.
func (m *Machine) PowerOn(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	// Consume keyboard buffer, check for strobe (KBDSTRB).
	go func() {
		var key byte
		memory := m.Bridge().Memory()
		keyMap := m.Bridge().KeyMap()
		keyBuf := m.Bridge().Channels().KeyBuffer()

		for ctx.Err() == nil {
			if key = keyMap.FromInput(<-keyBuf); key == 0x00 {
				continue
			}
			for memory.Read(0x00, 0xC0)&0x80 != 0x00 {
				time.Sleep(time.Millisecond * 1)
			}
			if ctx.Err() == nil {
				memory.Write(0x00, 0xC0, key)
			}
		}
	}()

	// ---

	batch := 1000
	cycles := uint(0)

loop:
	if ctx.Err() != nil {
		return ctx.Err()
	}
	now := time.Now()

	// Running each CPU step subsequently in "real time" will require to take
	// a sleeping break for a few hundred nanoseconds after each step.
	// The overhead associated with the sleep request is not proportional
	// with the gain (high CPU Load). Let's try in batches, so we can sleep
	// longer less often and the real CPU Load is reduced. Timing will break.
	for i := m.clock / batch; i > 0; {

		// Skips dynamically loaded DOS 3.3 wait routine.
		if m.cpu.PCH() == 0xBA && m.cpu.PCL() == 0x00 {
			m.cpu.PC(0x10, 0xBA)
		}
		if cycles, err = m.cpu.Step(); err != nil {
			return err
		}
		i -= int(cycles)
	}

	dur := time.Since(now).Nanoseconds()

	// Rough approximation.
	if delay := int(870*int64(batch) - dur); delay > 100_000 {
		retard(delay)
	}
	goto loop
}

// Reset resets the bridge (peripheral cards and CPU).
func (m *Machine) Reset() {
	m.bridge.Reset()
	m.cpu.Reset()
}

// Seems to work better than time.Sleep().
func retard(ns int) {
	var ts C.struct_timespec
	ts.tv_sec = C.time_t(0)
	ts.tv_nsec = C.long(ns)
	_, _ = C.nanosleep(&ts, nil)
}
