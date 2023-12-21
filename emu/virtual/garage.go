// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package virtual

import (
	cpu "github.com/dtgorski/m6502"
	"retro/emu/config"
	"retro/emu/device/builtin"
	"retro/emu/device/diskette"
	"retro/emu/device/language"
	"retro/emu/device/render"
	"retro/emu/files"
	"retro/emu/input"
	"retro/emu/memory"
)

// NewAppleTwo creates an Apple II setup.
func NewAppleTwo(conf *config.Config, keyMap *input.KeyMap, channels *Channels) *Machine {
	hz := int(conf.MHz * 1024 * 1024)

	// Main 64KB memory segment.
	mem := memory.NewMemory()

	// Display driver and I/O page (soft switches)
	renderer := render.NewDriver(createRenderModes(conf, mem.DMA()))
	keyboard := builtin.NewKeyboard(mem)
	paddle := builtin.NewPaddle(mem)

	// Delegates reads/writes to devices (I/O page, slots).
	mmu := memory.NewManager(mem, renderer, keyboard, paddle)

	// Onboard ROM, load Applesoft Basic and Monitor.
	mem.MustLoad(0xF800, files.MustOpen(files.ROM_APPLESOFT_BASIC_MON_F800))
	mem.MustLoad(0xF000, files.MustOpen(files.ROM_APPLESOFT_BASIC_F000))
	mem.MustLoad(0xE800, files.MustOpen(files.ROM_APPLESOFT_BASIC_E800))
	mem.MustLoad(0xE000, files.MustOpen(files.ROM_APPLESOFT_BASIC_E000))
	mem.MustLoad(0xD800, files.MustOpen(files.ROM_APPLESOFT_BASIC_D800))
	mem.MustLoad(0xD000, files.MustOpen(files.ROM_APPLESOFT_BASIC_D000))

	// Settings for Applesoft Basic in Zero Page.
	copy(mem.DMA()[0x67:], []byte{
		0x01, 0x08, // $67 - $68 | Start of program address
		0x03, 0x08, // $69 - $6A | Start of variable space address
		0x03, 0x08, // $6B - $6C | Start of array space address
		0x03, 0x08, // $6D - $6E | End of numeric storage address
	})
	copy(mem.DMA()[0x73:], []byte{
		0x00, 0x96, // $73 - $74 | HIMEM address
	})
	copy(mem.DMA()[0xAF:], []byte{
		0x03, 0x08, // $AF - $B0 | Pointer to end of Applesoft program
	})

	// Tweak ROM, modify paddle wait routine.
	copy(mem.DMA()[0xFB28:], []byte{
		0xA8, 0x60, // TAY, RTS
	})

	// Slot #0, mount Language Card.
	mmu.Mount(0, language.NewCard())

	// Slot #6, mount interface ROM only when disk image(s) provided.
	if len(conf.Disk.Drive1) > 0 || len(conf.Disk.Drive2) > 0 {
		card := diskette.NewCard(files.MustLoad(files.ROM_APPLE_DISK_II_16))
		mmu.Mount(6, card)
	}

	return NewMachine(NewBridge(mmu, renderer, keyMap, channels), cpu.New(mmu), hz)
}

// createRenderModes creates rendering modes.
func createRenderModes(conf *config.Config, mem []byte) render.Modes {

	// Text/graphics render modes.
	return render.Modes{
		render.ModeText: render.NewText(
			mem[0x0400:0x0800], // "page" 1
			mem[0x0800:0x0C00], // "page" 2
			render.NewFont(
				files.MustOpen(files.FONT_APPLE_II),
				conf.Render.Mono.Color,
			),
		),
		render.ModeLoRes: render.NewLoRes(
			mem[0x0400:0x0800], // "page" 1
			mem[0x0800:0x0C00], // "page" 2
			conf.Render.LoRes.Colors[:],
		),
		render.ModeHiRes: render.NewHiRes(
			mem[0x2000:0x4000], // "page" 1
			mem[0x4000:0x6000], // "page" 2
			conf.Render.HiRes.Colors[:],
		),
	}
}
