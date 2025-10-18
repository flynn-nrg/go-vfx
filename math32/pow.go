package math32

// Pow returns x**y, the base-x exponential of y.
//
// Special cases are (in order):
//
//	Pow(x, ±0) = 1 for any x
//	Pow(1, y) = 1 for any y
//	Pow(x, 1) = x for any x
//	Pow(NaN, y) = NaN
//	Pow(x, NaN) = NaN
//	Pow(±0, y) = ±Inf for y an odd integer < 0
//	Pow(±0, -Inf) = +Inf
//	Pow(±0, +Inf) = +0
//	Pow(±0, y) = +Inf for finite y < 0 and not an odd integer
//	Pow(±0, y) = ±0 for y an odd integer > 0
//	Pow(±0, y) = +0 for finite y > 0 and not an odd integer
//	Pow(-1, ±Inf) = 1
//	Pow(x, +Inf) = +Inf for |x| > 1
//	Pow(x, -Inf) = +Inf for |x| < 1
//	Pow(x, +Inf) = +0 for |x| < 1
//	Pow(x, -Inf) = +0 for |x| > 1
//	Pow(+Inf, y) = +Inf for y > 0
//	Pow(+Inf, y) = +0 for y < 0
//	Pow(-Inf, y) = Pow(-0, -y)
//	Pow(x, y) = NaN for finite x < 0 and finite non-integer y
func Pow(x, y float32) float32 {
	// Handle special cases

	// Pow(x, ±0) = 1 for any x (even NaN or Inf)
	if y == 0 {
		return 1
	}

	// Pow(1, y) = 1 for any y (even NaN or Inf)
	if x == 1 {
		return 1
	}

	// Pow(x, 1) = x
	if y == 1 {
		return x
	}

	// NaN propagation
	if IsNaN(x) || IsNaN(y) {
		return NaN()
	}

	// Handle x = 0
	if x == 0 {
		if y < 0 {
			// Pow(±0, y) for y < 0
			if isOddInteger(y) {
				// Return infinity with sign of x
				return Copysign(Inf(1), x)
			}
			return Inf(1)
		}
		// Pow(±0, y) for y > 0
		if isOddInteger(y) {
			// Return zero with sign of x
			return Copysign(0, x)
		}
		return 0
	}

	// Handle infinities
	if IsInf(y, 0) {
		if IsInf(y, 1) {
			// y = +Inf
			if x == -1 {
				return 1
			}
			absx := Abs(x)
			if absx > 1 {
				return Inf(1)
			}
			return 0
		}
		// y = -Inf
		if x == -1 {
			return 1
		}
		absx := Abs(x)
		if absx < 1 {
			return Inf(1)
		}
		return 0
	}

	if IsInf(x, 0) {
		if IsInf(x, 1) {
			// x = +Inf
			if y > 0 {
				return Inf(1)
			}
			return 0
		}
		// x = -Inf
		// Pow(-Inf, y) = Pow(-0, -y)
		return Pow(Copysign(0, -1), -y)
	}

	// Handle negative base
	if x < 0 {
		// Check if y is an integer
		if !isInteger(y) {
			// Pow(x, y) for x < 0 and non-integer y is undefined
			return NaN()
		}

		// For integer y, compute abs and fix sign later
		absResult := pow(Abs(x), y)

		// Determine sign of result
		if isOddInteger(y) {
			return -absResult
		}
		return absResult
	}

	// General case: x > 0, finite values
	return pow(x, y)
}

// pow is the internal implementation for positive x using exp(y * log(x))
func pow(x, y float32) float32 {
	// Handle exact cases for better performance
	if y == 2 {
		return x * x
	}
	if y == 0.5 {
		return Sqrt(x)
	}
	if y == -1 {
		return 1 / x
	}

	// Check for potential overflow/underflow before computing
	logx := Log(x)
	product := y * logx

	// Check bounds before calling Exp
	// These thresholds match those in exp.go
	const (
		expOverflow  = 88.72283905206835
		expUnderflow = -103.97207708399179
	)

	if product > expOverflow {
		return Inf(1)
	}
	if product < expUnderflow {
		return 0
	}

	// General case: x^y = exp(y * log(x))
	return Exp(product)
}

// isInteger returns true if x is an integer value
func isInteger(x float32) bool {
	if IsNaN(x) || IsInf(x, 0) {
		return false
	}
	return Floor(x) == x
}

// isOddInteger returns true if x is an odd integer
func isOddInteger(x float32) bool {
	if !isInteger(x) {
		return false
	}

	// For large values, all representable integers are even
	// (float32 loses precision for integers beyond 2^24)
	if Abs(x) > 16777216 { // 2^24
		return false
	}

	xi := int32(x)
	return xi&1 == 1
}
