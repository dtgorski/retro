// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package diskette

type (
	// Drive is an Apple II disk drive, kind of.
	Drive struct {
		image *Image
		noise *TrackReader
		phase [0x02]byte
		motor bool
	}
)

// NewDrive creates an Apple II disk drive.
func NewDrive() *Drive {
	return &Drive{
		noise: NewTrackReader([]byte{0x44, 0x54, 0x47, 0x00}),
	}
}

// Insert loads new image into this drive.
func (d *Drive) Insert(image *Image) {
	d.image = image
}

// Eject removes the image from this drive.
func (d *Drive) Eject() {
	d.image = nil
}

// Motor turns the diskette drive motor on/off.
func (d *Drive) Motor(state bool) {
	d.phase[0] = 0xFF
	d.phase[1] = 0xFF
	d.motor = state
}

// Phase updates the state of a drive's stepper motor based on the current phase.
// It maintains the previous phase state and determines the direction of movement.
func (d *Drive) Phase(phase byte, state bool) {

	// Do not care about the motor state, but check:
	if d.image == nil || !state {
		return
	}

	// Translate the four phases in stepper movements. Track 0 resides
	// at the outermost disk location. This implementation is incomplete.
	// I hope, DOS 3.3 recalibration forgives me.
	P0, P1 := byte(0), byte(1)
	P2, P3 := byte(2), byte(3)

	// Mapping between current phase, last phase, and movement direction.
	phases := []struct {
		this byte
		last byte
		dir  byte
	}{
		{P0, P1, 0}, {P1, P2, 0}, {P2, P3, 0}, {P3, P0, 0}, // outward
		{P0, P3, 1}, {P1, P0, 1}, {P2, P1, 1}, {P3, P2, 1}, // inward
	}
	movement := [2]func(){
		d.image.halfTrackOut,
		d.image.halfTrackIn,
	}

	// Unshift current phase.
	d.phase[1] = d.phase[0]
	d.phase[0] = phase & 0x03

	// Determine the direction of head movement based on the current
	// and previous phases. Last repelling magnet determines direction.
	for _, p := range phases {
		if d.phase[0] == p.this && d.phase[1] == p.last {
			movement[p.dir]()
			break
		}
	}
}

// TrackReader returns the reader for the current half-/track.
func (d *Drive) TrackReader() *TrackReader {
	if d.image == nil {
		return d.noise
	}
	return d.image.trackReader()
}
