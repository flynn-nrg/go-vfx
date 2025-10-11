//go:build arm64

package math32

// sqrt is implemented in sqrt_arm64.s using the FSQRTS instruction
func sqrt(x float32) float32
