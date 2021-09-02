package pkg

import "unsafe"

type Renderer struct {
	width  int
	height int
}

// for some reason, if these are inside the struct, they won't draw in Win32
var (
	frameA  RenderFrame
	frameB  RenderFrame
	targetA bool // if true, draw into frameA, else draw into frameB
)

type RenderFrame struct {
	Width      int
	Height     int
	bmp        []byte
	scanBytes  int
	frameBytes uint64
}

func NewRenderer(width, height int) *Renderer {
	frameA = makeFrame(width, height)
	frameB = makeFrame(width, height)
	targetA = true
	return &Renderer{
		width:   width,
		height:  height,
	}
}

func makeFrame(width, height int) RenderFrame {
	scanBytes := width * 4
	bmpSize := uint64(scanBytes) * uint64(height) // 32 bit argb
	bmp := make([]byte, bmpSize)
	return RenderFrame{
		Width:  width,
		Height: height,
		bmp:    bmp,
		scanBytes: scanBytes,
		frameBytes: bmpSize,
	}
}

var frameCount int
// Update should update the world state and draw into the target buffer.
// You should switch buffers when you're finished.
func (r *Renderer) Update(t int64) {
	defer r.SwitchBuffer()

	// any world logic can go in here
	frame := r.TargetFrame()

	//buf := frame.bmp
	buf := frame.wordArray()
	frameCount++

	var y,x int
	var i int
	var v uint32

	var white uint32 = 0xFFffFFff

	for y = 0; y < frame.Height; y++ {
		i = y * frame.scanBytes / 4
		for x = 0; x < frame.Width; x++ {
			v = white * uint32((x+y)%2)
			buf[i] = v
			i++
			/*v = byte( y + x + frameCount )
			buf[i+0] = v    // B
			buf[i+1] = v	// G
			buf[i+2] = v	// R
			buf[i+3] = 0	// A - mostly ignored

			i += 4*/
		}
	}
}

func (r Renderer) SwitchBuffer() {
	targetA = !targetA
}

// TargetFrame should return the buffer to draw into
func (r *Renderer)TargetFrame() RenderFrame{
	if targetA {return frameA}
	return frameB
}

// RenderFrame should return the non-target buffer
func (r *Renderer) RenderFrame() RenderFrame {
	if targetA {return frameB}
	return frameA
}

func (f *RenderFrame) GetBufferPointer() *byte { return &f.bmp[0] }

// wordArray puns the byte buffer into a 32bit word array
func (f *RenderFrame) wordArray() []uint32 {
	newLen := len(f.bmp) / 4
	ptr := (*uint32)(unsafe.Pointer(&(f.bmp[0])))
	return unsafe.Slice(ptr, newLen) // go 1.17+
}
