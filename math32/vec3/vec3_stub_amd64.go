//go:build amd64

package vec3

// dot is implemented in vec3_amd64.s using SSE instructions
func dot(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) float32

// add3 is implemented in vec3_amd64.s using SSE instructions
func add3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)

// sub3 is implemented in vec3_amd64.s using SSE instructions
func sub3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)

// mul3 is implemented in vec3_amd64.s using SSE instructions
func mul3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)

// scalarMul3 is implemented in vec3_amd64.s using SSE instructions
func scalarMul3(vX, vY, vZ, scalar float32) (float32, float32, float32)

// cross is implemented in vec3_amd64.s using SSE instructions
func cross(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
