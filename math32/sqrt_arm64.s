//go:build arm64

#include "textflag.h"

// func sqrt(x float32) float32
TEXT ·sqrt(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(RSP) and expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x directly from stack to float register
	FSQRTS	F0, F0           // F0 = sqrt(F0) using hardware instruction
	FMOVS	F0, 16(RSP)      // Store result directly from float register to stack
	RET

