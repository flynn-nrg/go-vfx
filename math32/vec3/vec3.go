// Package vec3 provides utility functions to work with vectors.
package vec3

import (
	"math"

	"github.com/flynn-nrg/go-vfx/math32"
	"github.com/flynn-nrg/go-vfx/math32/fastrandom"
)

// Vec3Impl defines a vector with its position and colour.
type Vec3Impl struct {
	X float32
	Y float32
	Z float32
	R float32
	G float32
	B float32
}

// Length returns the length of this vector.
func (v *Vec3Impl) Length() float32 {
	return math32.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
}

// SquaredLength returns the squared length of this vector.
func (v *Vec3Impl) SquaredLength() float32 {
	return (v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z)
}

// MakeUnitVector transform the vector into its unit representation.
func (v *Vec3Impl) MakeUnitVector() {
	l := v.Length()
	v.X = v.X / l
	v.Y = v.Y / l
	v.Z = v.Z / l
}

// Add returns the sum of two or more vectors.
func Add(v1 *Vec3Impl, args ...*Vec3Impl) *Vec3Impl {
	sum := &Vec3Impl{
		X: v1.X,
		Y: v1.Y,
		Z: v1.Z,
	}

	for i := range args {
		sum.X, sum.Y, sum.Z = add3(sum.X, sum.Y, sum.Z, args[i].X, args[i].Y, args[i].Z)
	}

	return sum
}

// Sub returns the subtraction of two or more vectors.
func Sub(v1 *Vec3Impl, args ...*Vec3Impl) *Vec3Impl {
	res := &Vec3Impl{
		X: v1.X,
		Y: v1.Y,
		Z: v1.Z,
	}

	for i := range args {
		res.X, res.Y, res.Z = sub3(res.X, res.Y, res.Z, args[i].X, args[i].Y, args[i].Z)
	}

	return res
}

// Mul returns the multiplication of two vectors.
func Mul(v1 *Vec3Impl, v2 *Vec3Impl) *Vec3Impl {
	x, y, z := mul3(v1.X, v1.Y, v1.Z, v2.X, v2.Y, v2.Z)
	return &Vec3Impl{
		X: x,
		Y: y,
		Z: z,
	}
}

// Div returns the division of two vectors.
func Div(v1 *Vec3Impl, v2 *Vec3Impl) *Vec3Impl {
	return &Vec3Impl{
		X: v1.X / v2.X,
		Y: v1.Y / v2.Y,
		Z: v1.Z / v2.Z,
	}
}

// ScalarMul returns the scalar multiplication of the given vector and scalar values.
func ScalarMul(v1 *Vec3Impl, t float32) *Vec3Impl {
	x, y, z := scalarMul3(v1.X, v1.Y, v1.Z, t)
	return &Vec3Impl{
		X: x,
		Y: y,
		Z: z,
	}
}

// ScalarMul returns the scalar division of the given vector and scalar values.
func ScalarDiv(v1 *Vec3Impl, t float32) *Vec3Impl {
	return &Vec3Impl{
		X: v1.X / t,
		Y: v1.Y / t,
		Z: v1.Z / t,
	}
}

// Dot computes the dot product of the two supplied vectors.
func Dot(v1 *Vec3Impl, v2 *Vec3Impl) float32 {
	return dot(v1.X, v1.Y, v1.Z, v2.X, v2.Y, v2.Z)
}

// Cross computes the cross product of the two supplied vectors.
func Cross(v1 *Vec3Impl, v2 *Vec3Impl) *Vec3Impl {
	x, y, z := cross(v1.X, v1.Y, v1.Z, v2.X, v2.Y, v2.Z)
	return &Vec3Impl{
		X: x,
		Y: y,
		Z: z,
	}
}

// UnitVector returns a unit vector representation of the supplied vector.
func UnitVector(v *Vec3Impl) *Vec3Impl {
	return ScalarDiv(v, v.Length())
}

// RandomCosineDirection returns a vector with a random cosine direction.
func RandomCosineDirection(random *fastrandom.XorShift) *Vec3Impl {
	r1 := random.Float32()
	r2 := random.Float32()
	z := math32.Sqrt(1 - r2)
	phi := 2 * math.Pi * r1
	x := math32.Cos(phi) * 2 * math32.Sqrt(r2)
	y := math32.Sin(phi) * 2 * math32.Sqrt(r2)
	return &Vec3Impl{X: x, Y: y, Z: z}
}

// RandomToSphere returns a new random sphere of the given radius at the given distance.
func RandomToSphere(radius float32, distanceSquared float32, random *fastrandom.XorShift) *Vec3Impl {
	r1 := random.Float32()
	r2 := random.Float32()
	z := 1 + r2*(math32.Sqrt(1-radius*radius/distanceSquared)-1)
	phi := 2 * math32.Pi * r1
	x := math32.Cos(phi) * math32.Sqrt(1-z*z)
	y := math32.Sin(phi) * math32.Sqrt(1-z*z)
	return &Vec3Impl{X: x, Y: y, Z: z}
}

// DeNAN ensures that the vector elements are numbers.
func DeNAN(v *Vec3Impl) *Vec3Impl {
	x := v.X
	y := v.Y
	z := v.Z
	if math32.IsNaN(x) || math32.IsInf(x, -1) || math32.IsInf(x, 1) {
		x = 0
	}

	if math32.IsNaN(y) || math32.IsInf(y, -1) || math32.IsInf(y, 1) {
		y = 0
	}

	if math32.IsNaN(z) || math32.IsInf(z, -1) || math32.IsInf(z, 1) {
		z = 0
	}

	return &Vec3Impl{X: x, Y: y, Z: z}
}

// Min3 returns a new vector with the minimum coordinates among the supplied ones.
func Min3(v0 *Vec3Impl, v1 *Vec3Impl, v2 *Vec3Impl) *Vec3Impl {
	xMin := float32(math32.MaxFloat32)
	yMin := float32(math32.MaxFloat32)
	zMin := float32(math32.MaxFloat32)

	if v0.X < xMin {
		xMin = v0.X
	}

	if v1.X < xMin {
		xMin = v1.X
	}

	if v2.X < xMin {
		xMin = v2.X
	}

	if v0.Y < yMin {
		yMin = v0.Y
	}

	if v1.Y < yMin {
		yMin = v1.Y
	}

	if v2.Y < yMin {
		yMin = v2.Y
	}

	if v0.Z < zMin {
		zMin = v0.Z
	}

	if v1.Z < zMin {
		zMin = v1.Z
	}

	if v2.Z < zMin {
		zMin = v2.Z
	}

	return &Vec3Impl{X: xMin, Y: yMin, Z: zMin}
}

// Max3 returns a new vector with the maximum coordinates among the supplied ones.
func Max3(v0 *Vec3Impl, v1 *Vec3Impl, v2 *Vec3Impl) *Vec3Impl {
	xMax := -float32(math32.MaxFloat32)
	yMax := -float32(math32.MaxFloat32)
	zMax := -float32(math32.MaxFloat32)

	if v0.X > xMax {
		xMax = v0.X
	}

	if v1.X > xMax {
		xMax = v1.X
	}

	if v2.X > xMax {
		xMax = v2.X
	}

	if v0.Y > yMax {
		yMax = v0.Y
	}

	if v1.Y > yMax {
		yMax = v1.Y
	}

	if v2.Y > yMax {
		yMax = v2.Y
	}

	if v0.Z > zMax {
		zMax = v0.Z
	}

	if v1.Z > zMax {
		zMax = v1.Z
	}

	if v2.Z > zMax {
		zMax = v2.Z
	}

	return &Vec3Impl{X: xMax, Y: yMax, Z: zMax}
}

// Lerp performs a linear interpolation between the two provided vectors.
func Lerp(v0, v1 *Vec3Impl, t float32) *Vec3Impl {
	return &Vec3Impl{
		X: (1-t)*v0.X + t*v1.X,
		Y: (1-t)*v0.Y + t*v1.Y,
		Z: (1-t)*v0.Z + t*v1.Z,
	}
}

// Equals returns whether two vectors are the same.
func Equals(v0, v1 *Vec3Impl) bool {
	return v0.X == v1.X &&
		v0.Y == v1.Y &&
		v0.Z == v1.Z
}
