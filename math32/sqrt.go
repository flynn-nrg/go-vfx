package math32

// Sqrt returns the square root of x.
//
// Special cases are:
//
//	Sqrt(+Inf) = +Inf
//	Sqrt(±0) = ±0
//	Sqrt(x < 0) = NaN
//	Sqrt(NaN) = NaN
//
// This function uses hardware instructions on ARM64 and AMD64 when available,
// falling back to software implementation on other architectures.
func Sqrt(x float32) float32 {
	// Handle special cases
	switch {
	case x == 0 || x != x: // ±0 or NaN
		return x
	case x < 0:
		return NaN()
	case IsInf(x, 1):
		return x
	}

	return sqrt(x)
}

// sqrt is implemented in assembly for ARM64 and AMD64 (see sqrt_arm64.s, sqrt_amd64.s).
// For other architectures, it uses the software fallback in sqrt_generic.go.
// The function declaration is in architecture-specific stub files.
