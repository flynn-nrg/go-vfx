//go:build !amd64.v3 && !arm64

package mat3

// matrixVectorMul provides a software fallback for architectures without assembly implementation.
func matrixVectorMul(a11, a12, a13, a21, a22, a23, a31, a32, a33, vX, vY, vZ float32) (float32, float32, float32) {
	return a11*vX + a12*vY + a13*vZ,
		a21*vX + a22*vY + a23*vZ,
		a31*vX + a32*vY + a33*vZ
}
