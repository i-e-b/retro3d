package pkg

type Scene struct {
	Camera   *SceneCam
	Geometry []RefTriangle
	Points   []Vec3t
	Textures []*Texture
	Time     float64
}

func (s *Scene) AddCube() {
	base := len(s.Points)

	// Add points
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z:  1, U: 1, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z:  1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z:  1, U: 0, V: 1})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z:  1, U: 1, V: 1})

	s.Points = append(s.Points, Vec3t{X: -1, Y:  1, Z: -1, U: 1, V: 1})
	s.Points = append(s.Points, Vec3t{X: -1, Y: -1, Z: -1, U: 1, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y: -1, Z: -1, U: 0, V: 0})
	s.Points = append(s.Points, Vec3t{X: 1, Y:  1, Z: -1, U: 0, V: 1})

	// Stitch triangles
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 0, B: base + 1, C: base + 2, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 2, B: base + 3, C: base + 0, Tex: 0})

	s.Geometry = append(s.Geometry, RefTriangle{A: base + 4, B: base + 5, C: base + 6, Tex: 0})
	s.Geometry = append(s.Geometry, RefTriangle{A: base + 5, B: base + 6, C: base + 7, Tex: 0})
	// TODO: other faces
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

func (s *Scene) Project(screenWidth, screenHeight float64) []Vec3t{ // TODO: should be returning projected triangles (later)
	up := Vec3{0,1,0}
	world := &Mat4x4{}
	world.setLookAt(s.Camera.Position, s.Camera.Target, up)

	projt := &Mat4x4{}
	projt.setProjectionMatrix(0.78,0.01, 10.0)

	halfWidth := screenWidth / 2.0
	halfHeight := screenHeight / 2.0

	in := s.Points
	out := make([]Vec3t, len(s.Points))
	for i := 0; i < len(in); i++ {
		world.mulVec3t(&in[i], &out[i])
		projt.mulVec3t(&out[i], &out[i])

		// scale from -1..1 range and centre for screen
		out[i].X = out[i].X /* * halfWidth*/ + halfWidth
		out[i].Y = out[i].Y /* * halfHeight*/ + halfHeight
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