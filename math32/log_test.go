package math32

import (
	"math"
	"testing"
)

func TestLog(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64
	}{
		// Special values
		{"one", 1, 0}, // log(1) = 0
		{"e", float32(E), 1},

		// Powers of 2 (should be exact with our method)
		{"2", 2, math.Log(2)},
		{"4", 4, math.Log(4)},
		{"8", 8, math.Log(8)},
		{"0.5", 0.5, math.Log(0.5)},
		{"0.25", 0.25, math.Log(0.25)},

		// Small values
		{"1.1", 1.1, math.Log(1.1)},
		{"1.5", 1.5, math.Log(1.5)},
		{"2.5", 2.5, math.Log(2.5)},

		// Medium values
		{"10", 10.0, math.Log(10.0)},
		{"100", 100.0, math.Log(100.0)},
		{"1000", 1000.0, math.Log(1000.0)},

		// Large values
		{"1e10", 1e10, math.Log(1e10)},
		{"1e20", 1e20, math.Log(1e20)},
		{"1e38", 1e38, math.Log(1e38)},

		// Small positive values
		{"0.1", 0.1, math.Log(0.1)},
		{"0.01", 0.01, math.Log(0.01)},
		{"1e-10", 1e-10, math.Log(1e-10)},
		{"1e-20", 1e-20, math.Log(1e-20)},
		{"1e-38", 1e-38, math.Log(1e-38)},

		// Values near 1
		{"1.01", 1.01, math.Log(1.01)},
		{"0.99", 0.99, math.Log(0.99)},
		{"1.001", 1.001, math.Log(1.001)},
		{"0.999", 0.999, math.Log(0.999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log(tt.input)
			expected := float32(tt.expected)

			// Use relative error except near zero
			// Log polynomial has slightly larger errors than trig functions
			tolerance := float32(5e-3) // 0.5% relative error
			var err float32
			if Abs(expected) > 1e-6 {
				err = Abs(got-expected) / Abs(expected)
			} else {
				err = Abs(got - expected)
			}

			if err > tolerance {
				t.Errorf("Log(%v) = %v, want %v (rel error: %e)",
					tt.input, got, expected, err)
			}
		})
	}
}

func TestLogSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Log(NaN()); !IsNaN(nan) {
		t.Errorf("Log(NaN) = %v, want NaN", nan)
	}

	// Test negative (should give NaN)
	if nan := Log(-1); !IsNaN(nan) {
		t.Errorf("Log(-1) = %v, want NaN", nan)
	}

	if nan := Log(-100); !IsNaN(nan) {
		t.Errorf("Log(-100) = %v, want NaN", nan)
	}

	// Test zero (should give -Inf)
	if inf := Log(0); !IsInf(inf, -1) {
		t.Errorf("Log(0) = %v, want -Inf", inf)
	}

	// Test +Inf
	if inf := Log(float32(math.Inf(1))); !IsInf(inf, 1) {
		t.Errorf("Log(+Inf) = %v, want +Inf", inf)
	}

	// Test log(1) = 0 exactly
	if result := Log(1); result != 0 {
		t.Errorf("Log(1) = %v, want exactly 0", result)
	}
}

func TestLogExpIdentity(t *testing.T) {
	// Test that log(exp(x)) = x
	testValues := []float32{-5, -2, -1, -0.5, 0, 0.5, 1, 2, 5, 10}

	for _, x := range testValues {
		expX := Exp(x)
		result := Log(expX)

		diff := Abs(result - x)
		if diff > 3e-3 {
			t.Errorf("log(exp(%v)) = %v, want %v (diff: %e)", x, result, x, diff)
		}
	}
}

func TestExpLogIdentity(t *testing.T) {
	// Test that exp(log(x)) = x
	testValues := []float32{0.1, 0.5, 1, 2, 5, 10, 100, 1000, 1e10}

	for _, x := range testValues {
		logX := Log(x)
		result := Exp(logX)

		relError := Abs(result-x) / x
		// Composed operations accumulate errors
		if relError > 2e-3 {
			t.Errorf("exp(log(%v)) = %v, want %v (rel error: %e)",
				x, result, x, relError)
		}
	}
}

func TestLogAccuracy(t *testing.T) {
	var maxRelError float32
	var maxErrorAt float32

	// Test a range of values across many orders of magnitude
	testPoints := []float32{
		1e-38, 1e-20, 1e-10, 1e-5, 0.001, 0.01, 0.1,
		0.5, 0.9, 0.99, 1.0, 1.01, 1.1, 1.5, 2,
		5, 10, 100, 1000, 1e10, 1e20, 1e38,
	}

	for _, x := range testPoints {
		got := Log(x)
		expected := float32(math.Log(float64(x)))

		// Use relative error for non-zero expected values
		var relError float32
		if Abs(expected) > 1e-6 {
			relError = Abs(got-expected) / Abs(expected)
		} else {
			relError = Abs(got - expected)
		}

		if relError > maxRelError {
			maxRelError = relError
			maxErrorAt = x
		}
	}

	t.Logf("Maximum relative error: %e at x=%v", maxRelError, maxErrorAt)

	// Log has slightly larger errors than other functions due to
	// the transformation and polynomial approximation
	if maxRelError > 5e-3 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxRelError, maxErrorAt)
	}
}

func TestLogMonotonicity(t *testing.T) {
	// Log should be strictly increasing for positive values
	prev := Log(1e-38)
	testValues := []float32{
		1e-38, 1e-20, 1e-10, 0.001, 0.01, 0.1, 0.5,
		1, 2, 5, 10, 100, 1000, 1e10, 1e20,
	}

	for _, x := range testValues {
		curr := Log(x)
		if curr < prev {
			t.Errorf("Log not monotonic: Log(%v)=%v < previous=%v", x, curr, prev)
		}
		prev = curr
	}
}

func TestLogPowersOf2(t *testing.T) {
	// log(2^k) should be exactly k·ln(2)
	for k := -20; k <= 20; k++ {
		var x float32
		absK := k
		if k < 0 {
			absK = -k
		}
		if k >= 0 {
			x = float32(int64(1) << uint(absK))
		} else {
			x = 1.0 / float32(int64(1)<<uint(absK))
		}

		got := Log(x)
		expected := float32(k) * ln2Hi

		// Allow small error for the ln2Lo component
		diff := Abs(got - expected)
		if diff > 3e-5 {
			t.Errorf("Log(2^%d) = %v, want %v·ln(2) = %v (diff: %e)",
				k, got, k, expected, diff)
		}
	}
}

// Benchmark Log
func BenchmarkLogSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Log(x)
	}
	_ = result
}

func BenchmarkLogNearOne(b *testing.B) {
	x := float32(1.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Log(x)
	}
	_ = result
}

func BenchmarkLogMedium(b *testing.B) {
	x := float32(100.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Log(x)
	}
	_ = result
}

func BenchmarkLogLarge(b *testing.B) {
	x := float32(1e20)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Log(x)
	}
	_ = result
}

func BenchmarkLogFloat64(b *testing.B) {
	x := 100.0
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Log(x)
	}
	_ = result
}
