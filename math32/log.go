package math32

import "math"

// Logarithm constants
const (
	// sqrt(2) for range reduction
	sqrt2 = 1.41421356237309504880
)

// Log returns the natural logarithm of x.
//
// Special cases are:
//
//	Log(+Inf) = +Inf
//	Log(0) = -Inf
//	Log(x < 0) = NaN
//	Log(NaN) = NaN
func Log(x float32) float32 {
	// Handle special cases
	if IsNaN(x) {
		return x
	}

	if x < 0 {
		return NaN() // log of negative is undefined
	}

	if x == 0 {
		return Inf(-1) // log(0) = -Inf
	}

	if IsInf(x, 1) {
		return x // log(+Inf) = +Inf
	}

	// Extract exponent and mantissa using bit manipulation
	bits := math.Float32bits(x)

	// Get exponent (bits 23-30)
	exp := int32((bits >> 23) & 0xFF)

	// Handle denormal numbers
	if exp == 0 {
		// Denormal - normalize it
		x *= float32(1 << 23)
		bits = math.Float32bits(x)
		exp = int32((bits>>23)&0xFF) - 23
	}

	// Unbias exponent (float32 uses bias of 127)
	k := exp - 127

	// Get mantissa and set exponent to 127 (value in [1, 2))
	bits = (bits & 0x807FFFFF) | (127 << 23)
	f := math.Float32frombits(bits)

	// Now f is in [1, 2)
	// Reduce to [sqrt(2)/2, sqrt(2)] for better polynomial accuracy
	if f > float32(sqrt2) {
		f *= 0.5
		k++
	}

	// Compute log(f) where f is in [sqrt(2)/2, sqrt(2)]
	// Transform to [-1/3, 1/3] using: s = (f-1)/(f+1)
	s := (f - 1) / (f + 1)

	// log(f) = 2s + 2s³/3 + 2s⁵/5 + ... = 2s(1 + s²/3 + s⁴/5 + ...)
	logF := logKernel(s)

	// Final result: log(x) = k·ln(2) + log(f)
	return float32(k)*ln2Hi + (float32(k)*ln2Lo + logF)
}

// logKernel evaluates log(f) for f in [sqrt(2)/2, sqrt(2)]
// Uses transformation s = (f-1)/(f+1) to reduce to smaller range
func logKernel(s float32) float32 {
	// For s = (f-1)/(f+1), we have:
	// log(f) = log((1+s)/(1-s)) = 2s(1 + s²/3 + s⁴/5 + s⁶/7 + ...)
	// This series converges rapidly for |s| < 1/3

	const (
		L1 = 6.666666666666735130e-01 // 2/3
		L2 = 3.999999999940941908e-01 // 2/5
		L3 = 2.857142874366239149e-01 // 2/7
		L4 = 2.222219843214978396e-01 // 2/9
		L5 = 1.818357216161805012e-01 // 2/11
		L6 = 1.531383769920937332e-01 // 2/13
	)

	s2 := s * s

	// Polynomial: 2s(1 + s²(L1 + s²(L2 + s²(L3 + s²(L4 + s²(L5 + s²·L6))))))
	poly := L1 + s2*(L2+s2*(L3+s2*(L4+s2*(L5+s2*L6))))

	return 2 * s * (1 + s2*poly)
}
