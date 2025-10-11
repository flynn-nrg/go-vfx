//go:build amd64

#include "textflag.h"

// func sqrt(x float32) float32
TEXT ·sqrt(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(RSP) and expects output at 16(RSP)
	MOVSS	8(RSP), X0       // Load x from 8(RSP) to XMM register
	SQRTSS	X0, X0           // X0 = sqrt(X0) using hardware instruction
	MOVSS	X0, 16(RSP)      // Store result to 16(RSP)
	RET

