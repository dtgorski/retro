// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

import (
	"io"
)

type (
	// Image is a DOS 3.3 disk (.dsk) image.
	// Track 0 is at the outermost location.
	Image struct {
		tracks    [35 * 2]Track
		readers   [35 * 2]Reader
		encoder   *Encoder
		decoder   *Decoder
		halfTrack int
	}

	// Track represents a disk track with sectors.
	Track struct {
		sectors [16]sector
		track   int
	}

	sector [0x100]byte
)

// NewStandardImage creates image with a standard encoder and decoder.
func NewStandardImage() *Image {
	return NewImage(NewEncoder(), NewDecoder())
}

// NewImage creates a new disk image.
func NewImage(encoder *Encoder, decoder *Decoder) *Image {
	return &Image{
		encoder: encoder,
		decoder: decoder,
	}
}

// Load loads a disk image and prepares track readers for use with Drive.
func (im *Image) Load(r io.Reader) error {
	buf := make([]byte, 0x100)

	// Populate tracks, skip half-tracks.
	for t, n := 0, len(im.tracks); t < n; t += 2 {
		for s := 0; s < 0x10; s++ {
			num, err := r.Read(buf)

			if num == 0 || err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			copy(im.tracks[t].sectors[s][:], buf)
		}
	}

	// Each track and (empty) half-track has its own reader.
	for t, n := 0, len(im.tracks); t < n; t++ {
		if t&0x01 == 0 {
			im.tracks[t].track = t >> 1
		} else {
			im.tracks[t].track = 0xFF - t
		}
		im.readers[t] = NewTrackReader(
			im.encoder.Encode(&im.tracks[t]),
		)
	}
	return nil
}

// MustLoad panics if the disk image can not be loaded.
func (im *Image) MustLoad(r io.Reader) *Image {
	if err := im.Load(r); err != nil {
		panic(err)
	}
	return im
}

func (im *Image) halfTrackIn() {
	if im.halfTrack++; im.halfTrack >= len(im.tracks) {
		im.halfTrack = len(im.tracks) - 1
	}
}

func (im *Image) halfTrackOut() {
	if im.halfTrack--; im.halfTrack < 0 {
		im.halfTrack = 0
	}
}

func (im *Image) trackReader() *TrackReader {
	return im.readers[im.halfTrack].(*TrackReader)
}
