//go:build amd64 && goexperiment.simd

package math32

import "simd/archsimd"

// sqrt uses the SQRTSS instruction via archsimd intrinsics.
func sqrt(x float32) float32 {
	v := archsimd.BroadcastFloat32x4(x)
	return v.Sqrt().GetElem(0)
}
