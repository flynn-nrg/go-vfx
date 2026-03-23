//go:build amd64 && goexperiment.simd

package math32

import "simd/archsimd"

// floor uses the ROUNDSS instruction via archsimd intrinsics.
func floor(x float32) float32 {
	v := archsimd.BroadcastFloat32x4(x)
	return v.Floor().GetElem(0)
}

// ceil uses the ROUNDSS instruction via archsimd intrinsics.
func ceil(x float32) float32 {
	v := archsimd.BroadcastFloat32x4(x)
	return v.Ceil().GetElem(0)
}

// round implements "round half away from zero" by adding 0.5 to the absolute
// value and truncating, then restoring the sign.
func round(x float32) float32 {
	v := archsimd.BroadcastFloat32x4(x)
	half := archsimd.BroadcastFloat32x4(0.5)

	// 0x7FFFFFFF clears the sign bit (for absolute value)
	absMask := archsimd.BroadcastInt32x4(0x7FFFFFFF)
	// 0x80000000 isolates the sign bit (math.MinInt32 = -2147483648)
	signMask := archsimd.BroadcastInt32x4(-0x80000000)

	// Get absolute value: v & 0x7FFFFFFF
	abs := v.AsInt32x4().And(absMask).AsFloat32x4()

	// Add 0.5 and truncate toward zero
	result := abs.Add(half).Trunc()

	// Extract sign bit from original: v & 0x80000000
	signBit := v.AsInt32x4().And(signMask)

	// Apply sign to result: result | signBit
	result = result.AsInt32x4().Or(signBit).AsFloat32x4()

	return result.GetElem(0)
}
