package math32

import (
	"math"
	"testing"
)

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float32
		expected float32
	}{
		// Basic comparisons
		{"1 vs 2", 1, 2, 2},
		{"2 vs 1", 2, 1, 2},
		{"equal", 5, 5, 5},

		// Negative numbers
		{"-1 vs -2", -1, -2, -1},
		{"-2 vs -1", -2, -1, -1},
		{"1 vs -1", 1, -1, 1},
		{"-1 vs 1", -1, 1, 1},

		// Zero cases
		{"0 vs 1", 0, 1, 1},
		{"1 vs 0", 1, 0, 1},
		{"0 vs 0", 0, 0, 0},
		{"-1 vs 0", -1, 0, 0},

		// Fractional values
		{"1.5 vs 2.5", 1.5, 2.5, 2.5},
		{"2.5 vs 1.5", 2.5, 1.5, 2.5},
		{"0.1 vs 0.2", 0.1, 0.2, 0.2},

		// Large values
		{"1e10 vs 1e20", 1e10, 1e20, 1e20},
		{"1e20 vs 1e10", 1e20, 1e10, 1e20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.x, tt.y)
			if got != tt.expected {
				t.Errorf("Max(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float32
		expected float32
	}{
		// Basic comparisons
		{"1 vs 2", 1, 2, 1},
		{"2 vs 1", 2, 1, 1},
		{"equal", 5, 5, 5},

		// Negative numbers
		{"-1 vs -2", -1, -2, -2},
		{"-2 vs -1", -2, -1, -2},
		{"1 vs -1", 1, -1, -1},
		{"-1 vs 1", -1, 1, -1},

		// Zero cases
		{"0 vs 1", 0, 1, 0},
		{"1 vs 0", 1, 0, 0},
		{"0 vs 0", 0, 0, 0},
		{"-1 vs 0", -1, 0, -1},

		// Fractional values
		{"1.5 vs 2.5", 1.5, 2.5, 1.5},
		{"2.5 vs 1.5", 2.5, 1.5, 1.5},
		{"0.1 vs 0.2", 0.1, 0.2, 0.1},

		// Large values
		{"1e10 vs 1e20", 1e10, 1e20, 1e10},
		{"1e20 vs 1e10", 1e20, 1e10, 1e10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Min(tt.x, tt.y)
			if got != tt.expected {
				t.Errorf("Min(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestMaxSpecialCases(t *testing.T) {
	nan := NaN()
	inf := float32(math.Inf(1))
	negInf := float32(math.Inf(-1))

	// NaN propagation
	if result := Max(nan, 1); !IsNaN(result) {
		t.Errorf("Max(NaN, 1) = %v, want NaN", result)
	}
	if result := Max(1, nan); !IsNaN(result) {
		t.Errorf("Max(1, NaN) = %v, want NaN", result)
	}
	if result := Max(nan, nan); !IsNaN(result) {
		t.Errorf("Max(NaN, NaN) = %v, want NaN", result)
	}

	// Infinity cases
	if result := Max(inf, 1); !IsInf(result, 1) {
		t.Errorf("Max(+Inf, 1) = %v, want +Inf", result)
	}
	if result := Max(1, inf); !IsInf(result, 1) {
		t.Errorf("Max(1, +Inf) = %v, want +Inf", result)
	}
	if result := Max(negInf, 1); result != 1 {
		t.Errorf("Max(-Inf, 1) = %v, want 1", result)
	}

	// Zero sign cases
	posZero := float32(0)
	negZero := float32(math.Copysign(0, -1))

	result := Max(posZero, negZero)
	if result != 0 || Signbit(result) {
		t.Errorf("Max(+0, -0) = %v (signbit=%v), want +0", result, Signbit(result))
	}

	result = Max(negZero, posZero)
	if result != 0 || Signbit(result) {
		t.Errorf("Max(-0, +0) = %v (signbit=%v), want +0", result, Signbit(result))
	}
}

func TestMinSpecialCases(t *testing.T) {
	nan := NaN()
	inf := float32(math.Inf(1))
	negInf := float32(math.Inf(-1))

	// NaN propagation
	if result := Min(nan, 1); !IsNaN(result) {
		t.Errorf("Min(NaN, 1) = %v, want NaN", result)
	}
	if result := Min(1, nan); !IsNaN(result) {
		t.Errorf("Min(1, NaN) = %v, want NaN", result)
	}
	if result := Min(nan, nan); !IsNaN(result) {
		t.Errorf("Min(NaN, NaN) = %v, want NaN", result)
	}

	// Infinity cases
	if result := Min(inf, 1); result != 1 {
		t.Errorf("Min(+Inf, 1) = %v, want 1", result)
	}
	if result := Min(1, inf); result != 1 {
		t.Errorf("Min(1, +Inf) = %v, want 1", result)
	}
	if result := Min(negInf, 1); !IsInf(result, -1) {
		t.Errorf("Min(-Inf, 1) = %v, want -Inf", result)
	}

	// Zero sign cases
	posZero := float32(0)
	negZero := float32(math.Copysign(0, -1))

	result := Min(posZero, negZero)
	if result != 0 || !Signbit(result) {
		t.Errorf("Min(+0, -0) = %v (signbit=%v), want -0", result, Signbit(result))
	}

	result = Min(negZero, posZero)
	if result != 0 || !Signbit(result) {
		t.Errorf("Min(-0, +0) = %v (signbit=%v), want -0", result, Signbit(result))
	}
}

func TestMaxMinCompatibility(t *testing.T) {
	testPairs := []struct {
		x, y float32
	}{
		{1, 2},
		{2, 1},
		{-1, -2},
		{1, -1},
		{0, 1},
		{0, -1},
		{0.5, 1.5},
		{1e10, 1e20},
	}

	for _, tp := range testPairs {
		// Test Max
		ourMax := Max(tp.x, tp.y)
		stdMax := float32(math.Max(float64(tp.x), float64(tp.y)))
		if ourMax != stdMax {
			t.Errorf("Max(%v, %v): ours=%v, math.Max=%v", tp.x, tp.y, ourMax, stdMax)
		}

		// Test Min
		ourMin := Min(tp.x, tp.y)
		stdMin := float32(math.Min(float64(tp.x), float64(tp.y)))
		if ourMin != stdMin {
			t.Errorf("Min(%v, %v): ours=%v, math.Min=%v", tp.x, tp.y, ourMin, stdMin)
		}
	}
}

func TestMaxMinRelationship(t *testing.T) {
	// For any x, y: Min(x,y) <= Max(x,y)
	testPairs := []struct {
		x, y float32
	}{
		{1, 2}, {2, 1}, {-1, -2}, {1, -1},
		{0.5, 1.5}, {1e10, 1e20}, {-5, 10},
	}

	for _, tp := range testPairs {
		minVal := Min(tp.x, tp.y)
		maxVal := Max(tp.x, tp.y)

		if minVal > maxVal {
			t.Errorf("Min(%v, %v)=%v > Max(%v, %v)=%v",
				tp.x, tp.y, minVal, tp.x, tp.y, maxVal)
		}
	}
}

// Benchmark Max
func BenchmarkMax(b *testing.B) {
	x := float32(1.5)
	y := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Max(x, y)
	}
	_ = result
}

func BenchmarkMaxFloat64(b *testing.B) {
	x := 1.5
	y := 2.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Max(x, y)
	}
	_ = result
}

// Benchmark Min
func BenchmarkMin(b *testing.B) {
	x := float32(1.5)
	y := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Min(x, y)
	}
	_ = result
}

func BenchmarkMinFloat64(b *testing.B) {
	x := 1.5
	y := 2.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Min(x, y)
	}
	_ = result
}
