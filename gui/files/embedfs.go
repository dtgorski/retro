// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package files

import (
	"embed"
	"io"
)

//go:embed "retro_16x16.png"
//go:embed "retro_24x24.png"
//go:embed "retro_32x32.png"
//go:embed "retro_48x48.png"
//go:embed "default.vert"
//go:embed "default.frag"
var fs embed.FS

// GUI resources.
const (
	WINDOW_ICON_16x = "retro_16x16.png"
	WINDOW_ICON_24x = "retro_24x24.png"
	WINDOW_ICON_32x = "retro_32x32.png"
	WINDOW_ICON_48x = "retro_48x48.png"

	SHADER_DEFAULT_VERT = "default.vert"
	SHADER_DEFAULT_FRAG = "default.frag"
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
