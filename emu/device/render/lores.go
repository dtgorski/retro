// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package render

type (
	// LoRes is responsible for Low Resolution (GR) rendering.
	LoRes struct {
		page1  []byte
		page2  []byte
		temp   []byte
		colors color16
	}

	color16 map[byte][]byte
)

// NewLoRes creates a new LoRes (low resolution) render object.
func NewLoRes(page1, page2 []byte, colors []int) *LoRes {

	colors16 := color16{}
	for i := byte(0); i < 16; i++ {
		colors16[i] = make([]byte, 7*4)

		for j := 0; j < 7; j++ {
			colors16[i][j<<2+0] = byte(colors[i] >> 24)
			colors16[i][j<<2+1] = byte(colors[i] >> 16)
			colors16[i][j<<2+2] = byte(colors[i] >> 8)
			colors16[i][j<<2+3] = byte(colors[i])
		}
	}

	return &LoRes{
		page1:  page1,
		page2:  page2,
		temp:   make([]byte, len(page1)),
		colors: colors16,
	}
}

// Render draws the content of selected memory page
// to the provided canvas byte array.
func (r *LoRes) Render(page byte, canvas []byte, _ bool) {
	if p := page & 0x01; p == 0 {
		copy(r.temp, r.page1)
	} else {
		copy(r.temp, r.page2)
	}
	for i := 0; i < 0x18; i++ {
		r.renderRow(i, canvas)
	}
}

func (r *LoRes) renderRow(row int, canvas []byte) {
	x := 0
	y := row * (Width * 8 * 4)
	p := (row>>3)*0x28 + (row&0x07)*0x80

	for i := 0; i < 0x28; i++ {
		top := r.colors[r.temp[p+i]&0x0F]
		btm := r.colors[r.temp[p+i]>>4]

		o := 0
		for h := 0; h < 4; h++ {
			copy(canvas[y+x+o:], top)
			o += Width * 4
		}
		for h := 0; h < 4; h++ {
			copy(canvas[y+x+o:], btm)
			o += Width * 4
		}
		x += 7 << 2
	}
}
