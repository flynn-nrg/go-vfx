package math32_test

import (
	"math"
	"testing"

	"github.com/flynn-nrg/go-vfx/math32"
	"github.com/flynn-nrg/go-vfx/math32/fastrandom"
	"github.com/flynn-nrg/go-vfx/math32/mat3"
	"github.com/flynn-nrg/go-vfx/math32/vec3"
)

// Float32 implementation using math32 package

type Ray32 struct {
	Origin    vec3.Vec3Impl
	Direction vec3.Vec3Impl
}

// reflect32 reflects a vector v around normal n
func reflect32(v, n vec3.Vec3Impl) vec3.Vec3Impl {
	// r = v - 2*dot(v,n)*n
	dot := vec3.Dot(v, n)
	return vec3.Sub(v, vec3.ScalarMul(n, 2*dot))
}

// refract32 refracts a vector through a surface
func refract32(v, n vec3.Vec3Impl, niOverNt float32) (vec3.Vec3Impl, bool) {
	uv := vec3.UnitVector(v)
	dt := vec3.Dot(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		// refracted = niOverNt * (uv - n*dt) - n*sqrt(discriminant)
		term1 := vec3.ScalarMul(vec3.Sub(uv, vec3.ScalarMul(n, dt)), niOverNt)
		term2 := vec3.ScalarMul(n, math32.Sqrt(discriminant))
		return vec3.Sub(term1, term2), true
	}
	return vec3.Vec3Impl{}, false
}

// schlick32 computes Schlick's approximation for Fresnel reflectance
func schlick32(cosine, refIdx float32) float32 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math32.Pow(1-cosine, 5)
}

// scatterDielectric32 simulates light scattering through a dielectric material (glass, water, etc.)
func scatterDielectric32(rayDir, normal vec3.Vec3Impl, refIdx float32, random *fastrandom.XorShift) (vec3.Vec3Impl, bool) {
	var outwardNormal vec3.Vec3Impl
	var niOverNt float32
	var cosine float32

	dotDN := vec3.Dot(rayDir, normal)

	if dotDN > 0 {
		outwardNormal = vec3.ScalarMul(normal, -1)
		niOverNt = refIdx
		cosine = refIdx * dotDN / vec3.UnitVector(rayDir).Length()
	} else {
		outwardNormal = normal
		niOverNt = 1.0 / refIdx
		cosine = -dotDN / vec3.UnitVector(rayDir).Length()
	}

	var reflectProb float32
	refracted, canRefract := refract32(rayDir, outwardNormal, niOverNt)

	if canRefract {
		reflectProb = schlick32(cosine, refIdx)
	} else {
		reflectProb = 1.0
	}

	if random.Float32() < reflectProb {
		reflected := reflect32(rayDir, normal)
		return reflected, true
	}

	return refracted, true
}

// pathTrace32 simulates a complete path tracing pass with multiple bounces
func pathTrace32(ray *Ray32, maxDepth int, random *fastrandom.XorShift) vec3.Vec3Impl {
	attenuation := vec3.Vec3Impl{X: 1.0, Y: 1.0, Z: 1.0}
	currentRay := ray

	for depth := 0; depth < maxDepth; depth++ {
		// Simulate ray-sphere intersection
		// Use some arbitrary values to simulate hitting a dielectric sphere
		t := 1.5 + random.Float32()*0.5
		hitPoint := vec3.Add(currentRay.Origin, vec3.ScalarMul(currentRay.Direction, t))

		// Calculate normal (for a sphere centered at origin)
		center := vec3.Vec3Impl{X: 0, Y: 0, Z: 0}
		normal := vec3.UnitVector(vec3.Sub(hitPoint, center))

		// Scatter through dielectric material
		refIdx := float32(1.5) // glass
		scattered, didScatter := scatterDielectric32(currentRay.Direction, normal, refIdx, random)

		if !didScatter {
			return vec3.Vec3Impl{X: 0, Y: 0, Z: 0}
		}

		// Apply some attenuation (wavelength-dependent for spectral rendering)
		wavelength := 400.0 + random.Float32()*300.0 // 400-700nm visible spectrum
		attenuationFactor := 0.95 + 0.05*math32.Sin(wavelength*0.01)
		attenuation = vec3.ScalarMul(attenuation, attenuationFactor)

		// Transform scattered direction using TBN matrix (common in normal mapping)
		tangent := vec3.Vec3Impl{X: normal.Y, Y: -normal.X, Z: 0}
		if tangent.SquaredLength() < 0.01 {
			tangent = vec3.Vec3Impl{X: 1, Y: 0, Z: 0}
		}
		tangent = vec3.UnitVector(tangent)
		bitangent := vec3.Cross(normal, tangent)
		tbn := mat3.NewTBN(tangent, bitangent, normal)

		// Apply TBN transformation
		scattered = mat3.MatrixVectorMul(tbn, scattered)
		scattered = vec3.UnitVector(scattered)

		// Prepare for next bounce
		currentRay = &Ray32{
			Origin:    hitPoint,
			Direction: scattered,
		}

		// Russian roulette for path termination
		survivalProb := float32(0.9)
		if random.Float32() > survivalProb {
			break
		}
		attenuation = vec3.ScalarDiv(attenuation, survivalProb)
	}

	// Return accumulated color with some computation to prevent optimization
	result := vec3.ScalarMul(attenuation, math32.Abs(math32.Sin(attenuation.X)+math32.Cos(attenuation.Y)))
	return result
}

// Float64 baseline implementation using standard math package

type Ray64 struct {
	Origin    [3]float64
	Direction [3]float64
}

func dot64(v1, v2 [3]float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

func length64(v [3]float64) float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func normalize64(v [3]float64) [3]float64 {
	l := length64(v)
	return [3]float64{v[0] / l, v[1] / l, v[2] / l}
}

func add64(v1, v2 [3]float64) [3]float64 {
	return [3]float64{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

func sub64(v1, v2 [3]float64) [3]float64 {
	return [3]float64{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

func scalarMul64(v [3]float64, s float64) [3]float64 {
	return [3]float64{v[0] * s, v[1] * s, v[2] * s}
}

func cross64(v1, v2 [3]float64) [3]float64 {
	return [3]float64{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}
}

func reflect64(v, n [3]float64) [3]float64 {
	d := dot64(v, n)
	return sub64(v, scalarMul64(n, 2*d))
}

func refract64(v, n [3]float64, niOverNt float64) ([3]float64, bool) {
	uv := normalize64(v)
	dt := dot64(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		term1 := scalarMul64(sub64(uv, scalarMul64(n, dt)), niOverNt)
		term2 := scalarMul64(n, math.Sqrt(discriminant))
		return sub64(term1, term2), true
	}
	return [3]float64{}, false
}

func schlick64(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

func scatterDielectric64(rayDir, normal [3]float64, refIdx float64, random *fastrandom.XorShift) ([3]float64, bool) {
	var outwardNormal [3]float64
	var niOverNt float64
	var cosine float64

	dotDN := dot64(rayDir, normal)

	if dotDN > 0 {
		outwardNormal = scalarMul64(normal, -1)
		niOverNt = refIdx
		cosine = refIdx * dotDN / length64(normalize64(rayDir))
	} else {
		outwardNormal = normal
		niOverNt = 1.0 / refIdx
		cosine = -dotDN / length64(normalize64(rayDir))
	}

	var reflectProb float64
	refracted, canRefract := refract64(rayDir, outwardNormal, niOverNt)

	if canRefract {
		reflectProb = schlick64(cosine, refIdx)
	} else {
		reflectProb = 1.0
	}

	if random.Float32() < float32(reflectProb) {
		reflected := reflect64(rayDir, normal)
		return reflected, true
	}

	return refracted, true
}

func matrixVectorMul64(m [9]float64, v [3]float64) [3]float64 {
	return [3]float64{
		m[0]*v[0] + m[1]*v[1] + m[2]*v[2],
		m[3]*v[0] + m[4]*v[1] + m[5]*v[2],
		m[6]*v[0] + m[7]*v[1] + m[8]*v[2],
	}
}

func pathTrace64(ray *Ray64, maxDepth int, random *fastrandom.XorShift) [3]float64 {
	attenuation := [3]float64{1.0, 1.0, 1.0}
	currentRay := ray

	for depth := 0; depth < maxDepth; depth++ {
		// Simulate ray-sphere intersection
		t := 1.5 + float64(random.Float32())*0.5
		hitPoint := add64(currentRay.Origin, scalarMul64(currentRay.Direction, t))

		// Calculate normal
		center := [3]float64{0, 0, 0}
		normal := normalize64(sub64(hitPoint, center))

		// Scatter through dielectric material
		refIdx := 1.5
		scattered, didScatter := scatterDielectric64(currentRay.Direction, normal, refIdx, random)

		if !didScatter {
			return [3]float64{0, 0, 0}
		}

		// Apply attenuation
		wavelength := 400.0 + float64(random.Float32())*300.0
		attenuationFactor := 0.95 + 0.05*math.Sin(wavelength*0.01)
		attenuation = scalarMul64(attenuation, attenuationFactor)

		// TBN transformation
		tangent := [3]float64{normal[1], -normal[0], 0}
		if dot64(tangent, tangent) < 0.01 {
			tangent = [3]float64{1, 0, 0}
		}
		tangent = normalize64(tangent)
		bitangent := cross64(normal, tangent)
		tbn := [9]float64{
			tangent[0], bitangent[0], normal[0],
			tangent[1], bitangent[1], normal[1],
			tangent[2], bitangent[2], normal[2],
		}

		scattered = matrixVectorMul64(tbn, scattered)
		scattered = normalize64(scattered)

		currentRay = &Ray64{
			Origin:    hitPoint,
			Direction: scattered,
		}

		// Russian roulette
		survivalProb := 0.9
		if random.Float32() > float32(survivalProb) {
			break
		}
		attenuation = scalarMul64(attenuation, 1.0/survivalProb)
	}

	result := scalarMul64(attenuation, math.Abs(math.Sin(attenuation[0])+math.Cos(attenuation[1])))
	return result
}

// Benchmarks

func BenchmarkPathTracer32_SingleBounce(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray32{
		Origin:    vec3.Vec3Impl{X: 0, Y: 0, Z: -5},
		Direction: vec3.Vec3Impl{X: 0, Y: 0, Z: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace32(ray, 1, random)
		if result.X < -1000 { // Prevent optimization
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkPathTracer64_SingleBounce(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray64{
		Origin:    [3]float64{0, 0, -5},
		Direction: [3]float64{0, 0, 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace64(ray, 1, random)
		if result[0] < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkPathTracer32_FiveBounces(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray32{
		Origin:    vec3.Vec3Impl{X: 0, Y: 0, Z: -5},
		Direction: vec3.Vec3Impl{X: 0, Y: 0, Z: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace32(ray, 5, random)
		if result.X < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkPathTracer64_FiveBounces(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray64{
		Origin:    [3]float64{0, 0, -5},
		Direction: [3]float64{0, 0, 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace64(ray, 5, random)
		if result[0] < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkPathTracer32_TenBounces(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray32{
		Origin:    vec3.Vec3Impl{X: 0, Y: 0, Z: -5},
		Direction: vec3.Vec3Impl{X: 0, Y: 0, Z: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace32(ray, 10, random)
		if result.X < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkPathTracer64_TenBounces(b *testing.B) {
	random := fastrandom.New(12345)
	ray := &Ray64{
		Origin:    [3]float64{0, 0, -5},
		Direction: [3]float64{0, 0, 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pathTrace64(ray, 10, random)
		if result[0] < -1000 {
			b.Fatal("unexpected")
		}
	}
}

// Individual operation benchmarks for detailed profiling

func BenchmarkDielectricScatter32(b *testing.B) {
	random := fastrandom.New(12345)
	rayDir := vec3.Vec3Impl{X: 0.5, Y: 0.5, Z: 1}
	normal := vec3.Vec3Impl{X: 0, Y: 0, Z: -1}
	refIdx := float32(1.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, ok := scatterDielectric32(rayDir, normal, refIdx, random)
		if !ok || result.X < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkDielectricScatter64(b *testing.B) {
	random := fastrandom.New(12345)
	rayDir := [3]float64{0.5, 0.5, 1}
	normal := [3]float64{0, 0, -1}
	refIdx := 1.5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, ok := scatterDielectric64(rayDir, normal, refIdx, random)
		if !ok || result[0] < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkVectorOps32(b *testing.B) {
	v1 := vec3.Vec3Impl{X: 1, Y: 2, Z: 3}
	v2 := vec3.Vec3Impl{X: 4, Y: 5, Z: 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r1 := vec3.Add(v1, v2)
		r2 := vec3.Sub(v1, v2)
		r3 := vec3.Dot(v1, v2)
		r4 := vec3.Cross(v1, v2)
		r5 := vec3.UnitVector(r1)
		if r5.X+r2.X+r4.X+r3 < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkVectorOps64(b *testing.B) {
	v1 := [3]float64{1, 2, 3}
	v2 := [3]float64{4, 5, 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r1 := add64(v1, v2)
		r2 := sub64(v1, v2)
		r3 := dot64(v1, v2)
		r4 := cross64(v1, v2)
		r5 := normalize64(r1)
		if r5[0]+r2[0]+r4[0]+r3 < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkTrigOps32(b *testing.B) {
	angles := make([]float32, 100)
	for i := range angles {
		angles[i] = float32(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := float32(0)
		for _, angle := range angles {
			sum += math32.Sin(angle) + math32.Cos(angle)
		}
		if sum < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkTrigOps64(b *testing.B) {
	angles := make([]float64, 100)
	for i := range angles {
		angles[i] = float64(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0.0
		for _, angle := range angles {
			sum += math.Sin(angle) + math.Cos(angle)
		}
		if sum < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkSqrtOps32(b *testing.B) {
	values := make([]float32, 100)
	for i := range values {
		values[i] = float32(i) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := float32(0)
		for _, v := range values {
			sum += math32.Sqrt(v)
		}
		if sum < -1000 {
			b.Fatal("unexpected")
		}
	}
}

func BenchmarkSqrtOps64(b *testing.B) {
	values := make([]float64, 100)
	for i := range values {
		values[i] = float64(i) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0.0
		for _, v := range values {
			sum += math.Sqrt(v)
		}
		if sum < -1000 {
			b.Fatal("unexpected")
		}
	}
}
