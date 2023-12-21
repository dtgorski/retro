// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

import (
	"bytes"
)

type (
	// Encoder is track/sector encoder.
	Encoder struct{}
)

var (
	addrPrologue = []byte{0xD5, 0xAA, 0x96}
	addrEpilogue = []byte{0xDE, 0xAA, 0xEB}
	dataPrologue = []byte{0xD5, 0xAA, 0xAD}
	dataEpilogue = []byte{0xDE, 0xAA, 0xEB}

	sixAndTwo = []byte{
		0x96, 0x97, 0x9A, 0x9B, 0x9D, 0x9E, 0x9F, 0xA6,
		0xA7, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB2, 0xB3,
		0xB4, 0xB5, 0xB6, 0xB7, 0xB9, 0xBA, 0xBB, 0xBC,
		0xBD, 0xBE, 0xBF, 0xCB, 0xCD, 0xCE, 0xCF, 0xD3,
		0xD6, 0xD7, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDE,
		0xDF, 0xE5, 0xE6, 0xE7, 0xE9, 0xEA, 0xEB, 0xEC,
		0xED, 0xEE, 0xEF, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6,
		0xF7, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF,
	}
)

// NewEncoder creates a new track/sector encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode translates the pure track data into a byte stream,
// that is suitable for reading by the Apple Disk II ROM.
func (e *Encoder) Encode(track *Track) []byte {
	vol := byte(0xFE)
	buf := bytes.Buffer{}

	var sec byte
	for num := range track.sectors {
		if sec = byte(num); num != 15 {
			sec = byte((num * 7) % 15)
		}
		_, _ = buf.Write(addrPrologue)
		_, _ = buf.Write(e.fourAndFour(vol))
		_, _ = buf.Write(e.fourAndFour(byte(track.track)))
		_, _ = buf.Write(e.fourAndFour(byte(num)))
		_, _ = buf.Write(e.fourAndFour(vol ^ byte(track.track) ^ byte(num)))
		_, _ = buf.Write(addrEpilogue)
		_, _ = buf.Write(dataPrologue)
		_, _ = buf.Write(e.sixAndTwo(track.sectors[sec][:]))
		_, _ = buf.Write(dataEpilogue)
	}
	return buf.Bytes()
}

// Borrowed from https://github.com/TomHarte/dsk2woz (MIT License)
func (*Encoder) sixAndTwo(b []byte) []byte {
	buf := [0x56 + 0x100 + 0x01]byte{}
	bit := []byte{0, 2, 1, 3}

	// Fill in byte values: the first 86 bytes contain shuffled and combined
	// copies of the bottom two bits of the sector contents.
	for i := 0; i < 0x54; i++ {
		buf[i] = bit[b[i]&0x03] | bit[b[i+0x56]&0x03]<<2 | bit[b[i+0xAC]&0x03]<<4
	}
	buf[0x54] = bit[b[0x54]&0x03] | bit[b[0xAA]&0x03]<<2
	buf[0x55] = bit[b[0x55]&0x03] | bit[b[0xAB]&0x03]<<2

	// These 256 bytes are the remaining six bits.
	for i := 0; i < 0x100; i++ {
		buf[i+0x56] = b[i] >> 2
	}

	// Exclusive OR each byte with the one before it.
	buf[0x156] = buf[0x155]
	for pos := 0x156; pos > 1; {
		pos--
		buf[pos] ^= buf[pos-1]
	}

	// Map six-bit values up to full bytes.
	for i := 0; i < 0x157; i++ {
		buf[i] = sixAndTwo[buf[i]]
	}
	return buf[:]
}

func (*Encoder) fourAndFour(b byte) []byte {
	return []byte{(b >> 1) | 0xAA, b | 0xAA}
}
