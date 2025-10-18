package math32

import (
	"math"
	"testing"
)

func TestPow(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float32
		expected float64 // use float64 for reference
	}{
		// Simple integer powers
		{"2^2", 2, 2, 4},
		{"2^3", 2, 3, 8},
		{"2^4", 2, 4, 16},
		{"2^8", 2, 8, 256},
		{"2^10", 2, 10, 1024},
		{"3^2", 3, 2, 9},
		{"3^3", 3, 3, 27},
		{"5^2", 5, 2, 25},
		{"10^3", 10, 3, 1000},

		// Negative powers
		{"2^-1", 2, -1, 0.5},
		{"2^-2", 2, -2, 0.25},
		{"2^-3", 2, -3, 0.125},
		{"10^-2", 10, -2, 0.01},

		// Fractional powers (roots)
		{"4^0.5", 4, 0.5, 2},
		{"9^0.5", 9, 0.5, 3},
		{"8^(1/3)", 8, 1.0 / 3, 2},
		{"27^(1/3)", 27, 1.0 / 3, 3},

		// Powers of 1
		{"1^2", 1, 2, 1},
		{"1^-5", 1, -5, 1},
		{"1^0.5", 1, 0.5, 1},
		{"1^100", 1, 100, 1},

		// Powers of 0
		{"0^1", 0, 1, 0},
		{"0^2", 0, 2, 0},
		{"0^10", 0, 10, 0},

		// Base e and natural powers
		{"e^1", float32(E), 1, E},
		{"e^2", float32(E), 2, math.Pow(E, 2)},
		{"e^0.5", float32(E), 0.5, math.Sqrt(E)},

		// Fractional bases
		{"0.5^2", 0.5, 2, 0.25},
		{"0.5^3", 0.5, 3, 0.125},
		{"0.5^-1", 0.5, -1, 2},
		{"0.25^0.5", 0.25, 0.5, 0.5},

		// Mixed fractional
		{"2.5^2", 2.5, 2, 6.25},
		{"1.5^3", 1.5, 3, math.Pow(1.5, 3)},
		{"3.7^2.3", 3.7, 2.3, math.Pow(3.7, 2.3)},

		// Very small bases
		{"0.1^2", 0.1, 2, 0.01},
		{"0.01^2", 0.01, 2, 0.0001},

		// Larger values
		{"100^2", 100, 2, 10000},
		{"1000^0.5", 1000, 0.5, math.Sqrt(1000)},

		// Mathematical identities
		{"sqrt(2)^2", float32(Sqrt2), 2, 2},
		{"pi^1", float32(Pi), 1, Pi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pow(tt.x, tt.y)
			expected := float32(tt.expected)

			// For float32 pow via exp/log chain, we expect modest accuracy
			// The error compounds through multiple transcendental function calls
			tolerance := float32(1e-4)
			diff := Abs(got - expected)
			relError := float32(0)
			if expected != 0 {
				relError = diff / Abs(expected)
			}

			// Use relative error for non-zero values, absolute error for near-zero
			if Abs(expected) > 1e-4 {
				if relError > 0.01 { // 1% relative error tolerance
					t.Errorf("Pow(%v, %v) = %v, want %v (diff: %v, rel error: %e)",
						tt.x, tt.y, got, expected, diff, relError)
				}
			} else {
				if diff > tolerance {
					t.Errorf("Pow(%v, %v) = %v, want %v (diff: %v)",
						tt.x, tt.y, got, expected, diff)
				}
			}
		})
	}
}

func TestPowSpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float32
		expected float32
		checkNaN bool
		checkInf bool
		infSign  int
	}{
		// x^0 = 1 for any x
		{"5^0", 5, 0, 1, false, false, 0},
		{"(-5)^0", -5, 0, 1, false, false, 0},
		{"0^0", 0, 0, 1, false, false, 0},
		{"Inf^0", float32(math.Inf(1)), 0, 1, false, false, 0},
		{"NaN^0", float32(math.NaN()), 0, 1, false, false, 0},

		// 1^y = 1 for any y
		{"1^5", 1, 5, 1, false, false, 0},
		{"1^(-5)", 1, -5, 1, false, false, 0},
		{"1^NaN", 1, float32(math.NaN()), 1, false, false, 0},
		{"1^Inf", 1, float32(math.Inf(1)), 1, false, false, 0},

		// x^1 = x
		{"5^1", 5, 1, 5, false, false, 0},
		{"(-5)^1", -5, 1, -5, false, false, 0},

		// NaN cases
		{"NaN^2", float32(math.NaN()), 2, 0, true, false, 0},
		{"2^NaN", 2, float32(math.NaN()), 0, true, false, 0},
		{"(-2)^0.5", -2, 0.5, 0, true, false, 0}, // negative base, non-integer exponent

		// Zero base cases
		{"0^(-1)", 0, -1, 0, false, true, 1},                   // +Inf
		{"0^(-2)", 0, -2, 0, false, true, 1},                   // +Inf (even)
		{"(-0)^(-1)", Copysign(0, -1), -1, 0, false, true, -1}, // -Inf (odd)
		{"(-0)^(-2)", Copysign(0, -1), -2, 0, false, true, 1},  // +Inf (even)

		// Infinity cases
		{"Inf^2", float32(math.Inf(1)), 2, 0, false, true, 1},
		{"Inf^(-2)", float32(math.Inf(1)), -2, 0, false, false, 0},
		{"(-Inf)^2", float32(math.Inf(-1)), 2, 0, false, true, 1},
		{"(-Inf)^3", float32(math.Inf(-1)), 3, 0, false, true, -1},

		// 2^Inf
		{"2^Inf", 2, float32(math.Inf(1)), 0, false, true, 1},
		{"2^(-Inf)", 2, float32(math.Inf(-1)), 0, false, false, 0},
		{"0.5^Inf", 0.5, float32(math.Inf(1)), 0, false, false, 0},
		{"0.5^(-Inf)", 0.5, float32(math.Inf(-1)), 0, false, true, 1},
		{"(-1)^Inf", -1, float32(math.Inf(1)), 1, false, false, 0},
		{"(-1)^(-Inf)", -1, float32(math.Inf(-1)), 1, false, false, 0},

		// Negative base with integer exponent
		{"(-2)^2", -2, 2, 4, false, false, 0},
		{"(-2)^3", -2, 3, -8, false, false, 0},
		{"(-3)^2", -3, 2, 9, false, false, 0},
		{"(-3)^3", -3, 3, -27, false, false, 0},
		{"(-0.5)^2", -0.5, 2, 0.25, false, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pow(tt.x, tt.y)

			if tt.checkNaN {
				if !IsNaN(got) {
					t.Errorf("Pow(%v, %v) = %v, want NaN", tt.x, tt.y, got)
				}
				return
			}

			if tt.checkInf {
				if !IsInf(got, tt.infSign) {
					infStr := "+Inf"
					if tt.infSign < 0 {
						infStr = "-Inf"
					}
					t.Errorf("Pow(%v, %v) = %v, want %s", tt.x, tt.y, got, infStr)
				}
				return
			}

			if got != tt.expected {
				// Check for sign-preserving zero
				if got == 0 && tt.expected == 0 {
					gotBits := math.Float32bits(got)
					expBits := math.Float32bits(tt.expected)
					if gotBits != expBits {
						t.Errorf("Pow(%v, %v) = %v, want %v (sign mismatch)",
							tt.x, tt.y, got, tt.expected)
					}
					return
				}

				tolerance := float32(1e-4)
				diff := Abs(got - tt.expected)
				relError := float32(0)
				if tt.expected != 0 {
					relError = diff / Abs(tt.expected)
				}

				// Use relative error for non-zero values
				if Abs(tt.expected) > 1e-4 {
					if relError > 0.01 { // 1% tolerance
						t.Errorf("Pow(%v, %v) = %v, want %v (diff: %v, rel error: %e)",
							tt.x, tt.y, got, tt.expected, diff, relError)
					}
				} else {
					if diff > tolerance {
						t.Errorf("Pow(%v, %v) = %v, want %v (diff: %v)",
							tt.x, tt.y, got, tt.expected, diff)
					}
				}
			}
		})
	}
}

func TestPowAccuracy(t *testing.T) {
	// Test accuracy across various ranges
	var maxError float32
	var maxErrorX, maxErrorY float32

	testCases := []struct {
		x, y float32
	}{
		// Small bases
		{0.1, 2}, {0.5, 3}, {0.9, 0.5},
		// Medium bases
		{2, 0.5}, {2, 2}, {2, 3}, {2, -1}, {2, -2},
		{3, 2}, {5, 2}, {7, 0.5},
		// Larger bases
		{10, 2}, {10, 3}, {100, 0.5},
		// Fractional exponents
		{4, 0.333333}, {8, 0.333333}, {16, 0.25},
		{100, 0.5}, {1000, 0.333333},
		// Mixed
		{1.5, 2.5}, {2.7, 1.8}, {3.14, 2.71},
	}

	for _, tc := range testCases {
		got := Pow(tc.x, tc.y)
		expected := float32(math.Pow(float64(tc.x), float64(tc.y)))

		if IsNaN(got) || IsInf(got, 0) || IsNaN(expected) || math.IsInf(float64(expected), 0) {
			continue
		}

		relError := Abs(got-expected) / Abs(expected)
		if relError > maxError {
			maxError = relError
			maxErrorX = tc.x
			maxErrorY = tc.y
		}
	}

	t.Logf("Maximum relative error: %e at Pow(%v, %v)", maxError, maxErrorX, maxErrorY)

	// For float32, we expect reasonable accuracy through exp/log chain
	// Allow more error than direct hardware instructions due to error compounding
	if maxError > 0.01 { // 1% relative error tolerance
		t.Errorf("Maximum relative error %e exceeds tolerance at Pow(%v, %v)",
			maxError, maxErrorX, maxErrorY)
	}
}

func TestPowIdentities(t *testing.T) {
	testValues := []float32{2, 3, 5, 10, 0.5, 1.5, 7}

	for _, x := range testValues {
		// x^2 = x * x
		got := Pow(x, 2)
		expected := x * x
		relError := Abs(got-expected) / expected
		if relError > 0.01 { // 1% tolerance
			t.Errorf("Pow(%v, 2) = %v, want %v", x, got, expected)
		}

		// x^0.5 = sqrt(x)
		got = Pow(x, 0.5)
		expected = Sqrt(x)
		relError = Abs(got-expected) / expected
		if relError > 0.01 { // 1% tolerance
			t.Errorf("Pow(%v, 0.5) = %v, want %v (sqrt)", x, got, expected)
		}

		// x^(-1) = 1/x
		got = Pow(x, -1)
		expected = 1 / x
		relError = Abs(got-expected) / expected
		if relError > 0.01 { // 1% tolerance
			t.Errorf("Pow(%v, -1) = %v, want %v", x, got, expected)
		}

		// (x^a)^b ≈ x^(a*b)
		// This test compounds errors significantly, so we need looser tolerance
		a := float32(2)
		b := float32(3)
		left := Pow(Pow(x, a), b)
		right := Pow(x, a*b)
		relError = Abs(left-right) / right
		if relError > 0.02 { // 2% tolerance due to compounding
			t.Errorf("(Pow(%v, %v))^%v = %v, want Pow(%v, %v) = %v (rel error: %e)",
				x, a, b, left, x, a*b, right, relError)
		}
	}
}

// Benchmark Pow with various inputs
func BenchmarkPow(b *testing.B) {
	x := float32(2.0)
	y := float32(3.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Pow(x, y)
	}
	_ = result
}

func BenchmarkPowSquare(b *testing.B) {
	x := float32(2.0)
	y := float32(2.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Pow(x, y)
	}
	_ = result
}

func BenchmarkPowSqrt(b *testing.B) {
	x := float32(4.0)
	y := float32(0.5)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Pow(x, y)
	}
	_ = result
}

func BenchmarkPowFractional(b *testing.B) {
	x := float32(2.5)
	y := float32(3.7)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Pow(x, y)
	}
	_ = result
}

func BenchmarkPowNegativeExp(b *testing.B) {
	x := float32(2.0)
	y := float32(-3.0)
	var result float32
	for i := 0; i < b.N; i++ {
		result = Pow(x, y)
	}
	_ = result
}

// Benchmark float64 pow for comparison
func BenchmarkPowFloat64(b *testing.B) {
	x := 2.0
	y := 3.0
	var result float64
	for i := 0; i < b.N; i++ {
		result = math.Pow(x, y)
	}
	_ = result
}
