//go:build amd64

#include "textflag.h"

// func floor(x float32) float32
TEXT ·floor(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(SP) and expects output at 16(SP)
	MOVSS	8(SP), X0        // Load x from stack
	ROUNDSS	$1, X0, X0       // Round toward -Inf (floor), mode = 0x01
	MOVSS	X0, 16(SP)       // Store result
	RET

// func ceil(x float32) float32
TEXT ·ceil(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(SP) and expects output at 16(SP)
	MOVSS	8(SP), X0        // Load x from stack
	ROUNDSS	$2, X0, X0       // Round toward +Inf (ceil), mode = 0x02
	MOVSS	X0, 16(SP)       // Store result
	RET

// func round(x float32) float32
// ROUNDSS doesn't have "ties away from zero" mode, so we implement it as:
// - For x >= 0: floor(x + 0.5)
// - For x < 0: ceil(x - 0.5)
TEXT ·round(SB),NOSPLIT,$0-8
	MOVSS	8(SP), X0        // Load x from stack
	
	// Check sign
	MOVSS	X0, X1           // Copy x
	XORPS	X2, X2           // X2 = 0
	CMPSS	$1, X2, X1       // X1 = (x < 0) ? 0xFFFFFFFF : 0
	
	// Create 0.5 with appropriate sign
	MOVL	$0x3F000000, AX  // 0.5 in float32 bits
	MOVL	AX, X3
	MOVL	$0xBF000000, AX  // -0.5 in float32 bits  
	MOVL	AX, X4
	
	// Select 0.5 or -0.5 based on sign
	ANDPS	X1, X4           // X4 = x<0 ? -0.5 : 0
	ANDNPS	X3, X1           // X1 = x>=0 ? 0.5 : 0
	ORPS	X4, X1           // X1 = x<0 ? -0.5 : 0.5
	
	// Add and truncate
	ADDSS	X1, X0           // X0 = x + (x<0 ? -0.5 : 0.5)
	ROUNDSS	$3, X0, X0       // Truncate (round toward zero)
	
	MOVSS	X0, 16(SP)       // Store result
	RET

