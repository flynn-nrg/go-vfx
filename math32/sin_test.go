package math32

import (
	"math"
	"testing"
)

func TestSin(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), 0},

		// Standard angles
		{"π/6 (30°)", float32(Pi / 6), math.Sin(math.Pi / 6)}, // 0.5
		{"π/4 (45°)", float32(Pi / 4), math.Sin(math.Pi / 4)}, // √2/2
		{"π/3 (60°)", float32(Pi / 3), math.Sin(math.Pi / 3)}, // √3/2
		{"π/2 (90°)", float32(Pi / 2), math.Sin(math.Pi / 2)}, // 1
		{"2π/3 (120°)", float32(2 * Pi / 3), math.Sin(2 * math.Pi / 3)},
		{"3π/4 (135°)", float32(3 * Pi / 4), math.Sin(3 * math.Pi / 4)},
		{"5π/6 (150°)", float32(5 * Pi / 6), math.Sin(5 * math.Pi / 6)},
		{"π (180°)", float32(Pi), math.Sin(math.Pi)}, // ~0

		// Negative angles
		{"-π/6", float32(-Pi / 6), math.Sin(-math.Pi / 6)},
		{"-π/4", float32(-Pi / 4), math.Sin(-math.Pi / 4)},
		{"-π/2", float32(-Pi / 2), math.Sin(-math.Pi / 2)},
		{"-π", float32(-Pi), math.Sin(-math.Pi)},

		// Large values (test range reduction)
		{"2π", float32(2 * Pi), math.Sin(2 * math.Pi)},
		{"3π", float32(3 * Pi), math.Sin(3 * math.Pi)},
		{"10π", float32(10 * Pi), math.Sin(10 * math.Pi)},
		{"100", 100.0, math.Sin(100.0)},

		// Small values
		{"0.1", 0.1, math.Sin(0.1)},
		{"0.01", 0.01, math.Sin(0.01)},
		{"-0.1", -0.1, math.Sin(-0.1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sin(tt.input)
			expected := float32(tt.expected)

			tolerance := float32(2e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Sin(%v) = %v, want %v (diff: %e)", tt.input, got, expected, diff)
			}
		})
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64
	}{
		// Special cases
		{"zero", 0, 1},

		// Standard angles
		{"π/6 (30°)", float32(Pi / 6), math.Cos(math.Pi / 6)}, // √3/2
		{"π/4 (45°)", float32(Pi / 4), math.Cos(math.Pi / 4)}, // √2/2
		{"π/3 (60°)", float32(Pi / 3), math.Cos(math.Pi / 3)}, // 0.5
		{"π/2 (90°)", float32(Pi / 2), math.Cos(math.Pi / 2)}, // ~0
		{"2π/3 (120°)", float32(2 * Pi / 3), math.Cos(2 * math.Pi / 3)},
		{"3π/4 (135°)", float32(3 * Pi / 4), math.Cos(3 * math.Pi / 4)},
		{"5π/6 (150°)", float32(5 * Pi / 6), math.Cos(5 * math.Pi / 6)},
		{"π (180°)", float32(Pi), math.Cos(math.Pi)}, // -1

		// Negative angles (cos is even: cos(-x) = cos(x))
		{"-π/6", float32(-Pi / 6), math.Cos(-math.Pi / 6)},
		{"-π/4", float32(-Pi / 4), math.Cos(-math.Pi / 4)},
		{"-π/2", float32(-Pi / 2), math.Cos(-math.Pi / 2)},
		{"-π", float32(-Pi), math.Cos(-math.Pi)},

		// Large values (test range reduction)
		{"2π", float32(2 * Pi), math.Cos(2 * math.Pi)},
		{"3π", float32(3 * Pi), math.Cos(3 * math.Pi)},
		{"10π", float32(10 * Pi), math.Cos(10 * math.Pi)},
		{"100", 100.0, math.Cos(100.0)},

		// Small values
		{"0.1", 0.1, math.Cos(0.1)},
		{"0.01", 0.01, math.Cos(0.01)},
		{"-0.1", -0.1, math.Cos(-0.1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cos(tt.input)
			expected := float32(tt.expected)

			tolerance := float32(2e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Cos(%v) = %v, want %v (diff: %e)", tt.input, got, expected, diff)
			}
		})
	}
}

func TestSinSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Sin(float32(math.NaN())); !IsNaN(nan) {
		t.Errorf("Sin(NaN) = %v, want NaN", nan)
	}

	// Test infinity
	if nan := Sin(float32(math.Inf(1))); !IsNaN(nan) {
		t.Errorf("Sin(+Inf) = %v, want NaN", nan)
	}
	if nan := Sin(float32(math.Inf(-1))); !IsNaN(nan) {
		t.Errorf("Sin(-Inf) = %v, want NaN", nan)
	}
}

func TestCosSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Cos(float32(math.NaN())); !IsNaN(nan) {
		t.Errorf("Cos(NaN) = %v, want NaN", nan)
	}

	// Test infinity
	if nan := Cos(float32(math.Inf(1))); !IsNaN(nan) {
		t.Errorf("Cos(+Inf) = %v, want NaN", nan)
	}
	if nan := Cos(float32(math.Inf(-1))); !IsNaN(nan) {
		t.Errorf("Cos(-Inf) = %v, want NaN", nan)
	}
}

func TestSinCosIdentity(t *testing.T) {
	// Test that sin²(x) + cos²(x) = 1 for various angles
	testValues := []float32{
		0, 0.1, 0.5, 1.0,
		float32(Pi / 6), float32(Pi / 4), float32(Pi / 3), float32(Pi / 2),
		float32(2 * Pi / 3), float32(3 * Pi / 4), float32(Pi),
		float32(2 * Pi), float32(3 * Pi),
		-0.5, float32(-Pi / 4), float32(-Pi / 2), float32(-Pi),
		10.0, 100.0, -100.0,
	}

	for _, x := range testValues {
		s := Sin(x)
		c := Cos(x)
		sum := s*s + c*c

		diff := Abs(sum - 1)
		if diff > 1e-6 {
			t.Errorf("sin²(%v) + cos²(%v) = %v, want 1.0 (diff: %e)", x, x, sum, diff)
		}
	}
}

func TestSinSymmetry(t *testing.T) {
	// Test that sin(-x) = -sin(x) (odd function)
	testValues := []float32{
		0.1, 0.5, 1.0,
		float32(Pi / 6), float32(Pi / 4), float32(Pi / 2),
		float32(2 * Pi), 10.0, 100.0,
	}

	for _, x := range testValues {
		sinPos := Sin(x)
		sinNeg := Sin(-x)

		diff := Abs(sinPos + sinNeg)
		if diff > 1e-6 {
			t.Errorf("sin(%v) = %v, sin(%v) = %v, expected sin(-x) = -sin(x) (diff: %e)",
				x, sinPos, -x, sinNeg, diff)
		}
	}
}

func TestCosSymmetry(t *testing.T) {
	// Test that cos(-x) = cos(x) (even function)
	testValues := []float32{
		0.1, 0.5, 1.0,
		float32(Pi / 6), float32(Pi / 4), float32(Pi / 2),
		float32(2 * Pi), 10.0, 100.0,
	}

	for _, x := range testValues {
		cosPos := Cos(x)
		cosNeg := Cos(-x)

		diff := Abs(cosPos - cosNeg)
		if diff > 1e-6 {
			t.Errorf("cos(%v) = %v, cos(%v) = %v, expected cos(-x) = cos(x) (diff: %e)",
				x, cosPos, -x, cosNeg, diff)
		}
	}
}

func TestSinPeriodicity(t *testing.T) {
	// Test that sin(x + 2π) = sin(x)
	testValues := []float32{0.5, 1.0, float32(Pi / 4), float32(Pi / 2)}

	for _, x := range testValues {
		sin0 := Sin(x)
		sin1 := Sin(x + float32(2*Pi))
		sin2 := Sin(x + float32(4*Pi))

		diff1 := Abs(sin0 - sin1)
		diff2 := Abs(sin0 - sin2)

		if diff1 > 2e-6 {
			t.Errorf("sin(%v) = %v, sin(%v+2π) = %v (diff: %e)", x, sin0, x, sin1, diff1)
		}
		if diff2 > 2e-6 {
			t.Errorf("sin(%v) = %v, sin(%v+4π) = %v (diff: %e)", x, sin0, x, sin2, diff2)
		}
	}
}

func TestCosPeriodicity(t *testing.T) {
	// Test that cos(x + 2π) = cos(x)
	testValues := []float32{0.5, 1.0, float32(Pi / 4), float32(Pi / 2)}

	for _, x := range testValues {
		cos0 := Cos(x)
		cos1 := Cos(x + float32(2*Pi))
		cos2 := Cos(x + float32(4*Pi))

		diff1 := Abs(cos0 - cos1)
		diff2 := Abs(cos0 - cos2)

		if diff1 > 2e-6 {
			t.Errorf("cos(%v) = %v, cos(%v+2π) = %v (diff: %e)", x, cos0, x, cos1, diff1)
		}
		if diff2 > 2e-6 {
			t.Errorf("cos(%v) = %v, cos(%v+4π) = %v (diff: %e)", x, cos0, x, cos2, diff2)
		}
	}
}

func TestSinAccuracy(t *testing.T) {
	var maxError float32
	var maxErrorAt float32

	// Test over several periods
	steps := 1000
	for i := -steps; i <= steps; i++ {
		x := float32(i) * float32(Pi) / 100.0
		got := Sin(x)
		expected := float32(math.Sin(float64(x)))

		err := Abs(got - expected)
		if err > maxError {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Sin: Maximum error: %e at x=%v", maxError, maxErrorAt)

	if maxError > 2e-6 {
		t.Errorf("Sin: Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

func TestCosAccuracy(t *testing.T) {
	var maxError float32
	var maxErrorAt float32

	// Test over several periods
	steps := 1000
	for i := -steps; i <= steps; i++ {
		x := float32(i) * float32(Pi) / 100.0
		got := Cos(x)
		expected := float32(math.Cos(float64(x)))

		err := Abs(got - expected)
		if err > maxError {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Cos: Maximum error: %e at x=%v", maxError, maxErrorAt)

	if maxError > 2e-6 {
		t.Errorf("Cos: Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

// Benchmark Sin
func BenchmarkSinSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sin(x)
	}
	_ = result
}

func BenchmarkSinStandard(b *testing.B) {
	x := float32(Pi / 4)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sin(x)
	}
	_ = result
}

func BenchmarkSinLarge(b *testing.B) {
	x := float32(100.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Sin(x)
	}
	_ = result
}

func BenchmarkSinFloat64(b *testing.B) {
	x := math.Pi / 4
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Sin(x)
	}
	_ = result
}

// Benchmark Cos
func BenchmarkCosSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Cos(x)
	}
	_ = result
}

func BenchmarkCosStandard(b *testing.B) {
	x := float32(Pi / 4)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Cos(x)
	}
	_ = result
}

func BenchmarkCosLarge(b *testing.B) {
	x := float32(100.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Cos(x)
	}
	_ = result
}

func BenchmarkCosFloat64(b *testing.B) {
	x := math.Pi / 4
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Cos(x)
	}
	_ = result
}

// Tan tests

func TestTan(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), 0},

		// Standard angles
		{"π/6 (30°)", float32(Pi / 6), math.Tan(math.Pi / 6)},           // 1/√3
		{"π/4 (45°)", float32(Pi / 4), math.Tan(math.Pi / 4)},           // 1
		{"π/3 (60°)", float32(Pi / 3), math.Tan(math.Pi / 3)},           // √3
		{"3π/4 (135°)", float32(3 * Pi / 4), math.Tan(3 * math.Pi / 4)}, // -1

		// Negative angles
		{"-π/6", float32(-Pi / 6), math.Tan(-math.Pi / 6)},
		{"-π/4", float32(-Pi / 4), math.Tan(-math.Pi / 4)},
		{"-π/3", float32(-Pi / 3), math.Tan(-math.Pi / 3)},

		// Near singularities (π/2 is undefined, test near it)
		{"π/2 - 0.1", float32(Pi/2 - 0.1), math.Tan(math.Pi/2 - 0.1)},
		{"-π/2 + 0.1", float32(-Pi/2 + 0.1), math.Tan(-math.Pi/2 + 0.1)},

		// Periodicity (tan has period π)
		{"π + π/4", float32(Pi + Pi/4), math.Tan(math.Pi + math.Pi/4)},
		{"2π + π/4", float32(2*Pi + Pi/4), math.Tan(2*math.Pi + math.Pi/4)},

		// Small values
		{"0.1", 0.1, math.Tan(0.1)},
		{"0.01", 0.01, math.Tan(0.01)},
		{"-0.1", -0.1, math.Tan(-0.1)},

		// Large values (test range reduction)
		{"10", 10.0, math.Tan(10.0)},
		{"100", 100.0, math.Tan(100.0)},
		{"-10", -10.0, math.Tan(-10.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tan(tt.input)
			expected := float32(tt.expected)

			// Tan can have large values near singularities
			// Use relative error for moderate values
			tolerance := float32(2e-6)
			var err float32
			if Abs(expected) < 100 {
				err = Abs(got - expected)
				if Abs(expected) > 0.001 {
					err = err / Abs(expected)
				}
			} else {
				// For large values, use absolute error
				err = Abs(got-expected) / Abs(expected)
			}

			if err > tolerance {
				t.Errorf("Tan(%v) = %v, want %v (error: %e)", tt.input, got, expected, err)
			}
		})
	}
}

func TestTanSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Tan(NaN()); !IsNaN(nan) {
		t.Errorf("Tan(NaN) = %v, want NaN", nan)
	}

	// Test infinity
	if nan := Tan(float32(math.Inf(1))); !IsNaN(nan) {
		t.Errorf("Tan(+Inf) = %v, want NaN", nan)
	}
	if nan := Tan(float32(math.Inf(-1))); !IsNaN(nan) {
		t.Errorf("Tan(-Inf) = %v, want NaN", nan)
	}
}

func TestTanSymmetry(t *testing.T) {
	// Test that tan(-x) = -tan(x) (odd function)
	testValues := []float32{
		0.1, 0.5, 1.0,
		float32(Pi / 6), float32(Pi / 4), float32(Pi / 3),
		10.0,
	}

	for _, x := range testValues {
		tanPos := Tan(x)
		tanNeg := Tan(-x)

		diff := Abs(tanPos + tanNeg)
		// Use relative error for moderate values
		threshold := float32(1e-6)
		if Abs(tanPos) > 0.001 {
			diff = diff / Abs(tanPos)
		}

		if diff > threshold {
			t.Errorf("tan(%v) = %v, tan(%v) = %v, expected tan(-x) = -tan(x) (diff: %e)",
				x, tanPos, -x, tanNeg, diff)
		}
	}
}

func TestTanPeriodicity(t *testing.T) {
	// Test that tan(x + π) = tan(x) (period is π, not 2π!)
	testValues := []float32{0.5, 1.0, float32(Pi / 4), float32(Pi / 3)}

	for _, x := range testValues {
		tan0 := Tan(x)
		tan1 := Tan(x + float32(Pi))
		tan2 := Tan(x + float32(2*Pi))

		diff1 := Abs(tan0 - tan1)
		diff2 := Abs(tan0 - tan2)

		// Relative error
		threshold := float32(2e-5)
		if Abs(tan0) > 0.01 {
			diff1 = diff1 / Abs(tan0)
			diff2 = diff2 / Abs(tan0)
		}

		if diff1 > threshold {
			t.Errorf("tan(%v) = %v, tan(%v+π) = %v (diff: %e)", x, tan0, x, tan1, diff1)
		}
		if diff2 > threshold {
			t.Errorf("tan(%v) = %v, tan(%v+2π) = %v (diff: %e)", x, tan0, x, tan2, diff2)
		}
	}
}

func TestTanIdentity(t *testing.T) {
	// Test that tan(x) = sin(x) / cos(x)
	testValues := []float32{0.1, 0.5, 1.0, float32(Pi / 6), float32(Pi / 4)}

	for _, x := range testValues {
		tanX := Tan(x)
		sinX := Sin(x)
		cosX := Cos(x)
		ratio := sinX / cosX

		diff := Abs(tanX - ratio)
		threshold := float32(1e-6)
		if Abs(tanX) > 0.01 {
			diff = diff / Abs(tanX)
		}

		if diff > threshold {
			t.Errorf("tan(%v) = %v, sin/cos = %v (diff: %e)", x, tanX, ratio, diff)
		}
	}
}

func TestTanAccuracy(t *testing.T) {
	var maxRelError float32
	var maxErrorAt float32

	// Test over range, avoiding near π/2 singularities
	steps := 100
	for i := -steps; i <= steps; i++ {
		// Skip values very close to π/2 + k·π (singularities)
		x := float32(i) * float32(Pi) / 100.0

		// Skip if close to singularity
		reduced := x
		for Abs(reduced) > Pi/2 {
			if reduced > 0 {
				reduced -= float32(Pi)
			} else {
				reduced += float32(Pi)
			}
		}
		if Abs(Abs(reduced)-float32(Pi/2)) < 0.1 {
			continue // Skip near singularity
		}

		got := Tan(x)
		expected := float32(math.Tan(float64(x)))

		// Skip if expected value is very large (near singularity)
		if Abs(expected) > 100 {
			continue
		}

		var err float32
		if Abs(expected) > 0.01 {
			err = Abs(got-expected) / Abs(expected)
		} else {
			err = Abs(got - expected)
		}

		if err > maxRelError {
			maxRelError = err
			maxErrorAt = x
		}
	}

	t.Logf("Tan: Maximum relative error: %e at x=%v", maxRelError, maxErrorAt)

	if maxRelError > 2e-5 {
		t.Errorf("Tan: Maximum error %e exceeds tolerance at x=%v", maxRelError, maxErrorAt)
	}
}

// Benchmark Tan
func BenchmarkTanSmall(b *testing.B) {
	x := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Tan(x)
	}
	_ = result
}

func BenchmarkTanStandard(b *testing.B) {
	x := float32(Pi / 4)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Tan(x)
	}
	_ = result
}

func BenchmarkTanLarge(b *testing.B) {
	x := float32(10.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Tan(x)
	}
	_ = result
}

func BenchmarkTanFloat64(b *testing.B) {
	x := math.Pi / 4
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Tan(x)
	}
	_ = result
}
