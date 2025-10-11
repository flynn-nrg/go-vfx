package math32

import (
	"math"
	"testing"
)

func TestAtan(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float64 // use float64 for reference
	}{
		// Special cases
		{"zero", 0, 0},
		{"negative zero", float32(math.Copysign(0, -1)), 0},

		// Small values (direct polynomial path)
		{"small positive", 0.1, math.Atan(0.1)},
		{"small negative", -0.1, math.Atan(-0.1)},
		{"0.2", 0.2, math.Atan(0.2)},
		{"0.4", 0.4, math.Atan(0.4)},

		// Medium values (π/4 + atan((x-1)/(x+1)) path)
		{"0.5", 0.5, math.Atan(0.5)},
		{"1.0", 1.0, math.Atan(1.0)},
		{"1.5", 1.5, math.Atan(1.5)},
		{"2.0", 2.0, math.Atan(2.0)},

		// Large values (π/2 - atan(1/x) path)
		{"3.0", 3.0, math.Atan(3.0)},
		{"5.0", 5.0, math.Atan(5.0)},
		{"10.0", 10.0, math.Atan(10.0)},
		{"100.0", 100.0, math.Atan(100.0)},

		// Negative values
		{"-0.5", -0.5, math.Atan(-0.5)},
		{"-1.0", -1.0, math.Atan(-1.0)},
		{"-2.0", -2.0, math.Atan(-2.0)},
		{"-10.0", -10.0, math.Atan(-10.0)},

		// Mathematical constants
		{"sqrt(3)", float32(math.Sqrt(3)), math.Atan(math.Sqrt(3))},
		{"1/sqrt(3)", float32(1 / math.Sqrt(3)), math.Atan(1 / math.Sqrt(3))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan(tt.input)
			expected := float32(tt.expected)

			// Allow for small floating point errors
			tolerance := float32(2e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Atan(%v) = %v, want %v (diff: %v)", tt.input, got, expected, diff)
			}
		})
	}
}

func TestAtanSpecialCases(t *testing.T) {
	// Test NaN
	if nan := Atan(float32(math.NaN())); !math.IsNaN(float64(nan)) {
		t.Errorf("Atan(NaN) = %v, want NaN", nan)
	}

	// Test +Inf
	result := Atan(float32(math.Inf(1)))
	expected := float32(Pi / 2)
	if Abs(result-expected) > 1e-6 {
		t.Errorf("Atan(+Inf) = %v, want π/2 = %v", result, expected)
	}

	// Test -Inf
	result = Atan(float32(math.Inf(-1)))
	expected = float32(-Pi / 2)
	if Abs(result-expected) > 1e-6 {
		t.Errorf("Atan(-Inf) = %v, want -π/2 = %v", result, expected)
	}

	// Test atan(1) = π/4
	result = Atan(1)
	expected = float32(Pi / 4)
	if Abs(result-expected) > 1e-6 {
		t.Errorf("Atan(1) = %v, want π/4 = %v", result, expected)
	}
}

func TestAtanAccuracy(t *testing.T) {
	// Test a range of values and measure maximum error
	var maxError float32
	var maxErrorAt float32

	testValues := []float32{
		-100, -50, -20, -10, -5, -3, -2, -1.5, -1, -0.5, -0.2, -0.1,
		0, 0.1, 0.2, 0.5, 1, 1.5, 2, 3, 5, 10, 20, 50, 100,
	}

	for _, x := range testValues {
		got := Atan(x)
		expected := float32(math.Atan(float64(x)))

		err := Abs(got - expected)
		if err > maxError {
			maxError = err
			maxErrorAt = x
		}
	}

	t.Logf("Maximum error: %e at x=%v", maxError, maxErrorAt)

	// For float32, we expect errors to be very small
	if maxError > 2e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at x=%v", maxError, maxErrorAt)
	}
}

func TestAtanMonotonicity(t *testing.T) {
	// Atan should be monotonically increasing
	prev := Atan(-100)
	testValues := []float32{
		-100, -10, -5, -2, -1, -0.5, -0.1, 0, 0.1, 0.5, 1, 2, 5, 10, 100,
	}

	for _, x := range testValues {
		curr := Atan(x)
		if curr < prev {
			t.Errorf("Atan not monotonic: Atan(%v)=%v < previous=%v", x, curr, prev)
		}
		prev = curr
	}
}

func TestAtan2(t *testing.T) {
	tests := []struct {
		name     string
		y, x     float32
		expected float64
	}{
		// Basic quadrants
		{"Q1: (1,1)", 1, 1, math.Pi / 4},
		{"Q2: (1,-1)", 1, -1, 3 * math.Pi / 4},
		{"Q3: (-1,-1)", -1, -1, -3 * math.Pi / 4},
		{"Q4: (-1,1)", -1, 1, -math.Pi / 4},

		// Axes
		{"positive x-axis", 0, 1, 0},
		{"negative x-axis", 0, -1, math.Pi},
		{"positive y-axis", 1, 0, math.Pi / 2},
		{"negative y-axis", -1, 0, -math.Pi / 2},

		// Origin
		{"origin", 0, 0, 0},

		// Various angles
		{"30 degrees", 0.5, float32(math.Sqrt(3) / 2), math.Pi / 6},
		{"60 degrees", float32(math.Sqrt(3) / 2), 0.5, math.Pi / 3},
		{"120 degrees", float32(math.Sqrt(3) / 2), -0.5, 2 * math.Pi / 3},
		{"150 degrees", 0.5, float32(-math.Sqrt(3) / 2), 5 * math.Pi / 6},

		// Large values
		{"large y", 1000, 1, math.Atan2(1000, 1)},
		{"large x", 1, 1000, math.Atan2(1, 1000)},
		{"both large", 1000, 1000, math.Pi / 4},

		// Small values
		{"small y", 0.001, 1, math.Atan2(0.001, 1)},
		{"small x", 1, 0.001, math.Atan2(1, 0.001)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan2(tt.y, tt.x)
			expected := float32(tt.expected)

			tolerance := float32(2e-6)
			diff := Abs(got - expected)

			if diff > tolerance {
				t.Errorf("Atan2(%v, %v) = %v, want %v (diff: %v)",
					tt.y, tt.x, got, expected, diff)
			}
		})
	}
}

func TestAtan2SpecialCases(t *testing.T) {
	// Test NaN propagation
	if nan := Atan2(float32(math.NaN()), 1); !math.IsNaN(float64(nan)) {
		t.Errorf("Atan2(NaN, 1) = %v, want NaN", nan)
	}
	if nan := Atan2(1, float32(math.NaN())); !math.IsNaN(float64(nan)) {
		t.Errorf("Atan2(1, NaN) = %v, want NaN", nan)
	}

	// Test infinities
	testCases := []struct {
		name     string
		y, x     float32
		expected float32
	}{
		{"+Inf, +Inf", float32(math.Inf(1)), float32(math.Inf(1)), float32(Pi / 4)},
		{"+Inf, -Inf", float32(math.Inf(1)), float32(math.Inf(-1)), float32(3 * Pi / 4)},
		{"-Inf, +Inf", float32(math.Inf(-1)), float32(math.Inf(1)), float32(-Pi / 4)},
		{"-Inf, -Inf", float32(math.Inf(-1)), float32(math.Inf(-1)), float32(-3 * Pi / 4)},
		{"+Inf, 1", float32(math.Inf(1)), 1, float32(Pi / 2)},
		{"-Inf, 1", float32(math.Inf(-1)), 1, float32(-Pi / 2)},
		{"1, +Inf", 1, float32(math.Inf(1)), 0},
		{"1, -Inf", 1, float32(math.Inf(-1)), float32(Pi)},
		{"-1, -Inf", -1, float32(math.Inf(-1)), float32(-Pi)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Atan2(tc.y, tc.x)
			if Abs(got-tc.expected) > 1e-6 {
				t.Errorf("Atan2(%v, %v) = %v, want %v", tc.y, tc.x, got, tc.expected)
			}
		})
	}

	// Test zero cases
	if result := Atan2(0, 1); result != 0 {
		t.Errorf("Atan2(0, 1) = %v, want 0", result)
	}

	if result := Atan2(0, -1); Abs(result-float32(Pi)) > 1e-6 {
		t.Errorf("Atan2(0, -1) = %v, want π", result)
	}
}

func TestAtan2Quadrants(t *testing.T) {
	// Verify all four quadrants return correct signs
	eps := float32(0.1)

	// Quadrant I (x>0, y>0): 0 to π/2
	q1 := Atan2(eps, eps)
	if q1 < 0 || q1 > float32(Pi/2) {
		t.Errorf("Q1: Atan2(%v, %v) = %v, expected in (0, π/2)", eps, eps, q1)
	}

	// Quadrant II (x<0, y>0): π/2 to π
	q2 := Atan2(eps, -eps)
	if q2 < float32(Pi/2) || q2 > float32(Pi) {
		t.Errorf("Q2: Atan2(%v, %v) = %v, expected in (π/2, π)", eps, -eps, q2)
	}

	// Quadrant III (x<0, y<0): -π to -π/2
	q3 := Atan2(-eps, -eps)
	if q3 > float32(-Pi/2) || q3 < float32(-Pi) {
		t.Errorf("Q3: Atan2(%v, %v) = %v, expected in (-π, -π/2)", -eps, -eps, q3)
	}

	// Quadrant IV (x>0, y<0): -π/2 to 0
	q4 := Atan2(-eps, eps)
	if q4 > 0 || q4 < float32(-Pi/2) {
		t.Errorf("Q4: Atan2(%v, %v) = %v, expected in (-π/2, 0)", -eps, eps, q4)
	}
}

func TestAtan2Accuracy(t *testing.T) {
	// Test accuracy over a grid of values
	var maxError float32
	var maxErrorY, maxErrorX float32

	testValues := []float32{-10, -5, -2, -1, -0.5, -0.1, 0, 0.1, 0.5, 1, 2, 5, 10}

	for _, y := range testValues {
		for _, x := range testValues {
			if x == 0 && y == 0 {
				continue
			}

			got := Atan2(y, x)
			expected := float32(math.Atan2(float64(y), float64(x)))

			err := Abs(got - expected)
			if err > maxError {
				maxError = err
				maxErrorY = y
				maxErrorX = x
			}
		}
	}

	t.Logf("Maximum error: %e at y=%v, x=%v", maxError, maxErrorY, maxErrorX)

	if maxError > 2e-6 {
		t.Errorf("Maximum error %e exceeds tolerance at y=%v, x=%v",
			maxError, maxErrorY, maxErrorX)
	}
}

// Benchmark Atan with small values
func BenchmarkAtanSmall(b *testing.B) {
	x := float32(0.3)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Atan(x)
	}
	_ = result
}

// Benchmark Atan with medium values
func BenchmarkAtanMedium(b *testing.B) {
	x := float32(1.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Atan(x)
	}
	_ = result
}

// Benchmark Atan with large values
func BenchmarkAtanLarge(b *testing.B) {
	x := float32(10.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Atan(x)
	}
	_ = result
}

// Benchmark float64 Atan for comparison
func BenchmarkAtanFloat64(b *testing.B) {
	x := 1.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Atan(x)
	}
	_ = result
}

// Benchmark Atan2
func BenchmarkAtan2(b *testing.B) {
	y := float32(1.0)
	x := float32(1.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Atan2(y, x)
	}
	_ = result
}

// Benchmark Atan2 with special quadrant
func BenchmarkAtan2Quadrant2(b *testing.B) {
	y := float32(1.0)
	x := float32(-1.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Atan2(y, x)
	}
	_ = result
}

// Benchmark float64 Atan2 for comparison
func BenchmarkAtan2Float64(b *testing.B) {
	y := 1.0
	x := 1.5
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Atan2(y, x)
	}
	_ = result
}
