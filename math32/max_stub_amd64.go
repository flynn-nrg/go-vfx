//go:build amd64

package math32

// max is implemented in max_amd64.s using the MAXSS instruction
func max(x, y float32) float32

// min is implemented in max_amd64.s using the MINSS instruction
func min(x, y float32) float32
