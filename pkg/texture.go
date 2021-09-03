package pkg

type Texture struct {
	Width     int
	Height    int
	Bmp       [][]uint32 // BGRA 8-bit-per-channel
	ScanBytes int
}

// NullTexture makes a new blank texture
func NullTexture() Texture {
	white := uint32(0x00ffFFff)
	light := uint32(0x00aaAAaa)
	dark_ := uint32(0x00776677)
	black := uint32(0x00000000)
	return Texture{
		Width:  4,
		Height: 4,
		Bmp: [][]uint32{
			{white, light, dark_, black},
			{light, dark_, black, dark_},
			{dark_, black, dark_, light},
			{black, dark_, light, white},
		},
		ScanBytes: 8,
	}
}

// LoadFromPng makes a texture from a PNG file
func LoadFromPng(fileName string) Texture{
	// TODO: implement
	return Texture{}
}
