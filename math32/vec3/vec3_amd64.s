//go:build amd64.v3

#include "textflag.h"

// func dot(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) float32
TEXT ·dot(SB),NOSPLIT,$0-28
	// Load v1 components
	MOVSS	v1X+0(FP), X0
	MOVSS	v1Y+4(FP), X1
	MOVSS	v1Z+8(FP), X2
	
	// Load v2 components
	MOVSS	v2X+12(FP), X3
	MOVSS	v2Y+16(FP), X4
	MOVSS	v2Z+20(FP), X5
	
	// Use FMA (requires GOAMD64=v3): VFMADD231SS multiplier, multiplicand, accumulator
	// VFMADD231SS: dest = dest + (src1 * src2)
	MULSS	X3, X0            // X0 = v1.X * v2.X
	VFMADD231SS	X4, X1, X0    // X0 = X0 + (X1 * X4) = X0 + (v1.Y * v2.Y)
	VFMADD231SS	X5, X2, X0    // X0 = X0 + (X2 * X5) = X0 + (v1.Z * v2.Z)
	
	MOVSS	X0, ret+24(FP)
	RET

// func add3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·add3(SB),NOSPLIT,$0-36
	// Load v1 into X0, X1, X2
	MOVSS	v1X+0(FP), X0
	MOVSS	v1Y+4(FP), X1
	MOVSS	v1Z+8(FP), X2
	
	// Add v2 components
	ADDSS	v2X+12(FP), X0
	ADDSS	v2Y+16(FP), X1
	ADDSS	v2Z+20(FP), X2
	
	// Store results
	MOVSS	X0, ret0+24(FP)
	MOVSS	X1, ret1+28(FP)
	MOVSS	X2, ret2+32(FP)
	RET

// func sub3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·sub3(SB),NOSPLIT,$0-36
	// Load v1 into X0, X1, X2
	MOVSS	v1X+0(FP), X0
	MOVSS	v1Y+4(FP), X1
	MOVSS	v1Z+8(FP), X2
	
	// Subtract v2 components
	SUBSS	v2X+12(FP), X0
	SUBSS	v2Y+16(FP), X1
	SUBSS	v2Z+20(FP), X2
	
	// Store results
	MOVSS	X0, ret0+24(FP)
	MOVSS	X1, ret1+28(FP)
	MOVSS	X2, ret2+32(FP)
	RET

// func mul3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·mul3(SB),NOSPLIT,$0-36
	// Load v1 into X0, X1, X2
	MOVSS	v1X+0(FP), X0
	MOVSS	v1Y+4(FP), X1
	MOVSS	v1Z+8(FP), X2
	
	// Multiply by v2 components
	MULSS	v2X+12(FP), X0
	MULSS	v2Y+16(FP), X1
	MULSS	v2Z+20(FP), X2
	
	// Store results
	MOVSS	X0, ret0+24(FP)
	MOVSS	X1, ret1+28(FP)
	MOVSS	X2, ret2+32(FP)
	RET

// func scalarMul3(vX, vY, vZ, scalar float32) (float32, float32, float32)
TEXT ·scalarMul3(SB),NOSPLIT,$0-28
	// Load scalar into X3
	MOVSS	scalar+12(FP), X3
	
	// Load vector components and multiply by scalar
	MOVSS	vX+0(FP), X0
	MOVSS	vY+4(FP), X1
	MOVSS	vZ+8(FP), X2
	
	MULSS	X3, X0
	MULSS	X3, X1
	MULSS	X3, X2
	
	// Store results
	MOVSS	X0, ret0+16(FP)
	MOVSS	X1, ret1+20(FP)
	MOVSS	X2, ret2+24(FP)
	RET

// func cross(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·cross(SB),NOSPLIT,$0-36
	// Load v1 and v2
	MOVSS	v1X+0(FP), X0     // v1.X
	MOVSS	v1Y+4(FP), X1     // v1.Y
	MOVSS	v1Z+8(FP), X2     // v1.Z
	MOVSS	v2X+12(FP), X3    // v2.X
	MOVSS	v2Y+16(FP), X4    // v2.Y
	MOVSS	v2Z+20(FP), X5    // v2.Z
	
	// Calculate X component: v1.Y * v2.Z - v1.Z * v2.Y
	MOVSS	X1, X6
	MULSS	X5, X6            // X6 = v1.Y * v2.Z
	MOVSS	X2, X7
	MULSS	X4, X7            // X7 = v1.Z * v2.Y
	SUBSS	X7, X6            // X6 = X6 - X7
	
	// Calculate Y component: -(v1.X * v2.Z - v1.Z * v2.X)
	MOVSS	X0, X7
	MULSS	X5, X7            // X7 = v1.X * v2.Z
	MOVSS	X2, X8
	MULSS	X3, X8            // X8 = v1.Z * v2.X
	SUBSS	X8, X7            // X7 = X7 - X8
	XORPS	X8, X8            // X8 = 0
	SUBSS	X7, X8            // X8 = -X7
	
	// Calculate Z component: v1.X * v2.Y - v1.Y * v2.X
	MULSS	X4, X0            // X0 = v1.X * v2.Y
	MULSS	X3, X1            // X1 = v1.Y * v2.X
	SUBSS	X1, X0            // X0 = X0 - X1
	
	// Store results
	MOVSS	X6, ret0+24(FP)
	MOVSS	X8, ret1+28(FP)
	MOVSS	X0, ret2+32(FP)
	RET

