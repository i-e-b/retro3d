package pkg

import (
	"math"
	"unsafe"
)

type Renderer struct {
	width  int
	height int
	scene *Scene
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
	bmp        []uint32
	scanBytes  int
	frameBytes uint64
}

func NewRenderer(width, height int) *Renderer {
	frameA = makeFrame(width, height)
	frameB = makeFrame(width, height)
	targetA = true

	basicScene := NewScene()
	basicScene.AddCube()

	return &Renderer{
		width:  width,
		height: height,
		scene:  basicScene,
	}
}

func makeFrame(width, height int) RenderFrame {
	scanBytes := width * 4
	bmpSize := uint64(scanBytes) * uint64(height) // 32 bit argb
	bmp := make([]uint32, bmpSize)
	return RenderFrame{
		Width:  width,
		Height: height,
		bmp:    bmp,
		scanBytes: scanBytes,
		frameBytes: bmpSize,
	}
}

// TODO: move Update to a different place

// Update should update the world state and draw into the target buffer.
// You should switch buffers when you're finished.
func (r *Renderer) Update(t int64) {
	defer r.SwitchBuffer()

	// any world logic can go in here
	frame := r.TargetFrame()
	frame.Clear(0) // shouldn't be needed when rendering complete scenes

	r.scene.Advance(t)
	r.scene.Camera.Position = Vec3{
		math.Cos(r.scene.Time)*5, math.Sin(r.scene.Time)-1, math.Sin(r.scene.Time)*5,
	}

	scene := r.scene.Project(float64(frame.Width), float64(frame.Height))

	buf := frame.bmp
	size := frame.Width*frame.Height

	var white uint32 = 0xFFffFFff
	buf[int(r.scene.Time)] = white
	for i := 0; i < len(scene); i++ {
		v := scene[i]
		pi := int(v.Y)*frame.Width + int(v.X)
		if pi >= 0 && pi < size{
			buf[pi] = white
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

func (f *RenderFrame) GetBufferPointer() uintptr {
	return uintptr(unsafe.Pointer(&f.bmp[0]))
}

func (f *RenderFrame) Clear(color uint32) {
	var y,x int
	var i int

	for y = 0; y < f.Height; y++ {
		i = y * f.Width
		for x = 0; x < f.Width; x++ {
			f.bmp[i] = color
			i++
		}
	}
}
