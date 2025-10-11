//go:build amd64

#include "textflag.h"

// func sqrt(x float32) float32
TEXT ·sqrt(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(SP) and expects output at 16(SP)
	MOVSS	8(SP), X0        // Load x from 8(SP) to XMM register
	SQRTSS	X0, X0           // X0 = sqrt(X0) using hardware instruction
	MOVSS	X0, 16(SP)       // Store result to 16(SP)
	RET

