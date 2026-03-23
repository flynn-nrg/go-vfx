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
	signMask := archsimd.BroadcastFloat32x4(-0.0).AsInt32x4()

	// Get absolute value: clear sign bit using ~signMask & v
	abs := signMask.AndNot(v.AsInt32x4()).AsFloat32x4()

	// Add 0.5 and truncate
	result := abs.Add(half).Trunc()

	// Restore original sign: copy sign bit from x to result
	signBit := v.AsInt32x4().And(signMask)
	result = result.AsInt32x4().Or(signBit).AsFloat32x4()

	return result.GetElem(0)
}
