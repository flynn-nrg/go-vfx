package math32

import (
	"math"
	"testing"
)

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64 // use float64 for reference
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), 0},
		{"one", 1, 1},
		{"four", 4, 2},
		{"nine", 9, 3},
		{"sixteen", 16, 4},

		// Fractional values
		{"0.25", 0.25, 0.5},
		{"0.5", 0.5, math.Sqrt(0.5)},
		{"0.75", 0.75, math.Sqrt(0.75)},

		// Small values
		{"1e-10", 1e-10, math.Sqrt(1e-10)},
		{"1e-20", 1e-20, math.Sqrt(1e-20)},

		// Large values
		{"1e10", 1e10, math.Sqrt(1e10)},
		{"1e20", 1e20, math.Sqrt(1e20)},

		// Mathematical constants
		{"2", 2, math.Sqrt(2)},
		{"3", 3, math.Sqrt(3)},
		{"5", 5, math.Sqrt(5)},
		{"10", 10, math.Sqrt(10)},

		// Perfect squares
		{"100", 100, 10},
		{"10000", 10000, 100},

		// Non-perfect squares
		{"7", 7, math.Sqrt(7)},
		{"13", 13, math.Sqrt(13)},
		{"42", 42, math.Sqrt(42)},
		{"1234.5", 1234.5, math.Sqrt(1234.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.input)
			expected := float32(tt.expected)

			// For float32, we expect very high accuracy (essentially perfect)
			// Hardware sqrt should be correctly rounded
			tolerance := float32(1e-6)
			diff := Abs(got - expected)
			relError := diff / expected

			if diff > tolerance && relError > 1e-6 {
				t.Errorf("Sqrt(%v) = %v, want %v (diff: %v, rel error: %e)",
					tt.input, got, expected, diff, relError)
			}
		})
	}
}

func TestSqrtSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Sqrt(float32(math.NaN())); !math.IsNaN(float64(nan)) {
		t.Errorf("Sqrt(NaN) = %v, want NaN", nan)
	}

	// Test negative numbers (should return NaN)
	if nan := Sqrt(-1); !math.IsNaN(float64(nan)) {
		t.Errorf("Sqrt(-1) = %v, want NaN", nan)
	}

	if nan := Sqrt(-100); !math.IsNaN(float64(nan)) {
		t.Errorf("Sqrt(-100) = %v, want NaN", nan)
	}

	// Test positive infinity
	if inf := Sqrt(float32(math.Inf(1))); !math.IsInf(float64(inf), 1) {
		t.Errorf("Sqrt(+Inf) = %v, want +Inf", inf)
	}

	// Test that sqrt(x)^2 ≈ x for various values
	testValues := []float32{0.1, 0.5, 1, 2, 3.5, 10, 100, 1000}
	for _, x := range testValues {
		sqrtX := Sqrt(x)
		reconstructed := sqrtX * sqrtX
		diff := Abs(reconstructed - x)
		relError := diff / x

		if relError > 1e-6 {
			t.Errorf("Sqrt(%v)^2 = %v, want %v (rel error: %e)",
				x, reconstructed, x, relError)
		}
	}
}

func TestSqrtAccuracy(t *testing.T) {
	// Test a range of values and measure maximum error
	var maxError float32
	var maxErrorAt float32

	testPoints := []float32{
		0, 1e-38, 1e-20, 1e-10, 1e-5, 0.001, 0.01, 0.1, 0.5,
		1, 2, 3, 5, 10, 100, 1000, 10000, 1e10, 1e20, 1e38,
	}

	for _, x := range testPoints {
		if x < 0 {
			continue
		}
		got := Sqrt(x)
		expected := float32(math.Sqrt(float64(x)))

		err := Abs(got - expected)
		if x > 0 {
			err = err / expected // relative error
		}

		if err > maxError && !math.IsInf(float64(x), 0) {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Maximum relative error: %e at x=%v", maxError, maxErrorAt)

	// Hardware sqrt should be very accurate (< 1 ULP for float32)
	if maxError > 1e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

func TestSqrtMonotonicity(t *testing.T) {
	// Sqrt should be monotonically increasing
	prev := float32(0)
	for i := 0; i < 1000; i++ {
		x := float32(i) / 10.0
		curr := Sqrt(x)
		if curr < prev {
			t.Errorf("Sqrt not monotonic: Sqrt(%v)=%v < Sqrt(%v)=%v",
				x, curr, x-0.1, prev)
		}
		prev = curr
	}
}

// Benchmark hardware-accelerated Sqrt
func BenchmarkSqrt(b *testing.B) {
	x := float32(2.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sqrt(x)
	}
	_ = result
}

// Benchmark with various input ranges
func BenchmarkSqrtSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sqrt(x)
	}
	_ = result
}

func BenchmarkSqrtLarge(b *testing.B) {
	x := float32(1e10)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sqrt(x)
	}
	_ = result
}

// Benchmark float64 sqrt for comparison
func BenchmarkSqrtFloat64(b *testing.B) {
	x := 2.0
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Sqrt(x)
	}
	_ = result
}

// Benchmark the old method (convert to float64 and back)
func BenchmarkSqrtViaFloat64(b *testing.B) {
	x := float32(2.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = float32(math.Sqrt(float64(x)))
	}
	_ = result
}
