package pkg

import "math"

type Vec3 struct {
	X,Y,Z float64
}

type Vec3t struct {
	X, Y, Z float64
	U, V    uint32
}

// DotV3 gives the dot-product / scalar-product of two 3D vectors
func DotV3(a,b Vec3) float64{
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// CrossV3 gives the cross-product / vector-product of two 3D vectors
func CrossV3(a,b Vec3) Vec3{
	return Vec3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y-b.X,
	}
}

// NormalV3 returns a normalised copy of a vector
func NormalV3(a Vec3) Vec3{
	length := math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	if length == 0 { return a }
	return Vec3{
		X: a.X / length,
		Y: a.Y / length,
		Z: a.Z / length,
	}
}

// SubV3 give component-wise: a - b
func SubV3(a,b Vec3) Vec3{
	return Vec3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}
// InvV3 give component-wise: -a
func InvV3(a Vec3) Vec3{
	return Vec3{
		X: -a.X,
		Y: -a.Y,
		Z: -a.Z,
	}
}

// LengthV3 give component-wise: |a|
func LengthV3(a Vec3) float64{
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}
// SignZV3 gives -1 or +1 based on the sign of the Z component
func SignZV3(a Vec3) float64{
	if a.Z >= 0 {return 1}
	return -1
}