package math32

// Sqrt returns the square root of x.
//
// On AMD64 with GOEXPERIMENT=simd, this uses the SQRTSS instruction via intrinsics.
// On other platforms, it uses float64 sqrt with conversion.
//
// Note: Special cases (NaN, Inf, negative values) are not explicitly handled
// for performance reasons. The behavior for these cases is determined by the
// underlying hardware/math.Sqrt implementation.
func Sqrt(x float32) float32 {
	return sqrt(x)
}
