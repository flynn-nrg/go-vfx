package math32

import (
	"math"
	"testing"
)

func TestAsin(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64 // use float64 for reference
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), 0},
		{"one", 1, math.Pi / 2},
		{"negative one", -1, -math.Pi / 2},

		// Small values (fast path: |x| <= 0.5)
		{"small positive", 0.1, math.Asin(0.1)},
		{"small negative", -0.1, math.Asin(-0.1)},
		{"quarter", 0.25, math.Asin(0.25)},
		{"half", 0.5, math.Asin(0.5)},
		{"negative half", -0.5, math.Asin(-0.5)},

		// Large values (sqrt path: 0.5 < |x| <= 1)
		{"0.7", 0.7, math.Asin(0.7)},
		{"0.9", 0.9, math.Asin(0.9)},
		{"0.99", 0.99, math.Asin(0.99)},
		{"0.999", 0.999, math.Asin(0.999)},
		{"negative 0.7", -0.7, math.Asin(-0.7)},
		{"negative 0.99", -0.99, math.Asin(-0.99)},

		// Edge cases near boundaries
		{"sqrt(2)/2", float32(math.Sqrt(2) / 2), math.Asin(math.Sqrt(2) / 2)},
		{"sqrt(3)/2", float32(math.Sqrt(3) / 2), math.Asin(math.Sqrt(3) / 2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asin(tt.input)
			expected := float32(tt.expected)

			// Allow for small floating point errors
			// For float32, we expect ~1 ULP accuracy
			tolerance := float32(1e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Asin(%v) = %v, want %v (diff: %v)", tt.input, got, expected, diff)
			}
		})
	}
}

func TestAsinSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Asin(float32(math.NaN())); !math.IsNaN(float64(nan)) {
		t.Errorf("Asin(NaN) = %v, want NaN", nan)
	}

	// Test out of range
	if nan := Asin(1.1); !math.IsNaN(float64(nan)) {
		t.Errorf("Asin(1.1) = %v, want NaN", nan)
	}

	if nan := Asin(-1.1); !math.IsNaN(float64(nan)) {
		t.Errorf("Asin(-1.1) = %v, want NaN", nan)
	}
}

func TestAsinAccuracy(t *testing.T) {
	// Test a range of values and measure maximum error
	var maxError float32
	var maxErrorAt float32

	steps := 1000
	for i := 0; i <= steps; i++ {
		x := -1.0 + 2.0*float32(i)/float32(steps)
		got := Asin(x)
		expected := float32(math.Asin(float64(x)))

		err := Abs(got - expected)
		if err > maxError {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Maximum error: %e at x=%v", maxError, maxErrorAt)

	// For float32, we expect errors to be very small (< 1e-6)
	if maxError > 1e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

// Benchmark the fast path (|x| <= 0.5)
func BenchmarkAsinSmall(b *testing.B) {
	x := float32(0.25)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Asin(x)
	}
	_ = result
}

// Benchmark the sqrt path (0.5 < |x| <= 1)
func BenchmarkAsinLarge(b *testing.B) {
	x := float32(0.9)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Asin(x)
	}
	_ = result
}

// Benchmark against float64 math.Asin for comparison
func BenchmarkAsinFloat64(b *testing.B) {
	x := 0.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Asin(x)
	}
	_ = result
}

func TestAcos(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64 // use float64 for reference
	}{
		// Special cases
		{"one", 1, 0},
		{"negative one", -1, math.Pi},
		{"zero", 0, math.Pi / 2},

		// Small values (|x| <= 0.5, fast path using asin)
		{"small positive", 0.1, math.Acos(0.1)},
		{"small negative", -0.1, math.Acos(-0.1)},
		{"quarter", 0.25, math.Acos(0.25)},
		{"negative quarter", -0.25, math.Acos(-0.25)},
		{"half", 0.5, math.Acos(0.5)},
		{"negative half", -0.5, math.Acos(-0.5)},

		// Large positive values (x > 0.5, direct computation)
		{"0.7", 0.7, math.Acos(0.7)},
		{"0.9", 0.9, math.Acos(0.9)},
		{"0.99", 0.99, math.Acos(0.99)},
		{"0.999", 0.999, math.Acos(0.999)},
		{"0.9999", 0.9999, math.Acos(0.9999)},

		// Large negative values (x < -0.5, direct computation)
		{"negative 0.7", -0.7, math.Acos(-0.7)},
		{"negative 0.9", -0.9, math.Acos(-0.9)},
		{"negative 0.99", -0.99, math.Acos(-0.99)},
		{"negative 0.999", -0.999, math.Acos(-0.999)},

		// Mathematical constants
		{"sqrt(2)/2", float32(math.Sqrt(2) / 2), math.Acos(math.Sqrt(2) / 2)},
		{"sqrt(3)/2", float32(math.Sqrt(3) / 2), math.Acos(math.Sqrt(3) / 2)},
		{"1/2", 0.5, math.Pi / 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acos(tt.input)
			expected := float32(tt.expected)

			// Allow for small floating point errors
			// Slightly higher tolerance than Asin due to multiple paths
			tolerance := float32(2e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Acos(%v) = %v, want %v (diff: %v)", tt.input, got, expected, diff)
			}
		})
	}
}

func TestAcosSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Acos(float32(math.NaN())); !math.IsNaN(float64(nan)) {
		t.Errorf("Acos(NaN) = %v, want NaN", nan)
	}

	// Test out of range
	if nan := Acos(1.1); !math.IsNaN(float64(nan)) {
		t.Errorf("Acos(1.1) = %v, want NaN", nan)
	}

	if nan := Acos(-1.1); !math.IsNaN(float64(nan)) {
		t.Errorf("Acos(-1.1) = %v, want NaN", nan)
	}

	// Test exact values
	if result := Acos(1); result != 0 {
		t.Errorf("Acos(1) = %v, want 0", result)
	}

	if result := Acos(-1); Abs(result-float32(Pi)) > 1e-6 {
		t.Errorf("Acos(-1) = %v, want π = %v", result, Pi)
	}

	if result := Acos(0); Abs(result-float32(Pi/2)) > 1e-6 {
		t.Errorf("Acos(0) = %v, want π/2 = %v", result, Pi/2)
	}
}

func TestAcosAsinRelationship(t *testing.T) {
	// Test that acos(x) + asin(x) = π/2 for all valid x
	testValues := []float32{-0.99, -0.9, -0.7, -0.5, -0.25, 0, 0.25, 0.5, 0.7, 0.9, 0.99}

	for _, x := range testValues {
		acos := Acos(x)
		asin := Asin(x)
		sum := acos + asin
		expected := float32(Pi / 2)
		diff := Abs(sum - expected)

		if diff > 1e-6 {
			t.Errorf("Acos(%v) + Asin(%v) = %v, want π/2 = %v (diff: %e)",
				x, x, sum, expected, diff)
		}
	}
}

func TestAcosAccuracy(t *testing.T) {
	// Test a range of values and measure maximum error
	var maxError float32
	var maxErrorAt float32

	steps := 1000
	for i := 0; i <= steps; i++ {
		x := -1.0 + 2.0*float32(i)/float32(steps)
		got := Acos(x)
		expected := float32(math.Acos(float64(x)))

		err := Abs(got - expected)
		if err > maxError {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Maximum error: %e at x=%v", maxError, maxErrorAt)

	// For float32, we expect errors to be very small
	// The multiple computation paths can accumulate slight errors
	if maxError > 2e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

func TestAcosMonotonicity(t *testing.T) {
	// Acos should be monotonically decreasing (larger input = smaller output)
	prev := float32(Pi)
	for i := 0; i <= 1000; i++ {
		x := -1.0 + 2.0*float32(i)/1000.0
		curr := Acos(x)
		if curr > prev {
			t.Errorf("Acos not monotonic: Acos(%v)=%v > Acos(%v)=%v",
				x, curr, x-0.002, prev)
		}
		prev = curr
	}
}

// Benchmark the fast path (|x| <= 0.5)
func BenchmarkAcosSmall(b *testing.B) {
	x := float32(0.25)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Acos(x)
	}
	_ = result
}

// Benchmark large positive (x > 0.5)
func BenchmarkAcosLargePositive(b *testing.B) {
	x := float32(0.9)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Acos(x)
	}
	_ = result
}

// Benchmark large negative (x < -0.5)
func BenchmarkAcosLargeNegative(b *testing.B) {
	x := float32(-0.9)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Acos(x)
	}
	_ = result
}

// Benchmark against float64 math.Acos for comparison
func BenchmarkAcosFloat64(b *testing.B) {
	x := 0.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Acos(x)
	}
	_ = result
}
