// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

type (
	// Decoder is track/sector decoder.
	Decoder struct{}
)

// NewDecoder return a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode is not implemented.
func (*Decoder) Decode(byte) *Track {
	panic("Decoder.Decode() not implemented")
}
