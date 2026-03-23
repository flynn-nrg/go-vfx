package vec3

import (
	"testing"
)

var (
	v1 = Vec3Impl{X: 1.0, Y: 2.0, Z: 3.0}
	v2 = Vec3Impl{X: 4.0, Y: 5.0, Z: 6.0}
	v3 = Vec3Impl{X: 0.5, Y: 0.5, Z: 0.5}
)

func BenchmarkAdd(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Add(v1, v2)
	}
	_ = result
}

func BenchmarkSub(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Sub(v1, v2)
	}
	_ = result
}

func BenchmarkMul(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Mul(v1, v2)
	}
	_ = result
}

func BenchmarkDiv(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Div(v1, v2)
	}
	_ = result
}

func BenchmarkScalarMul(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = ScalarMul(v1, 2.5)
	}
	_ = result
}

func BenchmarkScalarDiv(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = ScalarDiv(v1, 2.5)
	}
	_ = result
}

func BenchmarkDot(b *testing.B) {
	var result float32
	for i := 0; i < b.N; i++ {
		result = Dot(v1, v2)
	}
	_ = result
}

func BenchmarkCross(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Cross(v1, v2)
	}
	_ = result
}

func BenchmarkLength(b *testing.B) {
	var result float32
	for i := 0; i < b.N; i++ {
		result = v1.Length()
	}
	_ = result
}

func BenchmarkSquaredLength(b *testing.B) {
	var result float32
	for i := 0; i < b.N; i++ {
		result = v1.SquaredLength()
	}
	_ = result
}

func BenchmarkUnitVector(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = UnitVector(v1)
	}
	_ = result
}

func BenchmarkLerp(b *testing.B) {
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		result = Lerp(v1, v2, 0.5)
	}
	_ = result
}

// Compound operations typical in path tracing
func BenchmarkReflect(b *testing.B) {
	normal := UnitVector(v3)
	var result Vec3Impl
	for i := 0; i < b.N; i++ {
		// reflect = v - 2*dot(v,n)*n
		d := Dot(v1, normal)
		result = Sub(v1, ScalarMul(normal, 2*d))
	}
	_ = result
}

func BenchmarkNormalizeAndDot(b *testing.B) {
	var result float32
	for i := 0; i < b.N; i++ {
		n1 := UnitVector(v1)
		n2 := UnitVector(v2)
		result = Dot(n1, n2)
	}
	_ = result
}

// Test correctness
func TestDot(t *testing.T) {
	result := Dot(v1, v2)
	expected := float32(1*4 + 2*5 + 3*6) // 4 + 10 + 18 = 32
	if result != expected {
		t.Errorf("Dot(%v, %v) = %v, want %v", v1, v2, result, expected)
	}
}

func TestCross(t *testing.T) {
	a := Vec3Impl{X: 1, Y: 0, Z: 0}
	b := Vec3Impl{X: 0, Y: 1, Z: 0}
	result := Cross(a, b)
	expected := Vec3Impl{X: 0, Y: 0, Z: 1}
	if result != expected {
		t.Errorf("Cross(%v, %v) = %v, want %v", a, b, result, expected)
	}
}

func TestAdd(t *testing.T) {
	result := Add(v1, v2)
	expected := Vec3Impl{X: 5, Y: 7, Z: 9}
	if result != expected {
		t.Errorf("Add(%v, %v) = %v, want %v", v1, v2, result, expected)
	}
}

func TestSub(t *testing.T) {
	result := Sub(v1, v2)
	expected := Vec3Impl{X: -3, Y: -3, Z: -3}
	if result != expected {
		t.Errorf("Sub(%v, %v) = %v, want %v", v1, v2, result, expected)
	}
}

func TestScalarMul(t *testing.T) {
	result := ScalarMul(v1, 2)
	expected := Vec3Impl{X: 2, Y: 4, Z: 6}
	if result != expected {
		t.Errorf("ScalarMul(%v, 2) = %v, want %v", v1, result, expected)
	}
}

func TestLerp(t *testing.T) {
	a := Vec3Impl{X: 0, Y: 0, Z: 0}
	b := Vec3Impl{X: 10, Y: 20, Z: 30}
	result := Lerp(a, b, 0.5)
	expected := Vec3Impl{X: 5, Y: 10, Z: 15}
	if result != expected {
		t.Errorf("Lerp(%v, %v, 0.5) = %v, want %v", a, b, result, expected)
	}
}

func TestUnitVector(t *testing.T) {
	v := Vec3Impl{X: 3, Y: 0, Z: 0}
	result := UnitVector(v)
	// Should be {1, 0, 0}
	if result.X < 0.999 || result.X > 1.001 {
		t.Errorf("UnitVector(%v).X = %v, want ~1", v, result.X)
	}
	if result.Y != 0 || result.Z != 0 {
		t.Errorf("UnitVector(%v) = %v, want {1, 0, 0}", v, result)
	}
}
