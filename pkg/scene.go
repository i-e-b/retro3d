package pkg

type Scene struct {
	Camera SceneCam
	Geometry []RefTriangle
	Points []Vec3t
	Textures []Texture
}

func NewScene() *Scene {
	return &Scene{
		Camera:   SceneCam{
			Position: Vec3{0, 0, -10},
			Target:   Vec3{0, 0, 0},
		},
		Geometry: []RefTriangle{},
		Points:   []Vec3t{},
		Textures: []Texture{},
	}

}

type RefTriangle struct {
	// A, B, and C are indexes to the scene's Points list.
	A,B,C int
	// Tex is an index to the scene's Textures list.
	Tex int
}

type Texture struct {
	Width      int
	Height     int
	Bmp        [][]uint32
	ScanBytes  int
}

type SceneCam struct {
	Position Vec3
	Target Vec3
}