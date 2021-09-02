package pkg

type Texture struct {
	Width     int
	Height    int
	Bmp       [][]uint32 // BGRA 8-bit-per-channel
	ScanBytes int
}

// NullTexture makes a new blank texture
func NullTexture() Texture {
	return Texture{
		Width:  2,
		Height: 2,
		Bmp: [][]uint32{
			{0,          0x808080FF},
			{0x808080FF, 0xffFFffFF},
		},
		ScanBytes: 8,
	}
}
