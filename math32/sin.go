package math32

// Trigonometric reduction constants for Cody-Waite method
const (
	// We need to reduce x modulo 2π
	// Split π/2 into high and low parts for extended precision
	pi2Hi = 1.5703125        // High 24 bits of π/2
	pi2Lo = 4.83826794897e-4 // Low bits of π/2 (π/2 - pi2Hi)

	// Inverse of π/2 for determining quadrant
	invPi2 = 6.36619772367581382433e-01 // 2/π
)

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func Sin(x float32) float32 {
	// Handle special cases
	if x == 0 {
		return x // preserves sign of zero
	}

	if IsNaN(x) {
		return x
	}

	if IsInf(x, 0) {
		return NaN()
	}

	// Argument reduction
	neg := false
	if x < 0 {
		x = -x
		neg = true
	}

	// Reduce to [0, π/2] and determine quadrant
	j := uint32(x * invPi2)
	y := x - float32(j)*pi2Hi - float32(j)*pi2Lo

	// j&3 determines which quadrant (and thus which function and sign to use)
	// Quadrant 0: sin(y)
	// Quadrant 1: cos(y)
	// Quadrant 2: -sin(y)
	// Quadrant 3: -cos(y)

	var result float32
	switch j & 3 {
	case 0:
		result = sinKernel(y)
	case 1:
		result = cosKernel(y)
	case 2:
		result = -sinKernel(y)
	case 3:
		result = -cosKernel(y)
	}

	if neg {
		return -result
	}
	return result
}

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func Cos(x float32) float32 {
	// Handle special cases
	if IsNaN(x) {
		return x
	}

	if IsInf(x, 0) {
		return NaN()
	}

	// Make x positive (cos is even function)
	if x < 0 {
		x = -x
	}

	// Reduce to [0, π/2] and determine quadrant
	j := uint32(x * invPi2)
	y := x - float32(j)*pi2Hi - float32(j)*pi2Lo

	// For cos, quadrants are shifted by 1 from sin:
	// Quadrant 0: cos(y)
	// Quadrant 1: -sin(y)
	// Quadrant 2: -cos(y)
	// Quadrant 3: sin(y)

	switch j & 3 {
	case 0:
		return cosKernel(y)
	case 1:
		return -sinKernel(y)
	case 2:
		return -cosKernel(y)
	case 3:
		return sinKernel(y)
	}

	return 1 // unreachable
}

// sinKernel evaluates sin for arguments in [0, π/2]
// Uses minimax polynomial approximation
func sinKernel(x float32) float32 {
	// For small x, use Taylor series: sin(x) = x - x³/6 + x⁵/120 - ...
	// Rearranged as: sin(x) = x + x³·P(x²)

	const (
		S1 = -1.66666666666666324348e-01 // -1/6
		S2 = 8.33333333332248946124e-03  //  1/120
		S3 = -1.98412698298579493134e-04 // -1/5040
		S4 = 2.75573137070700676789e-06  //  1/362880
		S5 = -2.50507602534068634195e-08
	)

	z := x * x
	return x + x*z*(S1+z*(S2+z*(S3+z*(S4+z*S5))))
}

// cosKernel evaluates cos for arguments in [0, π/2]
// Uses minimax polynomial approximation
func cosKernel(x float32) float32 {
	// For small x, use Taylor series: cos(x) = 1 - x²/2 + x⁴/24 - ...
	// Rearranged as: cos(x) = 1 + x²·P(x²)

	const (
		C1 = 4.16666666666666019037e-02  //  1/24
		C2 = -1.38888888888741095749e-03 // -1/720
		C3 = 2.48015872894767294178e-05  //  1/40320
		C4 = -2.75573143513906633035e-07 // -1/3628800
		C5 = 2.08757232129817482790e-09
	)

	z := x * x
	return 1.0 - 0.5*z + z*z*(C1+z*(C2+z*(C3+z*(C4+z*C5))))
}

// Tan returns the tangent of the radian argument x.
//
// Special cases are:
//
//	Tan(±0) = ±0
//	Tan(±Inf) = NaN
//	Tan(NaN) = NaN
func Tan(x float32) float32 {
	// Handle special cases
	if x == 0 {
		return x // preserves sign of zero
	}

	if IsNaN(x) {
		return x
	}

	if IsInf(x, 0) {
		return NaN()
	}

	// Argument reduction
	neg := false
	if x < 0 {
		x = -x
		neg = true
	}

	// Reduce to [0, π/2] and determine quadrant
	j := uint32(x * invPi2)
	y := x - float32(j)*pi2Hi - float32(j)*pi2Lo

	// For tan, we use tan(x) = sin(x) / cos(x)
	// But we need to handle the quadrants carefully
	// Quadrant 0: tan(y) = sin(y) / cos(y)
	// Quadrant 1: tan(y) = cos(y) / (-sin(y)) = -cot(y)
	// And tan has period π, so quadrants 2,3 are same as 0,1

	var result float32
	if (j & 1) == 0 {
		// Even quadrant: use sin/cos directly
		s := sinKernel(y)
		c := cosKernel(y)
		result = s / c
	} else {
		// Odd quadrant: tan(π/2 - x) = cot(x) = cos/sin
		// But we want -cot, so negate
		s := sinKernel(y)
		c := cosKernel(y)
		result = -c / s
	}

	if neg {
		return -result
	}
	return result
}
