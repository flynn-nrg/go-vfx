//go:build !amd64 && !arm64

package math32

import "math"

// floor provides a software fallback for architectures without assembly implementation.
func floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

// ceil provides a software fallback for architectures without assembly implementation.
func ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

// round provides a software fallback for architectures without assembly implementation.
func round(x float32) float32 {
	return float32(math.Round(float64(x)))
}
