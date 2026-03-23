// Package vec3 provides utility functions to work with vectors.
package vec3

import (
	"github.com/flynn-nrg/go-vfx/math32"
	"github.com/flynn-nrg/go-vfx/math32/fastrandom"
)

// Vec3Impl defines a vector with its position and colour.
type Vec3Impl struct {
	X float32
	Y float32
	Z float32
}

// Length returns the length of this vector.
func (v Vec3Impl) Length() float32 {
	return lengthSIMD(v)
}

// SquaredLength returns the squared length of this vector.
func (v Vec3Impl) SquaredLength() float32 {
	return squaredLengthSIMD(v)
}

// MakeUnitVector transform the vector into its unit representation.
func (v Vec3Impl) MakeUnitVector() Vec3Impl {
	return unitVectorSIMD(v)
}

// Add returns the sum of two or more vectors.
func Add(v1 Vec3Impl, args ...Vec3Impl) Vec3Impl {
	sum := v1
	for i := range args {
		sum = addTwo(sum, args[i])
	}
	return sum
}

// Sub returns the subtraction of two or more vectors.
func Sub(v1 Vec3Impl, args ...Vec3Impl) Vec3Impl {
	res := v1
	for i := range args {
		res = subTwo(res, args[i])
	}
	return res
}

// Mul returns the multiplication of two vectors.
func Mul(v1 Vec3Impl, v2 Vec3Impl) Vec3Impl {
	return mulTwo(v1, v2)
}

// Div returns the division of two vectors.
func Div(v1 Vec3Impl, v2 Vec3Impl) Vec3Impl {
	return divTwo(v1, v2)
}

// ScalarMul returns the scalar multiplication of the given vector and scalar values.
func ScalarMul(v1 Vec3Impl, t float32) Vec3Impl {
	return scalarMulSIMD(v1, t)
}

// ScalarDiv returns the scalar division of the given vector and scalar values.
func ScalarDiv(v1 Vec3Impl, t float32) Vec3Impl {
	return scalarDivSIMD(v1, t)
}

// Dot computes the dot product of the two supplied vectors.
func Dot(v1 Vec3Impl, v2 Vec3Impl) float32 {
	return dotSIMD(v1, v2)
}

// Cross computes the cross product of the two supplied vectors.
func Cross(v1 Vec3Impl, v2 Vec3Impl) Vec3Impl {
	return crossSIMD(v1, v2)
}

// UnitVector returns a unit vector representation of the supplied vector.
func UnitVector(v Vec3Impl) Vec3Impl {
	return unitVectorSIMD(v)
}

// RandomCosineDirection returns a vector with a random cosine direction.
func RandomCosineDirection(random *fastrandom.XorShift) Vec3Impl {
	r1 := random.Float32()
	r2 := random.Float32()
	z := math32.Sqrt(1 - r2)
	phi := 2 * math32.Pi * r1
	x := math32.Cos(phi) * 2 * math32.Sqrt(r2)
	y := math32.Sin(phi) * 2 * math32.Sqrt(r2)
	return Vec3Impl{X: x, Y: y, Z: z}
}

// RandomToSphere returns a new random sphere of the given radius at the given distance.
func RandomToSphere(radius float32, distanceSquared float32, random *fastrandom.XorShift) Vec3Impl {
	r1 := random.Float32()
	r2 := random.Float32()
	z := 1 + r2*(math32.Sqrt(1-radius*radius/distanceSquared)-1)
	phi := 2 * math32.Pi * r1
	x := math32.Cos(phi) * math32.Sqrt(1-z*z)
	y := math32.Sin(phi) * math32.Sqrt(1-z*z)
	return Vec3Impl{X: x, Y: y, Z: z}
}

// DeNAN ensures that the vector elements are numbers.
func DeNAN(v Vec3Impl) Vec3Impl {
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

	return Vec3Impl{X: x, Y: y, Z: z}
}

// Min3 returns a new vector with the minimum coordinates among the supplied ones.
func Min3(v0 Vec3Impl, v1 Vec3Impl, v2 Vec3Impl) Vec3Impl {
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

	return Vec3Impl{X: xMin, Y: yMin, Z: zMin}
}

// Max3 returns a new vector with the maximum coordinates among the supplied ones.
func Max3(v0 Vec3Impl, v1 Vec3Impl, v2 Vec3Impl) Vec3Impl {
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

	return Vec3Impl{X: xMax, Y: yMax, Z: zMax}
}

// Lerp performs a linear interpolation between the two provided vectors.
func Lerp(v0, v1 Vec3Impl, t float32) Vec3Impl {
	return lerpSIMD(v0, v1, t)
}

// Equals returns whether two vectors are the same.
func Equals(v0, v1 Vec3Impl) bool {
	return v0.X == v1.X &&
		v0.Y == v1.Y &&
		v0.Z == v1.Z
}
