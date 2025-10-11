package math32

import "math"

// IsNaN reports whether x is an IEEE 754 "not-a-number" value.
//
// This function uses the IEEE 754 property that NaN is the only value
// that is not equal to itself. This compiles to a very efficient
// floating-point comparison on all architectures.
func IsNaN(x float32) bool {
	// IEEE 754 says that only NaNs satisfy x != x
	return x != x
}

// IsInf reports whether x is an infinity, according to sign.
// If sign > 0, IsInf reports whether x is positive infinity.
// If sign < 0, IsInf reports whether x is negative infinity.
// If sign == 0, IsInf reports whether x is either infinity.
func IsInf(x float32, sign int) bool {
	// Get the bit representation
	bits := math.Float32bits(x)

	// For float32:
	// - Sign bit: bit 31
	// - Exponent: bits 23-30 (8 bits)
	// - Mantissa: bits 0-22 (23 bits)
	//
	// Infinity has:
	// - Exponent = 0xFF (all 1s)
	// - Mantissa = 0

	// Mask for exponent and mantissa (bits 0-30, excluding sign bit)
	const expMantissaMask = 0x7FFFFFFF
	// Pattern for infinity (exponent all 1s, mantissa all 0s)
	const infPattern = 0x7F800000

	// Check if it's any infinity (ignore sign bit)
	if bits&expMantissaMask != infPattern {
		return false
	}

	// If sign == 0, any infinity is OK
	if sign == 0 {
		return true
	}

	// Check the sign bit (bit 31)
	// If sign bit is 1, it's negative infinity
	// If sign bit is 0, it's positive infinity
	isNegative := bits&0x80000000 != 0

	if sign > 0 {
		return !isNegative // positive infinity
	}
	return isNegative // negative infinity
}

// Signbit reports whether x is negative or negative zero.
//
// This function examines the sign bit directly, making it very fast.
// It returns true even for negative zero and negative NaN.
func Signbit(x float32) bool {
	// Get the bit representation
	bits := math.Float32bits(x)

	// Check bit 31 (the sign bit)
	// If set (1), the number is negative
	return bits&0x80000000 != 0
}

// NaN returns a float32 "not-a-number" value.
func NaN() float32 {
	return math.Float32frombits(0x7FC00000) // Standard quiet NaN
}
