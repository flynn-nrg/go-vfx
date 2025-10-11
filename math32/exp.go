package math32

import "math"

// Exponential function constants
const (
	// For range reduction: exp(x) = 2^k * exp(r) where x = k*ln(2) + r
	ln2Hi = 6.9314575195e-01 // High bits of ln(2)
	ln2Lo = 1.4286067653e-06 // Low bits of ln(2)

	invLn2 = 1.44269504088896338700e+00 // 1/ln(2)

	// Overflow/underflow thresholds for float32
	expOverflow  = 88.72283905206835   // ln(2^128)
	expUnderflow = -103.97207708399179 // ln(2^-149) slightly below denormal
)

// Exp returns e**x, the base-e exponential of x.
//
// Special cases are:
//
//	Exp(+Inf) = +Inf
//	Exp(-Inf) = 0
//	Exp(NaN) = NaN
//
// Very large values overflow to +Inf.
// Very small values underflow to 0.
func Exp(x float32) float32 {
	// Handle special cases
	if IsNaN(x) {
		return x
	}

	if IsInf(x, 1) {
		return x // +Inf
	}

	if IsInf(x, -1) {
		return 0 // exp(-Inf) = 0
	}

	// Check for overflow/underflow
	if x > expOverflow {
		return float32(math.Inf(1)) // Overflow to +Inf
	}

	if x < expUnderflow {
		return 0 // Underflow to 0
	}

	// Range reduction: exp(x) = 2^k * exp(r)
	// where x = k*ln(2) + r and |r| < ln(2)/2

	// Find k
	k := int32(x*invLn2 + 0.5)
	if x < 0 {
		k = int32(x*invLn2 - 0.5)
	}

	// Compute r = x - k*ln(2) with extended precision
	hi := x - float32(k)*ln2Hi
	lo := float32(k) * ln2Lo
	r := hi - lo

	// Evaluate exp(r) using polynomial
	// For small r: exp(r) ≈ 1 + r + r²/2 + r³/6 + ...
	result := expKernel(r)

	// Multiply by 2^k using ldexp (efficient bit manipulation)
	return ldexp32(result, k)
}

// expKernel evaluates exp(r) for small r using polynomial approximation
// Valid for |r| < ln(2)/2 ≈ 0.347
func expKernel(r float32) float32 {
	// Use Taylor series: exp(r) = 1 + r + r²/2! + r³/3! + r⁴/4! + r⁵/5! + r⁶/6!
	// Rearranged for Horner's method

	const (
		c2 = 5.0000000000e-01 // 1/2!
		c3 = 1.6666666667e-01 // 1/3!
		c4 = 4.1666666667e-02 // 1/4!
		c5 = 8.3333333333e-03 // 1/5!
		c6 = 1.3888888889e-03 // 1/6!
	)

	// Horner's method: exp(r) = 1 + r + r²·(c2 + r·(c3 + r·(c4 + r·(c5 + r·c6))))
	r2 := r * r
	poly := c2 + r*(c3+r*(c4+r*(c5+r*c6)))
	return 1 + r + r2*poly
}

// ldexp32 returns x * 2^exp using efficient bit manipulation
// This is much faster than actual multiplication
func ldexp32(x float32, exp int32) float32 {
	// Handle special cases
	if x == 0 || IsNaN(x) || IsInf(x, 0) {
		return x
	}

	// Get bit representation
	bits := math.Float32bits(x)

	// Extract current exponent (bits 23-30)
	currentExp := int32((bits >> 23) & 0xFF)

	// Handle denormal numbers
	if currentExp == 0 {
		// Denormal number - normalize it first
		x *= float32(1 << 23) // Multiply by 2^23
		bits = math.Float32bits(x)
		currentExp = int32((bits>>23)&0xFF) - 23
	}

	// Add the exponent adjustment
	newExp := currentExp + exp

	// Check for overflow
	if newExp >= 0xFF {
		// Overflow to infinity
		if bits&0x80000000 != 0 {
			return float32(math.Inf(-1))
		}
		return float32(math.Inf(1))
	}

	// Check for underflow
	if newExp <= 0 {
		// Underflow - might be denormal or zero
		if newExp <= -24 {
			// Too small, return zero with correct sign
			if bits&0x80000000 != 0 {
				return float32(math.Copysign(0, -1))
			}
			return 0
		}

		// Denormal result - shift mantissa
		shift := uint32(-newExp + 1)
		mantissa := (bits & 0x7FFFFF) | 0x800000 // Add implicit leading 1
		mantissa >>= shift

		// Reconstruct with exponent = 0 (denormal)
		result := (bits & 0x80000000) | mantissa
		return math.Float32frombits(result)
	}

	// Normal case: just update exponent
	bits = (bits & 0x807FFFFF) | (uint32(newExp) << 23)
	return math.Float32frombits(bits)
}
