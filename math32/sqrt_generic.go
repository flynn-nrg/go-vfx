//go:build !(amd64 && goexperiment.simd)

package math32

import "math"

// sqrt provides a software fallback using float64 conversion.
// On AMD64 with GOEXPERIMENT=simd, this is replaced by archsimd intrinsics.
func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
