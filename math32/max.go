package math32

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN) = Max(NaN, x) = NaN
//
// Note: The sign of zero (±0) is hardware-dependent and may vary.
func Max(x, y float32) float32 {
	// Handle NaN - if either is NaN, return NaN
	if IsNaN(x) || IsNaN(y) {
		return NaN()
	}

	// Handle infinities
	if IsInf(x, 1) || IsInf(y, 1) {
		if IsInf(x, 1) {
			return x
		}
		return y
	}

	// Use hardware instruction via assembly
	// Note: Signed zero behavior (+0 vs -0) is hardware-dependent
	// and doesn't matter for VFX/graphics applications
	return max(x, y)
}

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//
// Note: The sign of zero (±0) is hardware-dependent and may vary.
func Min(x, y float32) float32 {
	// Handle NaN - if either is NaN, return NaN
	if IsNaN(x) || IsNaN(y) {
		return NaN()
	}

	// Handle infinities
	if IsInf(x, -1) || IsInf(y, -1) {
		if IsInf(x, -1) {
			return x
		}
		return y
	}

	// Use hardware instruction via assembly
	// Note: Signed zero behavior (+0 vs -0) is hardware-dependent
	// and doesn't matter for VFX/graphics applications
	return min(x, y)
}
