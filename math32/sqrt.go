package math32

import "math"

// Sqrt returns the square root of x.
//
// This function uses float64 sqrt with conversion, which is faster than
// native float32 sqrt on most modern hardware (especially ARM64).
//
// Note: Special cases (NaN, Inf, negative values) are not explicitly handled
// for performance reasons. The behavior for these cases is determined by the
// underlying math.Sqrt implementation.
func Sqrt(x float32) float32 {
	// Use float64 sqrt and convert back - benchmarks show this is ~7x faster
	// on ARM64 and competitive on AMD64
	return float32(math.Sqrt(float64(x)))
}
