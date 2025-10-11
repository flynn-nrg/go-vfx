//go:build amd64

package math32

// sqrt is implemented in sqrt_amd64.s using the SQRTSS instruction
func sqrt(x float32) float32
