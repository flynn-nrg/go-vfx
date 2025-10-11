//go:build arm64

package math32

// floor is implemented in floor_arm64.s using the FRINTM instruction
func floor(x float32) float32

// ceil is implemented in floor_arm64.s using the FRINTP instruction
func ceil(x float32) float32

// round is implemented in floor_arm64.s using the FRINTN instruction
func round(x float32) float32
