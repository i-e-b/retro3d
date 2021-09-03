package pkg

type Scene struct {
	Camera   *SceneCam
	Geometry []RefTriangle
	Points   []Vec3t
	Textures []*Texture
	Time     float64
}

func (s *Scene)AddFancyCube(){
	texIdx := len(s.Textures)
	tex := LoadFromPng("img/tex.png")
	s.Textures = append(s.Textures, &tex)

	base := len(s.Points)
	w := tex.Width
	h := tex.Height

	// Add points
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z:  1, U: w, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z:  1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z:  1, U: 0, V: h})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z:  1, U: w, V: h})

	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z: -1, U: w, V: h})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z: -1, U: w, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z: -1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z: -1, U: 0, V: h})

	// Stitch triangles
	// back
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 0, B: base + 1, C: base + 2, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 2, B: base + 3, C: base + 0, Tex: texIdx})
	// front
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 4, B: base + 5, C: base + 6, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 7, C: base + 4, Tex: texIdx})
	// top
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 1, B: base + 4, C: base + 7, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 7, B: base + 0, C: base + 1, Tex: texIdx})
	// bottom
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 2, B: base + 5, C: base + 6, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 3, C: base + 2, Tex: texIdx})
	// left
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 1, B: base + 4, C: base + 5, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 5, B: base + 2, C: base + 1, Tex: texIdx})
	// right
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 0, B: base + 7, C: base + 6, Tex: texIdx})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 3, C: base + 0, Tex: texIdx})
}

func (s *Scene) AddCube() {
	base := len(s.Points)

	// Add points
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z:  1, U: 3, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z:  1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z:  1, U: 0, V: 3})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z:  1, U: 3, V: 3})

	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z: -1, U: 3, V: 3})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z: -1, U: 3, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z: -1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z: -1, U: 0, V: 3})

	// Stitch triangles
	// back
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 0, B: base + 1, C: base + 2, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 2, B: base + 3, C: base + 0, Tex: 0})
	// front
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 4, B: base + 5, C: base + 6, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 7, C: base + 4, Tex: 0})
	// top
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 1, B: base + 4, C: base + 7, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 7, B: base + 0, C: base + 1, Tex: 0})
	// bottom
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 2, B: base + 5, C: base + 6, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 3, C: base + 2, Tex: 0})
	// left
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 1, B: base + 4, C: base + 5, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 5, B: base + 2, C: base + 1, Tex: 0})
	// right
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 0, B: base + 7, C: base + 6, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 6, B: base + 3, C: base + 0, Tex: 0})
}

func NewScene() *Scene {
	nt := NullTexture()
	cam := SceneCam{
		Position: Vec3{0, 0, 5},
		Target:   Vec3{0, 0, 0},
	}
	return &Scene{
		Camera:   &cam,
		Geometry: []RefTriangle{},
		Points:   []Vec3t{},
		Textures: []*Texture{&nt},
	}
}

func (s *Scene) ProjectPoints(screenWidth, screenHeight float64) []Vec3t{
	up := Vec3{0,0,1}
	world := &Mat4x4{}
	world.setLookAt(s.Camera.Position, s.Camera.Target, up)

	projt := &Mat4x4{}
	projt.setProjectionMatrix(0.78,0.01, 10000.0)

	halfWidth := screenWidth / 2.0
	halfHeight := screenHeight / 2.0

	// Project all the points
	in := s.Points
	out := make([]Vec3t, len(s.Points))
	for i := 0; i < len(in); i++ {
		world.mulVec3t(&in[i], &out[i])
		upz := out[i].Z // save unprojected Z for texture mapping and depth sorting
		projt.mulVec3t(&out[i], &out[i])

		// scale from -1..1 range and centre for screen
		out[i].X = out[i].X /* * halfWidth*/ + halfWidth
		out[i].Y = out[i].Y /* * halfHeight*/ + halfHeight
		out[i].Z = upz
	}
	return out
}

func (s *Scene) Advance(t int64) {
	s.Time += float64(t) / 1000.0
}

type RefTriangle struct {
	// A, B, and C are indexes to the scene's Points list.
	A,B,C int
	// Tex is an index to the scene's Textures list.
	Tex int
}

type SceneCam struct {
	Position Vec3
	Target Vec3
}