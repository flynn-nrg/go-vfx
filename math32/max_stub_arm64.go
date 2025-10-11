//go:build arm64

package math32

// max is implemented in max_arm64.s using the FMAXS instruction
func max(x, y float32) float32

// min is implemented in max_arm64.s using the FMINS instruction
func min(x, y float32) float32
