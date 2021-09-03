package pkg

import (
	"image/png"
	"log"
	"os"
)

type Texture struct {
	Width     uint32
	Height    uint32
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
	texture, err := os.Open(fileName)
	if err != nil {
		wd, _ := os.Getwd()
		log.Println("Working directory:",wd)
		log.Fatal(err)
	}
	defer func(texture *os.File) {_ = texture.Close() }(texture)

	texImg, err := png.Decode(texture)
	if err != nil {
		log.Fatal(err)
	}

	width := texImg.Bounds().Dx()
	height := texImg.Bounds().Dy()

	// Make a square array and fill from the image
	bmp := make([][]uint32, height)
	for i:=0;i<height;i++ {
		bmp[i] = make([]uint32, width)
		for j := 0; j < width; j++ {
			r,g,b,_ := texImg.At(j,i).RGBA()
			r = r>>8; g = g>>8; b = b>>8

			bmp[i][j] = r<<16 | g << 8 | b
		}
	}

	return Texture{
		Width:     uint32(width),
		Height:    uint32(height),
		Bmp:       bmp,
		ScanBytes: width * 4,
	}
}
