//go:build !amd64 && !arm64

package math32

// max provides a software fallback for architectures without assembly implementation.
func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

// min provides a software fallback for architectures without assembly implementation.
func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}
