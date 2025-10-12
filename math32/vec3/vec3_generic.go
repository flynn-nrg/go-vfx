//go:build !amd64 && !arm64

package vec3

// dot provides a software fallback for architectures without assembly implementation.
func dot(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) float32 {
	return (v1X * v2X) + (v1Y * v2Y) + (v1Z * v2Z)
}

// add3 provides a software fallback for architectures without assembly implementation.
func add3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32) {
	return v1X + v2X, v1Y + v2Y, v1Z + v2Z
}

// sub3 provides a software fallback for architectures without assembly implementation.
func sub3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32) {
	return v1X - v2X, v1Y - v2Y, v1Z - v2Z
}

// mul3 provides a software fallback for architectures without assembly implementation.
func mul3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32) {
	return v1X * v2X, v1Y * v2Y, v1Z * v2Z
}

// scalarMul3 provides a software fallback for architectures without assembly implementation.
func scalarMul3(vX, vY, vZ, scalar float32) (float32, float32, float32) {
	return vX * scalar, vY * scalar, vZ * scalar
}

// cross provides a software fallback for architectures without assembly implementation.
func cross(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32) {
	return (v1Y*v2Z - v1Z*v2Y), -(v1X*v2Z - v1Z*v2X), (v1X*v2Y - v1Y*v2X)
}
