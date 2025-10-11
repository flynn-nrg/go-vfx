package math32

import (
	"math"
	"testing"
)

func TestIsNaN(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected bool
	}{
		// NaN cases
		{"NaN", float32(math.NaN()), true},

		// Non-NaN cases
		{"zero", 0, false},
		{"positive", 1.5, false},
		{"negative", -2.5, false},
		{"positive infinity", float32(math.Inf(1)), false},
		{"negative infinity", float32(math.Inf(-1)), false},
		{"large number", 1e38, false},
		{"small number", 1e-38, false},
		{"negative zero", float32(math.Copysign(0, -1)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNaN(tt.input)
			if got != tt.expected {
				t.Errorf("IsNaN(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsInf(t *testing.T) {
	posInf := float32(math.Inf(1))
	negInf := float32(math.Inf(-1))
	nan := float32(math.NaN())

	tests := []struct {
		name     string
		input    float32
		sign     int
		expected bool
	}{
		// Positive infinity
		{"positive infinity, sign=0", posInf, 0, true},
		{"positive infinity, sign=1", posInf, 1, true},
		{"positive infinity, sign=-1", posInf, -1, false},

		// Negative infinity
		{"negative infinity, sign=0", negInf, 0, true},
		{"negative infinity, sign=1", negInf, 1, false},
		{"negative infinity, sign=-1", negInf, -1, true},

		// Non-infinity values
		{"zero, sign=0", 0, 0, false},
		{"positive, sign=0", 1.5, 0, false},
		{"negative, sign=0", -2.5, 0, false},
		{"NaN, sign=0", nan, 0, false},
		{"large positive, sign=1", 1e38, 1, false},
		{"large negative, sign=-1", -1e38, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInf(tt.input, tt.sign)
			if got != tt.expected {
				t.Errorf("IsInf(%v, %d) = %v, want %v", tt.input, tt.sign, got, tt.expected)
			}
		})
	}
}

func TestSignbit(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected bool
	}{
		// Positive values
		{"positive zero", 0, false},
		{"positive", 1.5, false},
		{"positive small", 1e-38, false},
		{"positive large", 1e38, false},
		{"positive infinity", float32(math.Inf(1)), false},

		// Negative values
		{"negative zero", float32(math.Copysign(0, -1)), true},
		{"negative", -1.5, true},
		{"negative small", -1e-38, true},
		{"negative large", -1e38, true},
		{"negative infinity", float32(math.Inf(-1)), true},

		// NaN (both positive and negative)
		{"positive NaN", float32(math.NaN()), false},
		{"negative NaN", float32(math.Copysign(math.NaN(), -1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Signbit(tt.input)
			if got != tt.expected {
				t.Errorf("Signbit(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsNaNCompatibility(t *testing.T) {
	// Test that our IsNaN matches math.IsNaN behavior
	testValues := []float32{
		0, 1, -1, 1.5, -2.5,
		float32(math.Inf(1)),
		float32(math.Inf(-1)),
		float32(math.NaN()),
		1e38, -1e38, 1e-38, -1e-38,
	}

	for _, x := range testValues {
		our := IsNaN(x)
		std := math.IsNaN(float64(x))
		if our != std {
			t.Errorf("IsNaN(%v): ours=%v, math.IsNaN=%v", x, our, std)
		}
	}
}

func TestIsInfCompatibility(t *testing.T) {
	// Test that our IsInf matches math.IsInf behavior
	testValues := []float32{
		0, 1, -1, 1.5, -2.5,
		float32(math.Inf(1)),
		float32(math.Inf(-1)),
		float32(math.NaN()),
		1e38, -1e38,
	}

	signs := []int{-1, 0, 1}

	for _, x := range testValues {
		for _, sign := range signs {
			our := IsInf(x, sign)
			std := math.IsInf(float64(x), sign)
			if our != std {
				t.Errorf("IsInf(%v, %d): ours=%v, math.IsInf=%v", x, sign, our, std)
			}
		}
	}
}

func TestSignbitCompatibility(t *testing.T) {
	// Test that our Signbit matches math.Signbit behavior
	testValues := []float32{
		0,
		float32(math.Copysign(0, -1)),
		1, -1, 1.5, -2.5,
		float32(math.Inf(1)),
		float32(math.Inf(-1)),
		float32(math.NaN()),
		float32(math.Copysign(math.NaN(), -1)),
		1e38, -1e38, 1e-38, -1e-38,
	}

	for _, x := range testValues {
		our := Signbit(x)
		std := math.Signbit(float64(x))
		if our != std {
			t.Errorf("Signbit(%v): ours=%v, math.Signbit=%v", x, our, std)
		}
	}
}

func TestBitPatterns(t *testing.T) {
	// Test specific bit patterns to ensure correctness

	// Positive infinity: 0x7F800000
	posInf := math.Float32frombits(0x7F800000)
	if !IsInf(posInf, 1) {
		t.Error("Failed to recognize positive infinity bit pattern")
	}
	if Signbit(posInf) {
		t.Error("Positive infinity should not have sign bit set")
	}

	// Negative infinity: 0xFF800000
	negInf := math.Float32frombits(0xFF800000)
	if !IsInf(negInf, -1) {
		t.Error("Failed to recognize negative infinity bit pattern")
	}
	if !Signbit(negInf) {
		t.Error("Negative infinity should have sign bit set")
	}

	// Positive NaN: 0x7FC00000 (quiet NaN)
	posNaN := math.Float32frombits(0x7FC00000)
	if !IsNaN(posNaN) {
		t.Error("Failed to recognize positive NaN bit pattern")
	}
	if Signbit(posNaN) {
		t.Error("Positive NaN should not have sign bit set")
	}

	// Negative NaN: 0xFFC00000
	negNaN := math.Float32frombits(0xFFC00000)
	if !IsNaN(negNaN) {
		t.Error("Failed to recognize negative NaN bit pattern")
	}
	if !Signbit(negNaN) {
		t.Error("Negative NaN should have sign bit set")
	}

	// Negative zero: 0x80000000
	negZero := math.Float32frombits(0x80000000)
	if negZero != 0 {
		t.Error("Negative zero should compare equal to zero")
	}
	if !Signbit(negZero) {
		t.Error("Negative zero should have sign bit set")
	}
}

// Benchmark IsNaN
func BenchmarkIsNaN(b *testing.B) {
	x := float32(1.5)
	var result bool
	for i := 0; i < b.N; i++ {
		result = IsNaN(x)
	}
	_ = result
}

func BenchmarkIsNaNTrue(b *testing.B) {
	x := float32(math.NaN())
	var result bool
	for i := 0; i < b.N; i++ {
		result = IsNaN(x)
	}
	_ = result
}

// Benchmark IsInf
func BenchmarkIsInf(b *testing.B) {
	x := float32(1.5)
	var result bool
	for i := 0; i < b.N; i++ {
		result = IsInf(x, 0)
	}
	_ = result
}

func BenchmarkIsInfTrue(b *testing.B) {
	x := float32(math.Inf(1))
	var result bool
	for i := 0; i < b.N; i++ {
		result = IsInf(x, 1)
	}
	_ = result
}

// Benchmark Signbit
func BenchmarkSignbit(b *testing.B) {
	x := float32(1.5)
	var result bool
	for i := 0; i < b.N; i++ {
		result = Signbit(x)
	}
	_ = result
}

func BenchmarkSignbitNegative(b *testing.B) {
	x := float32(-1.5)
	var result bool
	for i := 0; i < b.N; i++ {
		result = Signbit(x)
	}
	_ = result
}

// Compare with standard library
func BenchmarkMathIsNaN(b *testing.B) {
	x := 1.5
	var result bool
	for i := 0; i < b.N; i++ {
		result = math.IsNaN(x)
	}
	_ = result
}

func BenchmarkMathIsInf(b *testing.B) {
	x := 1.5
	var result bool
	for i := 0; i < b.N; i++ {
		result = math.IsInf(x, 0)
	}
	_ = result
}

func BenchmarkMathSignbit(b *testing.B) {
	x := 1.5
	var result bool
	for i := 0; i < b.N; i++ {
		result = math.Signbit(x)
	}
	_ = result
}
