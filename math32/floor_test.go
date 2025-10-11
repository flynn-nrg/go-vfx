package math32

import (
	"math"
	"testing"
)

func TestFloor(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), float32(math.Copysign(0, -1))},
		{"positive infinity", float32(math.Inf(1)), float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1)), float32(math.Inf(-1))},

		// Positive integers
		{"1", 1, 1},
		{"2", 2, 2},
		{"10", 10, 10},
		{"100", 100, 100},

		// Negative integers
		{"-1", -1, -1},
		{"-2", -2, -2},
		{"-10", -10, -10},
		{"-100", -100, -100},

		// Positive fractions
		{"0.1", 0.1, 0},
		{"0.5", 0.5, 0},
		{"0.9", 0.9, 0},
		{"0.99", 0.99, 0},
		{"1.1", 1.1, 1},
		{"1.5", 1.5, 1},
		{"1.9", 1.9, 1},
		{"2.5", 2.5, 2},
		{"9.99", 9.99, 9},

		// Negative fractions
		{"-0.1", -0.1, -1},
		{"-0.5", -0.5, -1},
		{"-0.9", -0.9, -1},
		{"-0.99", -0.99, -1},
		{"-1.1", -1.1, -2},
		{"-1.5", -1.5, -2},
		{"-1.9", -1.9, -2},
		{"-2.5", -2.5, -3},
		{"-9.99", -9.99, -10},

		// Large values
		{"1e10", 1e10, 1e10},
		{"-1e10", -1e10, -1e10},
		{"1.5e10", 1.5e10, 1.5e10}, // Large enough that +0.5 doesn't change representation

		// Small values
		{"1e-10", 1e-10, 0},
		{"-1e-10", -1e-10, -1},

		// Values near integer boundaries
		{"2.0000001", 2.0000001, 2},
		{"1.9999999", 1.9999999, 1},
		{"-2.0000002", -2.0000002, -3}, // Slightly larger to ensure it's not exactly -2 in float32
		{"-1.9999999", -1.9999999, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Floor(tt.input)

			// For NaN, check if both are NaN
			if math.IsNaN(float64(tt.expected)) {
				if !math.IsNaN(float64(got)) {
					t.Errorf("Floor(%v) = %v, want NaN", tt.input, got)
				}
				return
			}

			if got != tt.expected {
				t.Errorf("Floor(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFloorNaN(t *testing.T) {
	nan := float32(math.NaN())
	result := Floor(nan)
	if !IsNaN(result) {
		t.Errorf("Floor(NaN) = %v, want NaN", result)
	}
}

func TestFloorCompatibility(t *testing.T) {
	// Test that our Floor matches math.Floor behavior
	testValues := []float32{
		0, -0, 1, -1, 0.5, -0.5, 1.5, -1.5,
		2.5, -2.5, 9.99, -9.99, 100.5, -100.5,
		0.1, -0.1, 0.9, -0.9, 1.1, -1.1,
		float32(math.Inf(1)), float32(math.Inf(-1)),
		1e10, -1e10, 1e-10, -1e-10,
	}

	for _, x := range testValues {
		our := Floor(x)
		std := float32(math.Floor(float64(x)))
		if our != std {
			t.Errorf("Floor(%v): ours=%v, math.Floor=%v", x, our, std)
		}
	}
}

func TestFloorPreservesZeroSign(t *testing.T) {
	// Test that Floor(+0) = +0
	posZero := float32(0)
	result := Floor(posZero)
	if result != 0 || Signbit(result) {
		t.Errorf("Floor(+0) = %v (signbit=%v), want +0", result, Signbit(result))
	}

	// Test that Floor(-0) = -0
	negZero := float32(math.Copysign(0, -1))
	result = Floor(negZero)
	if result != 0 || !Signbit(result) {
		t.Errorf("Floor(-0) = %v (signbit=%v), want -0", result, Signbit(result))
	}
}

func TestFloorRange(t *testing.T) {
	// Test a range of values to ensure correctness
	for i := -100; i <= 100; i++ {
		for _, fraction := range []float32{0, 0.1, 0.25, 0.5, 0.75, 0.9, 0.99} {
			x := float32(i) + fraction

			got := Floor(x)
			expected := float32(math.Floor(float64(x)))

			if got != expected {
				t.Errorf("Floor(%v) = %v, want %v", x, got, expected)
			}

			// Also test negative
			x = -x
			got = Floor(x)
			expected = float32(math.Floor(float64(x)))

			if got != expected {
				t.Errorf("Floor(%v) = %v, want %v", x, got, expected)
			}
		}
	}
}

func TestFloorMonotonicity(t *testing.T) {
	// Floor should be monotonically non-decreasing
	prev := Floor(-100)
	for i := -100.0; i <= 100.0; i += 0.1 {
		x := float32(i)
		curr := Floor(x)
		if curr < prev {
			t.Errorf("Floor not monotonic: Floor(%v)=%v < Floor(%v)=%v",
				x, curr, x-0.1, prev)
		}
		prev = curr
	}
}

func TestFloorIdempotent(t *testing.T) {
	// Floor(Floor(x)) should equal Floor(x)
	testValues := []float32{
		-10.5, -2.7, -1.1, -0.5, 0, 0.5, 1.1, 2.7, 10.5,
	}

	for _, x := range testValues {
		once := Floor(x)
		twice := Floor(once)

		if once != twice {
			t.Errorf("Floor not idempotent: Floor(%v)=%v, Floor(Floor(%v))=%v",
				x, once, x, twice)
		}
	}
}

// Benchmark Floor
func BenchmarkFloor(b *testing.B) {
	x := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Floor(x)
	}
	_ = result
}

func BenchmarkFloorNegative(b *testing.B) {
	x := float32(-2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Floor(x)
	}
	_ = result
}

func BenchmarkFloorInteger(b *testing.B) {
	x := float32(5.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Floor(x)
	}
	_ = result
}

func BenchmarkFloorFloat64(b *testing.B) {
	x := 2.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Floor(x)
	}
	_ = result
}

func BenchmarkFloorViaFloat64(b *testing.B) {
	x := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = float32(math.Floor(float64(x)))
	}
	_ = result
}

// Ceil tests

func TestCeil(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), float32(math.Copysign(0, -1))},
		{"positive infinity", float32(math.Inf(1)), float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1)), float32(math.Inf(-1))},

		// Positive integers
		{"1", 1, 1},
		{"2", 2, 2},
		{"10", 10, 10},

		// Negative integers
		{"-1", -1, -1},
		{"-2", -2, -2},
		{"-10", -10, -10},

		// Positive fractions
		{"0.1", 0.1, 1},
		{"0.5", 0.5, 1},
		{"0.9", 0.9, 1},
		{"1.1", 1.1, 2},
		{"1.5", 1.5, 2},
		{"1.9", 1.9, 2},
		{"2.5", 2.5, 3},

		// Negative fractions
		{"-0.1", -0.1, 0},
		{"-0.5", -0.5, 0},
		{"-0.9", -0.9, 0},
		{"-1.1", -1.1, -1},
		{"-1.5", -1.5, -1},
		{"-1.9", -1.9, -1},
		{"-2.5", -2.5, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ceil(tt.input)
			if got != tt.expected {
				t.Errorf("Ceil(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCeilCompatibility(t *testing.T) {
	testValues := []float32{
		0, -0, 1, -1, 0.5, -0.5, 1.5, -1.5,
		2.5, -2.5, 9.99, -9.99,
		float32(math.Inf(1)), float32(math.Inf(-1)),
	}

	for _, x := range testValues {
		our := Ceil(x)
		std := float32(math.Ceil(float64(x)))
		if our != std {
			t.Errorf("Ceil(%v): ours=%v, math.Ceil=%v", x, our, std)
		}
	}
}

// Round tests

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), float32(math.Copysign(0, -1))},
		{"positive infinity", float32(math.Inf(1)), float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1)), float32(math.Inf(-1))},

		// Positive integers
		{"1", 1, 1},
		{"2", 2, 2},
		{"10", 10, 10},

		// Negative integers
		{"-1", -1, -1},
		{"-2", -2, -2},
		{"-10", -10, -10},

		// Positive fractions
		{"0.1", 0.1, 0},
		{"0.4", 0.4, 0},
		{"0.5", 0.5, 1}, // Round half away from zero: rounds to 1
		{"0.6", 0.6, 1},
		{"1.1", 1.1, 1},
		{"1.5", 1.5, 2}, // Round half away from zero: rounds to 2
		{"1.9", 1.9, 2},
		{"2.5", 2.5, 3}, // Round half away from zero: rounds to 3
		{"3.5", 3.5, 4}, // Round half away from zero: rounds to 4

		// Negative fractions
		{"-0.1", -0.1, 0},
		{"-0.4", -0.4, 0},
		{"-0.5", -0.5, -1}, // Round half away from zero: rounds to -1
		{"-0.6", -0.6, -1},
		{"-1.1", -1.1, -1},
		{"-1.5", -1.5, -2}, // Round half away from zero: rounds to -2
		{"-1.9", -1.9, -2},
		{"-2.5", -2.5, -3}, // Round half away from zero: rounds to -3
		{"-3.5", -3.5, -4}, // Round half away from zero: rounds to -4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Round(tt.input)
			if got != tt.expected {
				t.Errorf("Round(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRoundCompatibility(t *testing.T) {
	testValues := []float32{
		0, -0, 1, -1, 0.5, -0.5, 1.5, -1.5,
		2.5, -2.5, 3.5, -3.5,
		float32(math.Inf(1)), float32(math.Inf(-1)),
	}

	for _, x := range testValues {
		our := Round(x)
		std := float32(math.Round(float64(x)))
		if our != std {
			t.Errorf("Round(%v): ours=%v, math.Round=%v", x, our, std)
		}
	}
}

// Benchmarks for Ceil

func BenchmarkCeil(b *testing.B) {
	x := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Ceil(x)
	}
	_ = result
}

func BenchmarkCeilNegative(b *testing.B) {
	x := float32(-2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Ceil(x)
	}
	_ = result
}

func BenchmarkCeilFloat64(b *testing.B) {
	x := 2.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Ceil(x)
	}
	_ = result
}

// Benchmarks for Round

func BenchmarkRound(b *testing.B) {
	x := float32(2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Round(x)
	}
	_ = result
}

func BenchmarkRoundNegative(b *testing.B) {
	x := float32(-2.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Round(x)
	}
	_ = result
}

func BenchmarkRoundFloat64(b *testing.B) {
	x := 2.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Round(x)
	}
	_ = result
}
