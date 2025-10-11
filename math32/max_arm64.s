//go:build arm64

#include "textflag.h"

// func max(x, y float32) float32
TEXT ·max(SB),NOSPLIT,$0-12
	// ABI wrapper stores inputs at 8(RSP) and 12(RSP), expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x
	FMOVS	12(RSP), F1      // Load y
	FMAXS	F0, F1, F0       // F0 = max(F0, F1)
	FMOVS	F0, 16(RSP)      // Store result
	RET

// func min(x, y float32) float32
TEXT ·min(SB),NOSPLIT,$0-12
	// ABI wrapper stores inputs at 8(RSP) and 12(RSP), expects output at 16(RSP)
	FMOVS	8(RSP), F0       // Load x
	FMOVS	12(RSP), F1      // Load y
	FMINS	F0, F1, F0       // F0 = min(F0, F1)
	FMOVS	F0, 16(RSP)      // Store result
	RET

