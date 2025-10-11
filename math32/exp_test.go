package math32

import (
	"math"
	"testing"
)

func TestExp(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64
	}{
		// Special cases
		{"zero", 0, 1},
		{"one", 1, math.E},

		// Small positive values
		{"0.1", 0.1, math.Exp(0.1)},
		{"0.5", 0.5, math.Exp(0.5)},
		{"1.5", 1.5, math.Exp(1.5)},
		{"2.0", 2.0, math.Exp(2.0)},

		// Small negative values
		{"-0.1", -0.1, math.Exp(-0.1)},
		{"-0.5", -0.5, math.Exp(-0.5)},
		{"-1.0", -1.0, math.Exp(-1.0)},
		{"-2.0", -2.0, math.Exp(-2.0)},

		// Medium values
		{"5", 5.0, math.Exp(5.0)},
		{"10", 10.0, math.Exp(10.0)},
		{"-5", -5.0, math.Exp(-5.0)},
		{"-10", -10.0, math.Exp(-10.0)},

		// Large values (near overflow)
		{"80", 80.0, math.Exp(80.0)},
		{"88", 88.0, math.Exp(88.0)},
		{"-80", -80.0, math.Exp(-80.0)},
		{"-100", -100.0, math.Exp(-100.0)},

		// ln(2) and multiples
		{"ln(2)", float32(Ln2), 2.0},
		{"2*ln(2)", float32(2 * Ln2), 4.0},
		{"10*ln(2)", float32(10 * Ln2), 1024.0},
		{"-ln(2)", float32(-Ln2), 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exp(tt.input)
			expected := float32(tt.expected)

			// Relative error for exp (absolute error grows with value)
			tolerance := float32(2e-6)
			relError := Abs(got-expected) / Abs(expected)

			if relError > tolerance && Abs(got-expected) > 1e-6 {
				t.Errorf("Exp(%v) = %v, want %v (rel error: %e)",
					tt.input, got, expected, relError)
			}
		})
	}
}

func TestExpSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Exp(NaN()); !IsNaN(nan) {
		t.Errorf("Exp(NaN) = %v, want NaN", nan)
	}

	// Test +Inf
	if inf := Exp(float32(math.Inf(1))); !IsInf(inf, 1) {
		t.Errorf("Exp(+Inf) = %v, want +Inf", inf)
	}

	// Test -Inf (should give 0)
	if result := Exp(float32(math.Inf(-1))); result != 0 {
		t.Errorf("Exp(-Inf) = %v, want 0", result)
	}

	// Test overflow
	if inf := Exp(1000); !IsInf(inf, 1) {
		t.Errorf("Exp(1000) = %v, want +Inf (overflow)", inf)
	}

	// Test underflow
	if result := Exp(-1000); result != 0 {
		t.Errorf("Exp(-1000) = %v, want 0 (underflow)", result)
	}
}

func TestExpIdentities(t *testing.T) {
	// Test exp(a + b) = exp(a) * exp(b)
	testPairs := []struct {
		a, b float32
	}{
		{1, 2},
		{0.5, 1.5},
		{-1, 3},
		{2, -1},
	}

	for _, tp := range testPairs {
		expSum := Exp(tp.a + tp.b)
		expProd := Exp(tp.a) * Exp(tp.b)

		relError := Abs(expSum-expProd) / Abs(expSum)
		if relError > 1e-5 {
			t.Errorf("exp(%v + %v) = %v, exp(%v)*exp(%v) = %v (rel error: %e)",
				tp.a, tp.b, expSum, tp.a, tp.b, expProd, relError)
		}
	}

	// Test exp(ln(x)) = x
	testValues := []float32{0.5, 1.0, 2.0, 10.0, 100.0}
	for _, x := range testValues {
		lnX := float32(math.Log(float64(x)))
		result := Exp(lnX)

		relError := Abs(result-x) / x
		if relError > 2e-6 {
			t.Errorf("exp(ln(%v)) = %v, want %v (rel error: %e)",
				x, result, x, relError)
		}
	}
}

func TestExpAccuracy(t *testing.T) {
	var maxRelError float32
	var maxErrorAt float32

	// Test range from -20 to 88 (avoid extreme underflow region)
	testPoints := []float32{
		-20, -10, -5, -2, -1, -0.5, -0.1,
		0, 0.1, 0.5, 1, 2, 5, 10, 20, 50, 80, 88,
	}

	for _, x := range testPoints {
		got := Exp(x)
		expected := float32(math.Exp(float64(x)))

		// Use relative error for exp (values vary by many orders of magnitude)
		var relError float32
		if Abs(expected) > 1e-20 {
			relError = Abs(got-expected) / Abs(expected)
		} else {
			// For extremely small values, use absolute error
			relError = Abs(got - expected)
		}

		if relError > maxRelError {
			maxRelError = relError
			maxErrorAt = x
		}
	}

	t.Logf("Maximum relative error: %e at x=%v", maxRelError, maxErrorAt)

	if maxRelError > 2e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxRelError, maxErrorAt)
	}
}

func TestExpMonotonicity(t *testing.T) {
	// Exp should be strictly increasing
	prev := Exp(-100)
	testValues := []float32{
		-100, -50, -20, -10, -5, -2, -1, -0.5, 0, 0.5, 1, 2, 5, 10, 20, 50, 80,
	}

	for _, x := range testValues {
		curr := Exp(x)
		if curr < prev {
			t.Errorf("Exp not monotonic: Exp(%v)=%v < previous=%v", x, curr, prev)
		}
		prev = curr
	}
}

func TestExpZero(t *testing.T) {
	// exp(0) should be exactly 1
	result := Exp(0)
	if result != 1.0 {
		t.Errorf("Exp(0) = %v, want 1.0", result)
	}
}

func TestLdexp32(t *testing.T) {
	// Test the ldexp32 helper function
	tests := []struct {
		name string
		x    float32
		exp  int32
		want float32
	}{
		{"1.0 * 2^0", 1.0, 0, 1.0},
		{"1.0 * 2^1", 1.0, 1, 2.0},
		{"1.0 * 2^2", 1.0, 2, 4.0},
		{"1.0 * 2^10", 1.0, 10, 1024.0},
		{"1.0 * 2^-1", 1.0, -1, 0.5},
		{"1.0 * 2^-2", 1.0, -2, 0.25},
		{"2.5 * 2^3", 2.5, 3, 20.0},
		{"0.5 * 2^-1", 0.5, -1, 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ldexp32(tt.x, tt.exp)
			if Abs(got-tt.want) > 1e-6 {
				t.Errorf("ldexp32(%v, %d) = %v, want %v", tt.x, tt.exp, got, tt.want)
			}
		})
	}
}

// Benchmark Exp
func BenchmarkExpZero(b *testing.B) {
	x := float32(0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Exp(x)
	}
	_ = result
}

func BenchmarkExpSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Exp(x)
	}
	_ = result
}

func BenchmarkExpMedium(b *testing.B) {
	x := float32(5.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Exp(x)
	}
	_ = result
}

func BenchmarkExpLarge(b *testing.B) {
	x := float32(50.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Exp(x)
	}
	_ = result
}

func BenchmarkExpNegative(b *testing.B) {
	x := float32(-5.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Exp(x)
	}
	_ = result
}

func BenchmarkExpFloat64(b *testing.B) {
	x := 5.0
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Exp(x)
	}
	_ = result
}
