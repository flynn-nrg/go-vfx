//go:build !amd64 && !arm64

package math32

import "math"

// sqrt provides a software fallback for architectures without assembly implementation.
// On ARM64 and AMD64, this is replaced by optimized assembly versions.
func sqrt(x float32) float32 {
	// Use float64 sqrt and convert back
	// This is slower but ensures correctness on all platforms
	return float32(math.Sqrt(float64(x)))
}
