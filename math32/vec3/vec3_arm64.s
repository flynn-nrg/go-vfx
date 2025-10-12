//go:build arm64

#include "textflag.h"

// func dot(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) float32
TEXT ·dot(SB),NOSPLIT,$0-28
	// Load v1 components
	FMOVS	v1X+0(FP), F0
	FMOVS	v1Y+4(FP), F1
	FMOVS	v1Z+8(FP), F2
	
	// Multiply by v2 components
	FMOVS	v2X+12(FP), F3
	FMULS	F3, F0, F0        // F0 = v1.X * v2.X
	FMOVS	v2Y+16(FP), F3
	FMULS	F3, F1, F1        // F1 = v1.Y * v2.Y
	FMOVS	v2Z+20(FP), F3
	FMULS	F3, F2, F2        // F2 = v1.Z * v2.Z
	
	// Sum the results
	FADDS	F1, F0, F0        // F0 = F0 + F1
	FADDS	F2, F0, F0        // F0 = F0 + F2
	
	FMOVS	F0, ret+24(FP)
	RET

// func add3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·add3(SB),NOSPLIT,$0-36
	// Load v1
	FMOVS	v1X+0(FP), F0
	FMOVS	v1Y+4(FP), F1
	FMOVS	v1Z+8(FP), F2
	
	// Add v2 components
	FMOVS	v2X+12(FP), F3
	FADDS	F3, F0, F0
	FMOVS	v2Y+16(FP), F3
	FADDS	F3, F1, F1
	FMOVS	v2Z+20(FP), F3
	FADDS	F3, F2, F2
	
	// Store results
	FMOVS	F0, ret0+24(FP)
	FMOVS	F1, ret1+28(FP)
	FMOVS	F2, ret2+32(FP)
	RET

// func sub3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·sub3(SB),NOSPLIT,$0-36
	// Load v1
	FMOVS	v1X+0(FP), F0
	FMOVS	v1Y+4(FP), F1
	FMOVS	v1Z+8(FP), F2
	
	// Subtract v2 components
	FMOVS	v2X+12(FP), F3
	FSUBS	F3, F0, F0
	FMOVS	v2Y+16(FP), F3
	FSUBS	F3, F1, F1
	FMOVS	v2Z+20(FP), F3
	FSUBS	F3, F2, F2
	
	// Store results
	FMOVS	F0, ret0+24(FP)
	FMOVS	F1, ret1+28(FP)
	FMOVS	F2, ret2+32(FP)
	RET

// func mul3(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·mul3(SB),NOSPLIT,$0-36
	// Load v1
	FMOVS	v1X+0(FP), F0
	FMOVS	v1Y+4(FP), F1
	FMOVS	v1Z+8(FP), F2
	
	// Multiply by v2 components
	FMOVS	v2X+12(FP), F3
	FMULS	F3, F0, F0
	FMOVS	v2Y+16(FP), F3
	FMULS	F3, F1, F1
	FMOVS	v2Z+20(FP), F3
	FMULS	F3, F2, F2
	
	// Store results
	FMOVS	F0, ret0+24(FP)
	FMOVS	F1, ret1+28(FP)
	FMOVS	F2, ret2+32(FP)
	RET

// func scalarMul3(vX, vY, vZ, scalar float32) (float32, float32, float32)
TEXT ·scalarMul3(SB),NOSPLIT,$0-28
	// Load scalar
	FMOVS	scalar+12(FP), F3
	
	// Load vector and multiply by scalar
	FMOVS	vX+0(FP), F0
	FMOVS	vY+4(FP), F1
	FMOVS	vZ+8(FP), F2
	
	FMULS	F3, F0, F0
	FMULS	F3, F1, F1
	FMULS	F3, F2, F2
	
	// Store results
	FMOVS	F0, ret0+16(FP)
	FMOVS	F1, ret1+20(FP)
	FMOVS	F2, ret2+24(FP)
	RET

// func cross(v1X, v1Y, v1Z, v2X, v2Y, v2Z float32) (float32, float32, float32)
TEXT ·cross(SB),NOSPLIT,$0-36
	// Load v1 and v2
	FMOVS	v1X+0(FP), F0     // v1.X
	FMOVS	v1Y+4(FP), F1     // v1.Y
	FMOVS	v1Z+8(FP), F2     // v1.Z
	FMOVS	v2X+12(FP), F3    // v2.X
	FMOVS	v2Y+16(FP), F4    // v2.Y
	FMOVS	v2Z+20(FP), F5    // v2.Z
	
	// Calculate X component: v1.Y * v2.Z - v1.Z * v2.Y
	FMULS	F5, F1, F6        // F6 = v1.Y * v2.Z
	FMULS	F4, F2, F7        // F7 = v1.Z * v2.Y
	FSUBS	F7, F6, F6        // F6 = F6 - F7
	
	// Calculate Y component: -(v1.X * v2.Z - v1.Z * v2.X)
	FMULS	F5, F0, F7        // F7 = v1.X * v2.Z
	FMULS	F3, F2, F8        // F8 = v1.Z * v2.X
	FSUBS	F8, F7, F7        // F7 = F7 - F8
	FNEGS	F7, F7            // F7 = -F7
	
	// Calculate Z component: v1.X * v2.Y - v1.Y * v2.X
	FMULS	F4, F0, F8        // F8 = v1.X * v2.Y
	FMULS	F3, F1, F9        // F9 = v1.Y * v2.X
	FSUBS	F9, F8, F8        // F8 = F8 - F9
	
	// Store results
	FMOVS	F6, ret0+24(FP)
	FMOVS	F7, ret1+28(FP)
	FMOVS	F8, ret2+32(FP)
	RET

