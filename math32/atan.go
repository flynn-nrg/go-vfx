package math32

import "math"

// Atan returns the arctangent, in radians, of x.
//
// Special cases are:
//
//	Atan(±0) = ±0
//	Atan(±Inf) = ±π/2
//	Atan(NaN) = NaN
//
// The result is in the range [-π/2, π/2].
func Atan(x float32) float32 {
	// Handle special cases
	if x == 0 {
		return x // preserves sign of zero
	}

	// Handle NaN
	if IsNaN(x) {
		return x
	}

	// Handle infinity
	if IsInf(x, 1) {
		return float32(Pi / 2)
	}
	if IsInf(x, -1) {
		return float32(-Pi / 2)
	}

	// Extract sign and work with absolute value
	sign := false
	if x < 0 {
		x = -x
		sign = true
	}

	var result float32

	// Range reduction using atan identities
	// This keeps the argument small for accurate polynomial approximation
	if x > 2.414213562373095 { // tan(3π/8) ≈ 2.414
		// For very large x, use: atan(x) = π/2 - atan(1/x)
		result = float32(Pi/2) - atanPoly(1/x)
	} else if x > 0.4142135623730951 { // tan(π/8) ≈ 0.414
		// For medium x, use: atan(x) = π/4 + atan((x-1)/(x+1))
		result = float32(Pi/4) + atanPoly((x-1)/(x+1))
	} else {
		// For small x, use direct polynomial approximation
		result = atanPoly(x)
	}

	if sign {
		return -result
	}
	return result
}

// atanPoly evaluates the polynomial approximation for atan
// Valid for |x| < tan(π/8) ≈ 0.414
// Uses minimax polynomial coefficients optimized for float32 precision
func atanPoly(x float32) float32 {
	// For small x, atan(x) ≈ x - x³/3 + x⁵/5 - x⁷/7 + ...
	// We use a rational approximation for better accuracy
	x2 := x * x

	// Minimax polynomial coefficients for atan on [-0.414, 0.414]
	// These coefficients give ~1 ULP accuracy for float32
	const (
		p0 = -3.333314528e-01 // ≈ -1/3
		p1 = 1.999355085e-01  // ≈ 1/5
		p2 = -1.420889944e-01 // ≈ -1/7
		p3 = 1.065626393e-01
		p4 = -7.522124857e-02
		p5 = 4.263936017e-02
		p6 = -1.480085629e-02
		p7 = 2.417283948e-03
	)

	// Evaluate using Horner's method
	// atan(x) ≈ x + x³·P(x²)
	poly := p0 + x2*(p1+x2*(p2+x2*(p3+x2*(p4+x2*(p5+x2*(p6+x2*p7))))))
	return x + x*x2*poly
}

// Atan2 returns the arc tangent of y/x, using the signs of the two to
// determine the quadrant of the return value.
//
// Special cases are (in order):
//
//	Atan2(y, NaN) = NaN
//	Atan2(NaN, x) = NaN
//	Atan2(±0, x>=0) = ±0
//	Atan2(±0, x<0) = ±π
//	Atan2(y>0, 0) = +π/2
//	Atan2(y<0, 0) = -π/2
//	Atan2(±Inf, +Inf) = ±π/4
//	Atan2(±Inf, -Inf) = ±3π/4
//	Atan2(y, +Inf) = 0
//	Atan2(y>0, -Inf) = +π
//	Atan2(y<0, -Inf) = -π
//	Atan2(±Inf, x) = ±π/2
//	Atan2(y, x>0) = Atan(y/x)
//	Atan2(y>=0, x<0) = Atan(y/x) + π
//	Atan2(y<0, x<0) = Atan(y/x) - π
//
// The result is in the range [-π, π].
func Atan2(y, x float32) float32 {
	// Handle NaN
	if IsNaN(y) || IsNaN(x) {
		return float32(math.NaN())
	}

	// Handle y = ±Inf
	if IsInf(y, 0) {
		if IsInf(x, 1) {
			// y = ±Inf, x = +Inf
			if y > 0 {
				return float32(Pi / 4) // +π/4
			}
			return float32(-Pi / 4) // -π/4
		}
		if IsInf(x, -1) {
			// y = ±Inf, x = -Inf
			if y > 0 {
				return float32(3 * Pi / 4) // +3π/4
			}
			return float32(-3 * Pi / 4) // -3π/4
		}
		// y = ±Inf, x = finite
		if y > 0 {
			return float32(Pi / 2) // +π/2
		}
		return float32(-Pi / 2) // -π/2
	}

	// Handle x = ±Inf (y is finite)
	if IsInf(x, 1) {
		// x = +Inf
		if y > 0 {
			return 0
		}
		if y < 0 {
			return float32(math.Copysign(0, -1)) // -0
		}
		return y // ±0
	}
	if IsInf(x, -1) {
		// x = -Inf
		if y > 0 {
			return float32(Pi)
		}
		if y < 0 {
			return float32(-Pi)
		}
		// y = ±0, x = -Inf
		if Signbit(y) {
			return float32(-Pi)
		}
		return float32(Pi)
	}

	// Handle x = 0
	if x == 0 {
		if y > 0 {
			return float32(Pi / 2)
		}
		if y < 0 {
			return float32(-Pi / 2)
		}
		// y = 0, x = 0
		return y // preserve sign of y
	}

	// Handle y = 0
	if y == 0 {
		if x > 0 {
			return y // ±0
		}
		// x < 0
		if Signbit(y) {
			return float32(-Pi)
		}
		return float32(Pi)
	}

	// General case: both x and y are finite and non-zero
	// Determine quadrant and compute angle
	q := Atan(y / x)

	if x > 0 {
		// Quadrants I and IV
		return q
	}

	// Quadrants II and III (x < 0)
	if y > 0 {
		// Quadrant II
		return q + float32(Pi)
	}
	// Quadrant III
	return q - float32(Pi)
}
