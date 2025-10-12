//go:build arm64

package mat3

// matrixVectorMul is implemented in mat3_arm64.s using NEON instructions
func matrixVectorMul(a11, a12, a13, a21, a22, a23, a31, a32, a33, vX, vY, vZ float32) (float32, float32, float32)
