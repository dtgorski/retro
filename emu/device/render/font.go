// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package render

import (
	"image/png"
	"io"
)

type (
	// Font provides 256 character glyphs.
	Font struct {
		glyphs [256]glyph
		r      byte
		g      byte
		b      byte
		a      byte
	}

	// glyph is a 64 bytes of four byte RGBA representation.
	glyph [64 * 4]byte
)

// NewFont creates a new Font for a 40x24 character resolution.
func NewFont(r io.Reader, color int) *Font {
	return (&Font{
		r: byte(color >> 24),
		g: byte(color >> 16),
		b: byte(color >> 8),
		a: byte(color),
	}).populate(r)
}

// glyph returns the glyph for the given byte.
func (f *Font) glyph(b byte) glyph {
	return f.glyphs[b]
}

func (f *Font) populate(r io.Reader) *Font {
	font, _ := png.Decode(r)

	// Expects sprite image to be 128x128px, containing 8x8px glyphs.
	for y := 0; y < font.Bounds().Dy(); y += 8 {
		for x := 0; x < font.Bounds().Dx(); x += 8 {
			i, j := y<<1+x>>3, 0

			for h := 0; h < 8; h++ {
				for w := 0; w < 8; w++ {
					r, _, _, _ := font.At(x+w, y+h).RGBA()

					if r == 0 {
						f.glyphs[i][j+0] = 0x00
						f.glyphs[i][j+1] = 0x00
						f.glyphs[i][j+2] = 0x00
						f.glyphs[i][j+3] = 0xFF
					} else {
						f.glyphs[i][j+0] = f.r
						f.glyphs[i][j+1] = f.g
						f.glyphs[i][j+2] = f.b
						f.glyphs[i][j+3] = f.a
					}
					j += 4
				}
			}
		}
	}
	return f
}
