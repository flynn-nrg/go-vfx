//go:build amd64

#include "textflag.h"

// func max(x, y float32) float32
TEXT ·max(SB),NOSPLIT,$0-12
	// ABI wrapper stores inputs at 8(RSP) and 12(RSP), expects output at 16(RSP)
	MOVSS	8(RSP), X0       // Load x
	MOVSS	12(RSP), X1      // Load y
	MAXSS	X1, X0           // X0 = max(X0, X1)
	MOVSS	X0, 16(RSP)      // Store result
	RET

// func min(x, y float32) float32
TEXT ·min(SB),NOSPLIT,$0-12
	// ABI wrapper stores inputs at 8(RSP) and 12(RSP), expects output at 16(RSP)
	MOVSS	8(RSP), X0       // Load x
	MOVSS	12(RSP), X1      // Load y
	MINSS	X1, X0           // X0 = min(X0, X1)
	MOVSS	X0, 16(RSP)      // Store result
	RET

