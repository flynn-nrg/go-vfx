package math32

import "math"

// Sqrt returns the square root of x.
//
// This function uses float64 sqrt with conversion, which is faster than
// native float32 sqrt on most modern hardware (especially ARM64).
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}
