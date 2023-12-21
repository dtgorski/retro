// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

type (
	// Reader is an endless stream reader.
	Reader interface {
		Read() byte
	}

	// TrackReader provides track data.
	TrackReader struct {
		buf []byte
		pos int
	}
)

// NewTrackReader creates a new track reader.
func NewTrackReader(buf []byte) *TrackReader {
	return &TrackReader{buf, 0}
}

// Reads a byte from the track stream. When the track reader has reached
// the end of its buffer, it starts again from the beginning, simulating
// an endless stream of bits (here as bytes) from a diskette.
func (r *TrackReader) Read() byte {
	b := r.buf[r.pos]
	if r.pos++; r.pos == len(r.buf) {
		r.pos = 0
	}
	return b
}
