//go:build amd64

package math32

// floor is implemented in floor_amd64.s using the ROUNDSS instruction
func floor(x float32) float32

// ceil is implemented in floor_amd64.s using the ROUNDSS instruction
func ceil(x float32) float32

// round is implemented in floor_amd64.s using the ROUNDSS instruction
func round(x float32) float32
