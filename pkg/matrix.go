package pkg

import "math"

type Mat4x4 struct {
	m00, m01, m02, m03 float64
	m10, m11, m12, m13 float64
	m20, m21, m22, m23 float64
	m30, m31, m32, m33 float64
}

func (M *Mat4x4)setProjectionMatrix(viewAngle float64, nearDistance float64, farDistance float64) {
	// set the projection matrix
	scale := 1.0 / math.Tan(viewAngle*0.5*math.Pi/180.0)
	fieldOfView := farDistance - nearDistance
	M.m00 = scale                                     // scale the X coordinates of the projected point
	M.m11 = scale                                     // scale the Y coordinates of the projected point
	M.m22 = -farDistance / fieldOfView                // used to remap Z to [0,1]
	M.m32 = -farDistance * nearDistance / fieldOfView // used to remap Z [0,1]
	M.m23 = -1                                        // set w = -Z
	M.m33 = 0
}

func (M *Mat4x4) mulVec3t(in, out *Vec3t) {
	// out = in * M;
	out.X = in.X*M.m00 + in.Y*M.m10 + in.Z*M.m20 + /* in.Z = 1 */ M.m30
	out.Y = in.X*M.m01 + in.Y*M.m11 + in.Z*M.m21 + /* in.Z = 1 */ M.m31
	out.Z = in.X*M.m02 + in.Y*M.m12 + in.Z*M.m22 + /* in.Z = 1 */ M.m32
	w    := in.X*M.m03 + in.Y*M.m13 + in.Z*M.m23 + /* in.Z = 1 */ M.m33

	// normalize if w is not 1 (convert to Cartesian coordinates)
	if w != 1 {
		out.X /= w
		out.Y /= w
		out.Z /= w
	}
	out.U = in.U
	out.V = in.V
}

func (M *Mat4x4) FPSViewRH(eye Vec3, pitch, yaw float64) {
	// I assume the values are already converted to radians.
	cosPitch := math.Cos(pitch)
	sinPitch := math.Sin(pitch)
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	xaxis := Vec3{cosYaw, 0, -sinYaw}
	yaxis := Vec3{sinYaw * sinPitch, cosPitch, cosYaw * sinPitch}
	zaxis := Vec3{sinYaw * cosPitch, -sinPitch, cosPitch * cosYaw}

	eX := DotV3(xaxis, eye)
	eY := DotV3(yaxis, eye)
	eZ := DotV3(zaxis, eye)

	// Create a 4x4 view matrix from the right, up, forward and eye position vectors
	M.m00 = xaxis.X;		M.m01 = yaxis.X;		M.m02 = zaxis.X;		M.m03 = 0
	M.m10 = xaxis.Y;		M.m11 = yaxis.Y;		M.m12 = zaxis.Y;		M.m13 = 0
	M.m20 = xaxis.Z;		M.m21 = yaxis.Z;		M.m22 = zaxis.Z;		M.m23 = 0
	M.m30 = eX;				M.m31 = eY;				M.m32 = eZ;				M.m33 = 1
}

// setLookAt sets the matrix to translate by position, then rotate toward target and up.
func (M *Mat4x4)setLookAt(position, target Vec3){
	vv := SubV3(target, position) // view vector
	zx := NormalV3(vv)

	// calculate the 'up' vector. This is not in world space, but
	// is perpendicular to the view vector.
	horizon := Vec3{1,0,1}
	up := NormalV3(CrossV3(vv, horizon))
	//if up.Y > 0 {up.Y = -up.Y}
	//up := Vec3{0,1,0}

	xx := NormalV3(CrossV3(up, position))
	yx := CrossV3(zx, xx)
	ne := InvV3(position)

	eX := - DotV3(xx, ne)
	eY := - DotV3(yx, ne)
	eZ := - DotV3(zx, ne)

	M.m00 = xx.X;		M.m01 = yx.X;		M.m02 = zx.X;		M.m03 = 0
	M.m10 = xx.Y;		M.m11 = yx.Y;		M.m12 = zx.Y;		M.m13 = 0
	M.m20 = xx.Z;		M.m21 = yx.Z;		M.m22 = zx.Z;		M.m23 = 0
	M.m30 = eX;			M.m31 = eY;			M.m32 = eZ;			M.m33 = 1
}