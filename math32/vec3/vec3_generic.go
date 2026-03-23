//go:build !(amd64 && goexperiment.simd)

package vec3

import "github.com/flynn-nrg/go-vfx/math32"

// addTwo returns the sum of two vectors (scalar implementation).
func addTwo(v1, v2 Vec3Impl) Vec3Impl {
	return Vec3Impl{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

// subTwo returns the difference of two vectors (scalar implementation).
func subTwo(v1, v2 Vec3Impl) Vec3Impl {
	return Vec3Impl{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

// mulTwo returns the element-wise product of two vectors (scalar implementation).
func mulTwo(v1, v2 Vec3Impl) Vec3Impl {
	return Vec3Impl{
		X: v1.X * v2.X,
		Y: v1.Y * v2.Y,
		Z: v1.Z * v2.Z,
	}
}

// divTwo returns the element-wise division of two vectors (scalar implementation).
func divTwo(v1, v2 Vec3Impl) Vec3Impl {
	return Vec3Impl{
		X: v1.X / v2.X,
		Y: v1.Y / v2.Y,
		Z: v1.Z / v2.Z,
	}
}

// scalarMulSIMD returns the scalar multiplication (scalar implementation).
func scalarMulSIMD(v Vec3Impl, t float32) Vec3Impl {
	return Vec3Impl{
		X: v.X * t,
		Y: v.Y * t,
		Z: v.Z * t,
	}
}

// scalarDivSIMD returns the scalar division (scalar implementation).
func scalarDivSIMD(v Vec3Impl, t float32) Vec3Impl {
	return Vec3Impl{
		X: v.X / t,
		Y: v.Y / t,
		Z: v.Z / t,
	}
}

// dotSIMD computes dot product (scalar implementation).
func dotSIMD(v1, v2 Vec3Impl) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// crossSIMD computes cross product (scalar implementation).
func crossSIMD(v1, v2 Vec3Impl) Vec3Impl {
	return Vec3Impl{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

// lerpSIMD performs linear interpolation (scalar implementation).
func lerpSIMD(v0, v1 Vec3Impl, t float32) Vec3Impl {
	return Vec3Impl{
		X: (1-t)*v0.X + t*v1.X,
		Y: (1-t)*v0.Y + t*v1.Y,
		Z: (1-t)*v0.Z + t*v1.Z,
	}
}

// squaredLengthSIMD returns the squared length (scalar implementation).
func squaredLengthSIMD(v Vec3Impl) float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// lengthSIMD returns the length (scalar implementation).
func lengthSIMD(v Vec3Impl) float32 {
	return math32.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// unitVectorSIMD returns a normalized vector (scalar implementation).
func unitVectorSIMD(v Vec3Impl) Vec3Impl {
	l := math32.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vec3Impl{
		X: v.X / l,
		Y: v.Y / l,
		Z: v.Z / l,
	}
}
