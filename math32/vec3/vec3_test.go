package vec3

import (
	"math"
	"testing"
)

const epsilon = 1e-6

func almostEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < epsilon
}

func TestDot(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3Impl
		v2       *Vec3Impl
		expected float32
	}{
		{
			name:     "Unit vectors parallel",
			v1:       &Vec3Impl{X: 1, Y: 0, Z: 0},
			v2:       &Vec3Impl{X: 1, Y: 0, Z: 0},
			expected: 1,
		},
		{
			name:     "Unit vectors perpendicular",
			v1:       &Vec3Impl{X: 1, Y: 0, Z: 0},
			v2:       &Vec3Impl{X: 0, Y: 1, Z: 0},
			expected: 0,
		},
		{
			name:     "General case",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 4, Y: 5, Z: 6},
			expected: 32, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
		},
		{
			name:     "Negative values",
			v1:       &Vec3Impl{X: -1, Y: 2, Z: -3},
			v2:       &Vec3Impl{X: 4, Y: -5, Z: 6},
			expected: -32, // -1*4 + 2*-5 + -3*6 = -4 - 10 - 18 = -32
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Dot(tt.v1, tt.v2)
			if !almostEqual(result, tt.expected) {
				t.Errorf("Dot() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3Impl
		v2       *Vec3Impl
		expected *Vec3Impl
	}{
		{
			name:     "Simple addition",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 4, Y: 5, Z: 6},
			expected: &Vec3Impl{X: 5, Y: 7, Z: 9},
		},
		{
			name:     "Adding zero vector",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 0, Y: 0, Z: 0},
			expected: &Vec3Impl{X: 1, Y: 2, Z: 3},
		},
		{
			name:     "Negative values",
			v1:       &Vec3Impl{X: -1, Y: 2, Z: -3},
			v2:       &Vec3Impl{X: 1, Y: -2, Z: 3},
			expected: &Vec3Impl{X: 0, Y: 0, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.v1, tt.v2)
			if !almostEqual(result.X, tt.expected.X) ||
				!almostEqual(result.Y, tt.expected.Y) ||
				!almostEqual(result.Z, tt.expected.Z) {
				t.Errorf("Add() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3Impl
		v2       *Vec3Impl
		expected *Vec3Impl
	}{
		{
			name:     "Simple subtraction",
			v1:       &Vec3Impl{X: 5, Y: 7, Z: 9},
			v2:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			expected: &Vec3Impl{X: 4, Y: 5, Z: 6},
		},
		{
			name:     "Subtracting zero vector",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 0, Y: 0, Z: 0},
			expected: &Vec3Impl{X: 1, Y: 2, Z: 3},
		},
		{
			name:     "Result with negative values",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 4, Y: 5, Z: 6},
			expected: &Vec3Impl{X: -3, Y: -3, Z: -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sub(tt.v1, tt.v2)
			if !almostEqual(result.X, tt.expected.X) ||
				!almostEqual(result.Y, tt.expected.Y) ||
				!almostEqual(result.Z, tt.expected.Z) {
				t.Errorf("Sub() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3Impl
		v2       *Vec3Impl
		expected *Vec3Impl
	}{
		{
			name:     "Simple multiplication",
			v1:       &Vec3Impl{X: 2, Y: 3, Z: 4},
			v2:       &Vec3Impl{X: 5, Y: 6, Z: 7},
			expected: &Vec3Impl{X: 10, Y: 18, Z: 28},
		},
		{
			name:     "Multiply by one",
			v1:       &Vec3Impl{X: 2, Y: 3, Z: 4},
			v2:       &Vec3Impl{X: 1, Y: 1, Z: 1},
			expected: &Vec3Impl{X: 2, Y: 3, Z: 4},
		},
		{
			name:     "Multiply by zero",
			v1:       &Vec3Impl{X: 2, Y: 3, Z: 4},
			v2:       &Vec3Impl{X: 0, Y: 0, Z: 0},
			expected: &Vec3Impl{X: 0, Y: 0, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mul(tt.v1, tt.v2)
			if !almostEqual(result.X, tt.expected.X) ||
				!almostEqual(result.Y, tt.expected.Y) ||
				!almostEqual(result.Z, tt.expected.Z) {
				t.Errorf("Mul() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestScalarMul(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3Impl
		scalar   float32
		expected *Vec3Impl
	}{
		{
			name:     "Multiply by 2",
			v:        &Vec3Impl{X: 1, Y: 2, Z: 3},
			scalar:   2,
			expected: &Vec3Impl{X: 2, Y: 4, Z: 6},
		},
		{
			name:     "Multiply by 0.5",
			v:        &Vec3Impl{X: 2, Y: 4, Z: 6},
			scalar:   0.5,
			expected: &Vec3Impl{X: 1, Y: 2, Z: 3},
		},
		{
			name:     "Multiply by -1",
			v:        &Vec3Impl{X: 1, Y: 2, Z: 3},
			scalar:   -1,
			expected: &Vec3Impl{X: -1, Y: -2, Z: -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScalarMul(tt.v, tt.scalar)
			if !almostEqual(result.X, tt.expected.X) ||
				!almostEqual(result.Y, tt.expected.Y) ||
				!almostEqual(result.Z, tt.expected.Z) {
				t.Errorf("ScalarMul() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestCross(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3Impl
		v2       *Vec3Impl
		expected *Vec3Impl
	}{
		{
			name:     "X cross Y = Z",
			v1:       &Vec3Impl{X: 1, Y: 0, Z: 0},
			v2:       &Vec3Impl{X: 0, Y: 1, Z: 0},
			expected: &Vec3Impl{X: 0, Y: 0, Z: 1},
		},
		{
			name:     "Y cross Z = X",
			v1:       &Vec3Impl{X: 0, Y: 1, Z: 0},
			v2:       &Vec3Impl{X: 0, Y: 0, Z: 1},
			expected: &Vec3Impl{X: 1, Y: 0, Z: 0},
		},
		{
			name:     "Z cross X = Y",
			v1:       &Vec3Impl{X: 0, Y: 0, Z: 1},
			v2:       &Vec3Impl{X: 1, Y: 0, Z: 0},
			expected: &Vec3Impl{X: 0, Y: 1, Z: 0},
		},
		{
			name:     "General case",
			v1:       &Vec3Impl{X: 1, Y: 2, Z: 3},
			v2:       &Vec3Impl{X: 4, Y: 5, Z: 6},
			expected: &Vec3Impl{X: -3, Y: 6, Z: -3}, // (2*6-3*5, -(1*6-3*4), 1*5-2*4)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Cross(tt.v1, tt.v2)
			if !almostEqual(result.X, tt.expected.X) ||
				!almostEqual(result.Y, tt.expected.Y) ||
				!almostEqual(result.Z, tt.expected.Z) {
				t.Errorf("Cross() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestLength(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3Impl
		expected float32
	}{
		{
			name:     "Unit vector X",
			v:        &Vec3Impl{X: 1, Y: 0, Z: 0},
			expected: 1,
		},
		{
			name:     "3-4-5 triangle",
			v:        &Vec3Impl{X: 0, Y: 3, Z: 4},
			expected: 5,
		},
		{
			name:     "General case",
			v:        &Vec3Impl{X: 1, Y: 2, Z: 2},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Length()
			if !almostEqual(result, tt.expected) {
				t.Errorf("Length() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance comparison

func BenchmarkDot(b *testing.B) {
	v1 := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	v2 := &Vec3Impl{X: 4.2, Y: 5.1, Z: 6.8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Dot(v1, v2)
	}
}

func BenchmarkAdd(b *testing.B) {
	v1 := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	v2 := &Vec3Impl{X: 4.2, Y: 5.1, Z: 6.8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Add(v1, v2)
	}
}

func BenchmarkSub(b *testing.B) {
	v1 := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	v2 := &Vec3Impl{X: 4.2, Y: 5.1, Z: 6.8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sub(v1, v2)
	}
}

func BenchmarkMul(b *testing.B) {
	v1 := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	v2 := &Vec3Impl{X: 4.2, Y: 5.1, Z: 6.8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Mul(v1, v2)
	}
}

func BenchmarkScalarMul(b *testing.B) {
	v := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	scalar := float32(2.5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ScalarMul(v, scalar)
	}
}

func BenchmarkCross(b *testing.B) {
	v1 := &Vec3Impl{X: 1.5, Y: 2.3, Z: 3.7}
	v2 := &Vec3Impl{X: 4.2, Y: 5.1, Z: 6.8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Cross(v1, v2)
	}
}
