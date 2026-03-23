//go:build amd64 && goexperiment.simd

package math32

import "simd/archsimd"

// max uses the MAXSS instruction via archsimd intrinsics.
func max(x, y float32) float32 {
	vx := archsimd.BroadcastFloat32x4(x)
	vy := archsimd.BroadcastFloat32x4(y)
	return vx.Max(vy).GetElem(0)
}

// min uses the MINSS instruction via archsimd intrinsics.
func min(x, y float32) float32 {
	vx := archsimd.BroadcastFloat32x4(x)
	vy := archsimd.BroadcastFloat32x4(y)
	return vx.Min(vy).GetElem(0)
}
