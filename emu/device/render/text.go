// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package render

type (
	// Text is responsible for Text rendering.
	Text struct {
		page1 []byte
		page2 []byte
		temp  []byte
		font  *Font
	}
)

// NewText creates a new Text, responsible for rendering 40x24 frames.
func NewText(page1, page2 []byte, font *Font) *Text {
	return &Text{
		page1: page1,
		page2: page2,
		temp:  make([]byte, len(page1)),
		font:  font,
	}
}

// Render draws the content of selected memory page to the provided canvas byte array.
func (t *Text) Render(page byte, canvas []byte, flash bool) {
	if p := page & 0x01; p == 0 {
		copy(t.temp, t.page1)
	} else {
		copy(t.temp, t.page2)
	}
	for i := 0; i < 24; i++ {
		t.renderRow(i, canvas, flash)
	}
}

// Mixed draws the last four text lines (mixed mode) to the provided canvas byte array.
func (t *Text) Mixed(page byte, canvas []byte, flash bool) {
	if p := page & 0x01; p == 0 {
		copy(t.temp, t.page1)
	} else {
		copy(t.temp, t.page2)
	}
	for i := 20; i < 24; i++ {
		t.renderRow(i, canvas, flash)
	}
}

func (t *Text) renderRow(row int, canvas []byte, flash bool) {
	x := 0
	y := row * (Width << 5)
	p := (row>>3)*0x28 + (row&0x07)*0x80

	for i := 0; i < 0x28; i++ {
		b := t.temp[p+i]

		if flash && b&0x40 != 0 {
			b |= 0x80
		}

		g := t.font.glyph(b)

		for h, o := 0, 0; h < 8; h++ {
			copy(canvas[y+x+o:], g[h<<5:h<<5+28])
			o += Width << 2
		}
		x += 7 << 2
	}
}
