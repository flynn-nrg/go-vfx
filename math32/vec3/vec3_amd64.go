//go:build amd64 && goexperiment.simd

package vec3

import (
	"simd/archsimd"
)

// load packs a Vec3Impl into a Float32x4 (4th element is 0).
func load(v Vec3Impl) archsimd.Float32x4 {
	arr := [4]float32{v.X, v.Y, v.Z, 0}
	return archsimd.LoadFloat32x4(&arr)
}

// store extracts a Vec3Impl from a Float32x4.
func store(v archsimd.Float32x4) Vec3Impl {
	return Vec3Impl{
		X: v.GetElem(0),
		Y: v.GetElem(1),
		Z: v.GetElem(2),
	}
}

// addTwo returns the sum of two vectors using SIMD.
func addTwo(v1, v2 Vec3Impl) Vec3Impl {
	return store(load(v1).Add(load(v2)))
}

// subTwo returns the difference of two vectors using SIMD.
func subTwo(v1, v2 Vec3Impl) Vec3Impl {
	return store(load(v1).Sub(load(v2)))
}

// mulTwo returns the element-wise product of two vectors using SIMD.
func mulTwo(v1, v2 Vec3Impl) Vec3Impl {
	return store(load(v1).Mul(load(v2)))
}

// divTwo returns the element-wise division of two vectors using SIMD.
func divTwo(v1, v2 Vec3Impl) Vec3Impl {
	return store(load(v1).Div(load(v2)))
}

// scalarMulSIMD returns the scalar multiplication using SIMD.
func scalarMulSIMD(v Vec3Impl, t float32) Vec3Impl {
	vv := load(v)
	tv := archsimd.BroadcastFloat32x4(t)
	return store(vv.Mul(tv))
}

// scalarDivSIMD returns the scalar division using SIMD.
func scalarDivSIMD(v Vec3Impl, t float32) Vec3Impl {
	vv := load(v)
	tv := archsimd.BroadcastFloat32x4(t)
	return store(vv.Div(tv))
}

// dotSIMD computes dot product using SIMD multiply and horizontal add.
func dotSIMD(v1, v2 Vec3Impl) float32 {
	a := load(v1)
	b := load(v2)
	prod := a.Mul(b)
	// Horizontal add: sum all elements
	// prod = [x1*x2, y1*y2, z1*z2, 0]
	// We need x1*x2 + y1*y2 + z1*z2
	return prod.GetElem(0) + prod.GetElem(1) + prod.GetElem(2)
}

// crossSIMD computes cross product using SIMD.
// cross(a,b) = [a.y*b.z - a.z*b.y, a.z*b.x - a.x*b.z, a.x*b.y - a.y*b.x]
func crossSIMD(v1, v2 Vec3Impl) Vec3Impl {
	// For cross product we need:
	// result.x = a.y*b.z - a.z*b.y
	// result.y = a.z*b.x - a.x*b.z
	// result.z = a.x*b.y - a.y*b.x

	// Create shuffled versions and compute using SIMD:
	// a_yzx = [y1, z1, x1, 0], b_zxy = [z2, x2, y2, 0]
	// a_zxy = [z1, x1, y1, 0], b_yzx = [y2, z2, x2, 0]
	a_yzx := archsimd.LoadFloat32x4(&[4]float32{v1.Y, v1.Z, v1.X, 0})
	a_zxy := archsimd.LoadFloat32x4(&[4]float32{v1.Z, v1.X, v1.Y, 0})
	b_zxy := archsimd.LoadFloat32x4(&[4]float32{v2.Z, v2.X, v2.Y, 0})
	b_yzx := archsimd.LoadFloat32x4(&[4]float32{v2.Y, v2.Z, v2.X, 0})

	// cross = a_yzx * b_zxy - a_zxy * b_yzx
	return store(a_yzx.Mul(b_zxy).Sub(a_zxy.Mul(b_yzx)))
}

// lerpSIMD performs linear interpolation using SIMD FMA.
// lerp(v0, v1, t) = v0 + t*(v1 - v0) = v0*(1-t) + v1*t
func lerpSIMD(v0, v1 Vec3Impl, t float32) Vec3Impl {
	a := load(v0)
	b := load(v1)
	tv := archsimd.BroadcastFloat32x4(t)
	oneMinusT := archsimd.BroadcastFloat32x4(1 - t)

	// result = a*(1-t) + b*t
	// Using MulAdd: a*oneMinusT + b*t = MulAdd(b, t, a*(1-t))
	aScaled := a.Mul(oneMinusT)
	result := b.MulAdd(tv, aScaled)
	return store(result)
}

// squaredLengthSIMD returns the squared length using SIMD.
func squaredLengthSIMD(v Vec3Impl) float32 {
	return dotSIMD(v, v)
}

// lengthSIMD returns the length using SIMD.
func lengthSIMD(v Vec3Impl) float32 {
	vv := load(v)
	prod := vv.Mul(vv)
	sumSquares := prod.GetElem(0) + prod.GetElem(1) + prod.GetElem(2)
	// Use SIMD sqrt
	sqrtVec := archsimd.BroadcastFloat32x4(sumSquares).Sqrt()
	return sqrtVec.GetElem(0)
}

// unitVectorSIMD returns a normalized vector using SIMD.
func unitVectorSIMD(v Vec3Impl) Vec3Impl {
	vv := load(v)
	prod := vv.Mul(vv)
	sumSquares := prod.GetElem(0) + prod.GetElem(1) + prod.GetElem(2)

	// Compute 1/sqrt(sumSquares) using reciprocal sqrt
	rsqrt := archsimd.BroadcastFloat32x4(sumSquares).ReciprocalSqrt()
	invLen := rsqrt.GetElem(0)

	return store(vv.Mul(archsimd.BroadcastFloat32x4(invLen)))
}
