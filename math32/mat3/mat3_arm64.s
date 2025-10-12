//go:build arm64

#include "textflag.h"

// func matrixVectorMul(a11, a12, a13, a21, a22, a23, a31, a32, a33, vX, vY, vZ float32) (float32, float32, float32)
TEXT ·matrixVectorMul(SB),NOSPLIT,$0-52
	// Load vector components
	FMOVS	vX+36(FP), F3
	FMOVS	vY+40(FP), F4
	FMOVS	vZ+44(FP), F5
	
	// Calculate first row: a11*vX + a12*vY + a13*vZ
	FMOVS	a11+0(FP), F0
	FMULS	F3, F0, F0        // F0 = a11 * vX
	FMOVS	a12+4(FP), F1
	FMULS	F4, F1, F1        // F1 = a12 * vY
	FADDS	F1, F0, F0        // F0 = F0 + F1
	FMOVS	a13+8(FP), F1
	FMULS	F5, F1, F1        // F1 = a13 * vZ
	FADDS	F1, F0, F0        // F0 = final X
	
	// Calculate second row: a21*vX + a22*vY + a23*vZ
	FMOVS	a21+12(FP), F1
	FMULS	F3, F1, F1        // F1 = a21 * vX
	FMOVS	a22+16(FP), F2
	FMULS	F4, F2, F2        // F2 = a22 * vY
	FADDS	F2, F1, F1        // F1 = F1 + F2
	FMOVS	a23+20(FP), F2
	FMULS	F5, F2, F2        // F2 = a23 * vZ
	FADDS	F2, F1, F1        // F1 = final Y
	
	// Calculate third row: a31*vX + a32*vY + a33*vZ
	FMOVS	a31+24(FP), F2
	FMULS	F3, F2, F2        // F2 = a31 * vX
	FMOVS	a32+28(FP), F6
	FMULS	F4, F6, F6        // F6 = a32 * vY
	FADDS	F6, F2, F2        // F2 = F2 + F6
	FMOVS	a33+32(FP), F6
	FMULS	F5, F6, F6        // F6 = a33 * vZ
	FADDS	F6, F2, F2        // F2 = final Z
	
	// Store results
	FMOVS	F0, ret0+48(FP)
	FMOVS	F1, ret1+52(FP)
	FMOVS	F2, ret2+56(FP)
	RET

