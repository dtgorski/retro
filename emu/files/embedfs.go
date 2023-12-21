// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package files

import (
	"embed"
	"io"
)

//go:embed "APPLESOFT_BASIC_D000.bin"
//go:embed "APPLESOFT_BASIC_D800.bin"
//go:embed "APPLESOFT_BASIC_E000.bin"
//go:embed "APPLESOFT_BASIC_E800.bin"
//go:embed "APPLESOFT_BASIC_F000.bin"
//go:embed "APPLESOFT_BASIC_AUTOSTART_MONITOR_F800.bin"
//go:embed "APPLE_DISK_II_16_SECTOR_INTERFACE_CARD_ROM_P5_(BOOT).bin"
//go:embed "font.png"
var fs embed.FS

// ROMs and stuff.
const (
	FONT_APPLE_II = "font.png"

	ROM_APPLESOFT_BASIC_D000     = "APPLESOFT_BASIC_D000.bin"
	ROM_APPLESOFT_BASIC_D800     = "APPLESOFT_BASIC_D800.bin"
	ROM_APPLESOFT_BASIC_E000     = "APPLESOFT_BASIC_E000.bin"
	ROM_APPLESOFT_BASIC_E800     = "APPLESOFT_BASIC_E800.bin"
	ROM_APPLESOFT_BASIC_F000     = "APPLESOFT_BASIC_F000.bin"
	ROM_APPLESOFT_BASIC_MON_F800 = "APPLESOFT_BASIC_AUTOSTART_MONITOR_F800.bin"

	ROM_APPLE_DISK_II_16 = "APPLE_DISK_II_16_SECTOR_INTERFACE_CARD_ROM_P5_(BOOT).bin"
)

// MustOpen panics if it can not open embedded file.
func MustOpen(asset string) io.Reader {
	file, err := fs.Open(asset)
	if err != nil {
		panic(err)
	}
	return file
}

// MustLoad panics if it can not load embedded file.
func MustLoad(asset string) []byte {
	b, err := io.ReadAll(MustOpen(asset))
	if err != nil {
		panic(err)
	}
	return b
}
