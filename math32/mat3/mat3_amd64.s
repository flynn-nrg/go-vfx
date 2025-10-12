//go:build amd64

#include "textflag.h"

// func matrixVectorMul(a11, a12, a13, a21, a22, a23, a31, a32, a33, vX, vY, vZ float32) (float32, float32, float32)
TEXT ·matrixVectorMul(SB),NOSPLIT,$0-52
	// Load vector components
	MOVSS	vX+36(FP), X3
	MOVSS	vY+40(FP), X4
	MOVSS	vZ+44(FP), X5
	
	// Calculate first row: a11*vX + a12*vY + a13*vZ
	MOVSS	a11+0(FP), X0
	MULSS	X3, X0            // X0 = a11 * vX
	MOVSS	a12+4(FP), X1
	MULSS	X4, X1            // X1 = a12 * vY
	ADDSS	X1, X0            // X0 = X0 + X1
	MOVSS	a13+8(FP), X1
	MULSS	X5, X1            // X1 = a13 * vZ
	ADDSS	X1, X0            // X0 = final X
	
	// Calculate second row: a21*vX + a22*vY + a23*vZ
	MOVSS	a21+12(FP), X1
	MULSS	X3, X1            // X1 = a21 * vX
	MOVSS	a22+16(FP), X2
	MULSS	X4, X2            // X2 = a22 * vY
	ADDSS	X2, X1            // X1 = X1 + X2
	MOVSS	a23+20(FP), X2
	MULSS	X5, X2            // X2 = a23 * vZ
	ADDSS	X2, X1            // X1 = final Y
	
	// Calculate third row: a31*vX + a32*vY + a33*vZ
	MOVSS	a31+24(FP), X2
	MULSS	X3, X2            // X2 = a31 * vX
	MOVSS	a32+28(FP), X6
	MULSS	X4, X6            // X6 = a32 * vY
	ADDSS	X6, X2            // X2 = X2 + X6
	MOVSS	a33+32(FP), X6
	MULSS	X5, X6            // X6 = a33 * vZ
	ADDSS	X6, X2            // X2 = final Z
	
	// Store results
	MOVSS	X0, ret0+48(FP)
	MOVSS	X1, ret1+52(FP)
	MOVSS	X2, ret2+56(FP)
	RET

