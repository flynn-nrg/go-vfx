//go:build arm64

#include "textflag.h"

// func floor(x float32) float32
TEXT ·floor(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(RSP) and expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x from stack
	FRINTMS	F0, F0           // Round toward minus infinity (floor)
	FMOVS	F0, 16(RSP)      // Store result
	RET

// func ceil(x float32) float32
TEXT ·ceil(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(RSP) and expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x from stack
	FRINTPS	F0, F0           // Round toward plus infinity (ceil)
	FMOVS	F0, 16(RSP)      // Store result
	RET

// func round(x float32) float32
TEXT ·round(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(RSP) and expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x from stack
	FRINTAS	F0, F0           // Round to nearest, ties away from zero
	FMOVS	F0, 16(RSP)      // Store result
	RET

