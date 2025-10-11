package math32

// Asin returns the arcsine, in radians, of x.
//
// Special cases are:
//
//	Asin(±0) = ±0
//	Asin(x) = NaN if x < -1 or x > 1
func Asin(x float32) float32 {
	// Handle special cases
	if x == 0 {
		return x // preserves sign of zero
	}

	// Handle NaN
	if IsNaN(x) {
		return x
	}

	// Extract sign and work with absolute value
	sign := false
	if x < 0 {
		x = -x
		sign = true
	}

	// Out of domain
	if x > 1 {
		return NaN()
	}

	var result float32

	if x <= 0.5 {
		// For small x, use direct polynomial approximation
		// asin(x) ≈ x + x³·P(x²)
		// This is the fast path that avoids sqrt
		x2 := x * x
		result = x + x*x2*asinPoly(x2)
	} else {
		// For large x (0.5 < x ≤ 1), use the identity:
		// asin(x) = π/2 - 2·asin(sqrt((1-x)/2))
		// This avoids the singularity at x = 1
		z := (1 - x) * 0.5
		s := Sqrt(z) // Use hardware-accelerated float32 sqrt
		x2 := z
		result = float32(Pi/2) - 2*(s+s*x2*asinPoly(x2))
	}

	if sign {
		return -result
	}
	return result
}

// asinPoly evaluates the polynomial approximation for asin on [0, 0.5]
// Uses minimax polynomial coefficients optimized for float32 precision
// Formula: asin(x) ≈ x + x³·P(x²) where P is this polynomial
func asinPoly(x2 float32) float32 {
	// Minimax polynomial coefficients for asin on [0, 0.5]
	// These coefficients give ~1 ULP accuracy for float32
	// Derived from Remez algorithm approximation
	const (
		p0 = 1.6666586697e-01 // ≈ 1/6
		p1 = 7.4953002686e-02
		p2 = 4.5470025998e-02
		p3 = 2.4181311049e-02
		p4 = 4.2163199048e-02
	)

	// Evaluate using Horner's method for numerical stability and efficiency
	// Takes advantage of FMA instructions on modern CPUs
	return p0 + x2*(p1+x2*(p2+x2*(p3+x2*p4)))
}

// Acos returns the arccosine, in radians, of x.
//
// Special cases are:
//
//	Acos(x) = NaN if x < -1 or x > 1
//	Acos(1) = 0
//	Acos(-1) = π
//
// The result is in the range [0, π].
func Acos(x float32) float32 {
	// Handle NaN
	if IsNaN(x) {
		return x
	}

	// Out of domain
	if x < -1 || x > 1 {
		return NaN()
	}

	// Special cases for exact values
	if x == 1 {
		return 0
	}
	if x == -1 {
		return float32(Pi)
	}

	// For |x| > 0.5, use direct computation to avoid cancellation error
	// This is more accurate than π/2 - asin(x) near x = ±1
	if x > 0.5 {
		// For x > 0.5, use: acos(x) = 2*asin(sqrt((1-x)/2))
		z := (1 - x) * 0.5
		s := Sqrt(z)
		x2 := z
		return 2 * (s + s*x2*asinPoly(x2))
	}

	if x < -0.5 {
		// For x < -0.5, use: acos(x) = π - 2*asin(sqrt((1+x)/2))
		z := (1 + x) * 0.5
		s := Sqrt(z)
		x2 := z
		return float32(Pi) - 2*(s+s*x2*asinPoly(x2))
	}

	// For |x| <= 0.5, use: acos(x) = π/2 - asin(x)
	// This avoids sqrt and uses the fast polynomial path of asin
	return float32(Pi/2) - Asin(x)
}
