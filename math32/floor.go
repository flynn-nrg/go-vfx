package math32

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN
func Floor(x float32) float32 {
	// Handle special cases
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}

	return floor(x)
}

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//
//	Ceil(±0) = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN) = NaN
func Ceil(x float32) float32 {
	// Handle special cases
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}

	return ceil(x)
}

// Round returns the nearest integer, rounding half away from zero.
//
// Special cases are:
//
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round(x float32) float32 {
	// Handle special cases
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}

	return round(x)
}
