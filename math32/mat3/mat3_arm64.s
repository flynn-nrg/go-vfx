//go:build arm64

#include "textflag.h"

// func matrixVectorMul(a11, a12, a13, a21, a22, a23, a31, a32, a33, vX, vY, vZ float32) (float32, float32, float32)
TEXT ·matrixVectorMul(SB),NOSPLIT,$0-52
	// Load vector components
	FMOVS	vX+36(FP), F3
	FMOVS	vY+40(FP), F4
	FMOVS	vZ+44(FP), F5
	
	// Calculate first row: a11*vX + a12*vY + a13*vZ using FMA
	FMOVS	a11+0(FP), F0
	FMOVS	a12+4(FP), F1
	FMOVS	a13+8(FP), F2
	FMULS	F3, F0, F0        // F0 = a11 * vX
	FMADDS	F4, F0, F1, F0    // F0 = F0 + (F1 * F4) = F0 + (a12 * vY)
	FMADDS	F5, F0, F2, F0    // F0 = F0 + (F2 * F5) = F0 + (a13 * vZ)
	
	// Calculate second row: a21*vX + a22*vY + a23*vZ using FMA
	FMOVS	a21+12(FP), F1
	FMOVS	a22+16(FP), F2
	FMOVS	a23+20(FP), F6
	FMULS	F3, F1, F1        // F1 = a21 * vX
	FMADDS	F4, F1, F2, F1    // F1 = F1 + (F2 * F4) = F1 + (a22 * vY)
	FMADDS	F5, F1, F6, F1    // F1 = F1 + (F6 * F5) = F1 + (a23 * vZ)
	
	// Calculate third row: a31*vX + a32*vY + a33*vZ using FMA
	FMOVS	a31+24(FP), F2
	FMOVS	a32+28(FP), F6
	FMOVS	a33+32(FP), F7
	FMULS	F3, F2, F2        // F2 = a31 * vX
	FMADDS	F4, F2, F6, F2    // F2 = F2 + (F6 * F4) = F2 + (a32 * vY)
	FMADDS	F5, F2, F7, F2    // F2 = F2 + (F7 * F5) = F2 + (a33 * vZ)
	
	// Store results
	FMOVS	F0, ret0+48(FP)
	FMOVS	F1, ret1+52(FP)
	FMOVS	F2, ret2+56(FP)
	RET

