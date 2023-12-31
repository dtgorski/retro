// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package input

type (
	// KeyMap is responsible for keyboard input translation.
	KeyMap struct{}

	keyCode [8]byte
	keyMap  map[int]keyCode
)

// NewKeyMap creates a key input translation map.
func NewKeyMap() *KeyMap {
	return &KeyMap{}
}

// FromInput takes an input stroke, translates it to an Apple II character.
func (KeyMap) FromInput(e KeyInput) byte {
	if e.act == 1 || e.act == 2 {
		if code, ok := botched[e.key]; ok {
			off := e.mod & 0x03
			return code[off]
		}
	}
	return 0x00
}

// FromASCII takes an ASCII input (paste), translates it to an Apple II character.
func (KeyMap) FromASCII(b byte) KeyInput {
	for key, mapping := range botched {
		if mapping[4] != 0x00 && b == mapping[4] {
			return KeyInput{key: key, act: 1, mod: 0}
		}
		if mapping[5] != 0x00 && b == mapping[5] {
			return KeyInput{key: key, act: 1, mod: 1}
		}
	}
	return KeyInput{}
}

var botched = keyMap{
	//       KEY   SHIFT CTRL  BOTH     Paste: ASCII SHIFT
	0x0020: {0xA0, 0xA0, 0xA0, 0xA0 /*     */, 0x20, 0x20},

	0x002C: {0xAC, 0xBB, 0xAC, 0xBB /* , ; */, 0x2C, 0x3B},
	0x002D: {0x00, 0xBF, 0x00, 0xBF /*  ?  */, 0x00, 0x3F},
	0x002E: {0xAE, 0xBA, 0xAE, 0xBA /* . : */, 0x2E, 0x3A},
	0x002F: {0xAD, 0x00, 0xAD, 0x00 /*  -  */, 0x2D, 0x00},

	0x0030: {0xB0, 0xBD, 0xB0, 0xBD /* 0 = */, 0x30, 0x3D},
	0x0031: {0xB1, 0xA1, 0xB1, 0xA1 /* 1 ! */, 0x31, 0x21},
	0x0032: {0xB2, 0xA2, 0xB2, 0xA2 /* 2 " */, 0x32, 0x22},
	0x0033: {0xB3, 0xA3, 0xB3, 0xA3 /* 3 # */, 0x33, 0x23},
	0x0034: {0xB4, 0xA4, 0xB4, 0xA4 /* 4 $ */, 0x34, 0x24},
	0x0035: {0xB5, 0xA5, 0xB5, 0xA5 /* 5 % */, 0x35, 0x25},
	0x0036: {0xB6, 0xA6, 0xB6, 0xA6 /* 6 & */, 0x36, 0x26},
	0x0037: {0xB7, 0xAF, 0xB7, 0xAF /* 7 / */, 0x37, 0x2F},
	0x0038: {0xB8, 0xA8, 0xB8, 0xA8 /* 8 ( */, 0x38, 0x28},
	0x0039: {0xB9, 0xA9, 0xB9, 0xA9 /* 9 ) */, 0x39, 0x29},

	0x0041: {0xC1, 0xC1, 0x81, 0x81 /*  A  */, 0x41, 0x41},
	0x0042: {0xC2, 0xC2, 0x82, 0x82 /*  B  */, 0x42, 0x42},
	0x0043: {0xC3, 0xC3, 0x83, 0x83 /*  C  */, 0x43, 0x43},
	0x0044: {0xC4, 0xC4, 0x84, 0x84 /*  D  */, 0x44, 0x44},
	0x0045: {0xC5, 0xC5, 0x85, 0x85 /*  E  */, 0x45, 0x45},
	0x0046: {0xC6, 0xC6, 0x86, 0x86 /*  F  */, 0x46, 0x46},
	0x0047: {0xC7, 0xC7, 0x87, 0x87 /*  G  */, 0x47, 0x47},
	0x0048: {0xC8, 0xC8, 0x88, 0x88 /*  H  */, 0x48, 0x48},
	0x0049: {0xC9, 0xC9, 0x89, 0x89 /*  I  */, 0x49, 0x49},
	0x004A: {0xCA, 0xCA, 0x8A, 0x8A /*  J  */, 0x4A, 0x4A},
	0x004B: {0xCB, 0xCB, 0x8B, 0x8B /*  K  */, 0x4B, 0x4B},
	0x004C: {0xCC, 0xCC, 0x8C, 0x8C /*  L  */, 0x4C, 0x4C},
	0x004D: {0xCD, 0xDD, 0x8D, 0x9D /* M ] */, 0x4D, 0x5D},
	0x004E: {0xCE, 0xDE, 0x8E, 0x9E /* N ^ */, 0x4E, 0x5E},
	0x004F: {0xCF, 0xCF, 0x8F, 0x8F /*  O  */, 0x4F, 0x4F},
	0x0050: {0xD0, 0xC0, 0x90, 0x80 /* P @ */, 0x50, 0x40},
	0x0051: {0xD1, 0xD1, 0x91, 0x91 /*  Q  */, 0x51, 0x51},
	0x0052: {0xD2, 0xD2, 0x92, 0x00 /*  R  */, 0x52, 0x52},
	0x0053: {0xD3, 0xD3, 0x93, 0x93 /*  S  */, 0x53, 0x53},
	0x0054: {0xD4, 0xD4, 0x94, 0x94 /*  T  */, 0x54, 0x54},
	0x0055: {0xD5, 0xD5, 0x95, 0x95 /*  U  */, 0x55, 0x55},
	0x0056: {0xD6, 0xD6, 0x00, 0x96 /*  V  */, 0x56, 0x56},
	0x0057: {0xD7, 0xD7, 0x97, 0x97 /*  W  */, 0x57, 0x57},
	0x0058: {0xD8, 0xD8, 0x98, 0x98 /*  X  */, 0x58, 0x58},
	0x0059: {0xDA, 0xDA, 0x9A, 0x9A /*  Y  */, 0x5A, 0x5A},
	0x005A: {0xD9, 0xD9, 0x99, 0x99 /*  Z  */, 0x59, 0x59},

	0x005C: {0xA3, 0xA7, 0xA3, 0xA7 /* # ' */, 0x23, 0x27},
	0x005D: {0xAB, 0xAA, 0xAB, 0xAA /* + * */, 0x2B, 0x2A},
	0x00A1: {0xBC, 0xBE, 0xBC, 0xBE /* < > */, 0x3C, 0x3E},
	0x0100: {0x9B, 0x9B, 0x9B, 0x9B /* Esc */, 0x00, 0x00},
	0x0101: {0x8D, 0x8D, 0x8D, 0x8D /* Ret */, 0x0A, 0x0A},
	0x0103: {0x88, 0x88, 0x88, 0x88 /* BS  */, 0x00, 0x00},
	0x0106: {0x95, 0x95, 0x95, 0x95 /* ->  */, 0x00, 0x00},
	0x0107: {0x88, 0x88, 0x88, 0x88 /* <-  */, 0x00, 0x00},
}
