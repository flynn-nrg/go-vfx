package math32

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN) = Max(NaN, x) = NaN
//	Max(+0, ±0) = Max(±0, +0) = +0
//	Max(-0, -0) = -0
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
	return max(x, y)
}

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//	Min(-0, ±0) = Min(±0, -0) = -0
//	Min(+0, +0) = +0
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
	return min(x, y)
}
