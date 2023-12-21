// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package render

type (
	// HiRes is responsible for High Resolution (HGR) rendering.
	HiRes struct {
		page1  []byte
		page2  []byte
		temp   []byte
		colors color8
	}

	color8 map[byte][]byte
)

// NewHiRes creates a new HiRes (high resolution) render object.
func NewHiRes(page1, page2 []byte, colors []int) *HiRes {

	colors8 := color8{}
	for i := byte(0); i < 8; i++ {
		r, g := byte(colors[i]>>24), byte(colors[i]>>16)
		b, a := byte(colors[i]>>8), byte(colors[i])
		dim := byte(float32(a) * 0.95)
		colors8[i] = []byte{r, g, b, dim, r, g, b, a}
	}

	return &HiRes{
		page1:  page1,
		page2:  page2,
		temp:   make([]byte, len(page1)),
		colors: colors8,
	}
}

// Render draws the content of selected memory page
// to the provided canvas byte array.
func (r *HiRes) Render(page byte, canvas []byte, _ bool) {
	if p := page & 0x01; p == 0 {
		copy(r.temp, r.page1)
	} else {
		copy(r.temp, r.page2)
	}
	for i := 0; i < 0x40; i++ {
		r.renderLine(i, canvas)
	}
}

func (r *HiRes) renderLine(line int, canvas []byte) {
	y := line * (Width * 4)
	p := (line>>3)*0x80 + (line&0x07)*0x400

	for i := 0; i < 3; i++ {
		x := 0
		for j := 0; j < 0x28; j += 2 {
			b0, b1 := r.temp[p+j], r.temp[p+j+1]

			c0 := (b0 >> 5) & 0x04
			c1 := (b1 >> 5) & 0x04

			// Extracting color information from the bytes
			// and mapping them to the corresponding colors.
			px := []byte{
				c0 | ((b0 >> 0) & 0x03),
				c0 | ((b0 >> 2) & 0x03),
				c0 | ((b0 >> 4) & 0x03),
				c0 | ((b0>>6)&0x01 | (b1&0x01)<<1),
				c1 | ((b1 >> 1) & 0x03),
				c1 | ((b1 >> 3) & 0x03),
				c1 | ((b1 >> 5) & 0x03),
			}
			for k := 0; k < 7; k++ {
				copy(canvas[y+x:], r.colors[px[k]])
				x += 8
			}
		}
		y += Width * 4 * 0x40
		p += 0x28
	}
}
